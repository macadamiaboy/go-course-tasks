package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func RecoveryMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					logger.Error("panic recovered", "panic", p)
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func safeHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func panicHandler(_ http.ResponseWriter, _ *http.Request) {
	panic("something went wrong")
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.Handle("GET /safe", Chain(http.HandlerFunc(safeHandler), RecoveryMiddleware(logger)))
	mux.Handle("GET /panic", Chain(http.HandlerFunc(panicHandler), RecoveryMiddleware(logger)))

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
