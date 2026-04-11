package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware func(http.Handler) http.Handler

type contextKey string

const userClaimsKey contextKey = "user-claims-key"

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

type apiError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// token verifying interface and structs for jwt n paseto
type TokenVerifier interface {
	Verify(token string) (jwt.Claims, error)
}

type JWTVerifier struct {
	secretKey []byte
}

func NewJWTVerifier(secret string) *JWTVerifier {
	return &JWTVerifier{secretKey: []byte(secret)}
}

type PasetoVerifier struct {
	symmetricKey paseto.V4SymmetricKey
	parser       paseto.Parser
}

func NewPasetoVerifier(hexKey string) (*PasetoVerifier, error) {
	key, err := paseto.V4SymmetricKeyFromHex(hexKey)
	if err != nil {
		return nil, err
	}

	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	return &PasetoVerifier{
		symmetricKey: key,
		parser:       parser,
	}, nil
}

// 5.2.3
func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing the token header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("wrong token header format")
	}

	res := strings.TrimSpace(headerParts[1])
	return res, nil
}

// 5.2.1
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func testMiddleware(name string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[%s] before\n", name)
			next.ServeHTTP(w, r)
			fmt.Printf("[%s] after\n", name)
		})
	}
}

// 5.2.2
func loggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sRecorder := statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(&sRecorder, r)

			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", sRecorder.status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

// 5.2.4
func authMiddleware(verifier TokenVerifier, logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := extractBearerToken(r.Header.Get("Authorization"))
			if err != nil {
				logger.Error("failed to get the token", "err", err)
				writeJSON(w, http.StatusBadRequest, apiError{"failed to get the token"})
				return
			}

			claims, err := verifier.Verify(token)
			if err != nil {
				logger.Error("failed to verify the token", "err", err)
				writeJSON(w, http.StatusUnauthorized, apiError{"failed to verify the token"})
				return
			}

			ctxWithClaims := context.WithValue(r.Context(), userClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctxWithClaims))
		})
	}
}

// 5.2.5 Two implementations for jwt n paseto
func (v *JWTVerifier) Verify(tokenString string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return v.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (v *PasetoVerifier) Verify(token string) (jwt.Claims, error) {
	parsedToken, err := v.parser.ParseV4Local(v.symmetricKey, token, nil)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	for key, value := range parsedToken.Claims() {
		claims[key] = value
	}

	return claims, nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "everything's ok")
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, r.Context().Value(userClaimsKey))
}

func main() {
	mux := http.NewServeMux()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	//handler := Chain(http.HandlerFunc(healthHandler), testMiddleware("mw1"), testMiddleware("mw2"))
	jwtVerifier := NewJWTVerifier("super_secret_salt")
	pasetoVerifier, err := NewPasetoVerifier("707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f")
	if err != nil {
		logger.Error("failed to create the paseto verifier", "err", err)
		return
	}
	_ = pasetoVerifier

	meHandler := Chain(http.HandlerFunc(meHandler), loggingMiddleware(logger), authMiddleware(jwtVerifier, logger))
	//meHandler := Chain(http.HandlerFunc(healthHandler), loggingMiddleware(logger), authMiddleware(pasetoVerifier, logger))
	healthHandler := Chain(http.HandlerFunc(healthHandler), testMiddleware("mw1"), testMiddleware("mw2"), loggingMiddleware(logger))

	mux.Handle("GET /health", healthHandler)

	mux.Handle("GET /api/v1/me", meHandler)

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
