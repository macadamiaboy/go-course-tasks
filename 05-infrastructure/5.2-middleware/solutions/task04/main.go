package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Claims struct {
	UserID string
	Role   string
}

type contextKey string

const claimsKey contextKey = "claims"

type TokenVerifier interface {
	Verify(token string) (Claims, error)
}

type mockVerifier struct{}

func (v *mockVerifier) Verify(token string) (Claims, error) {
	if token == "valid-token" {
		return Claims{UserID: "user-123", Role: "admin"}, nil
	}
	return Claims{}, errors.New("invalid token")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header is empty")
	}
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok || strings.TrimSpace(token) == "" {
		return "", errors.New("invalid authorization header format")
	}
	return strings.TrimSpace(token), nil
}

func AuthMiddleware(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, err := extractBearerToken(header)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
				return
			}

			claims, err := verifier.Verify(token)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(Claims)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"user_id": claims.UserID,
		"role":    claims.Role,
	})
}

func main() {
	verifier := &mockVerifier{}

	mux := http.NewServeMux()
	mux.Handle("GET /api/v1/me", AuthMiddleware(verifier)(http.HandlerFunc(meHandler)))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("server error:", err)
	}
}
