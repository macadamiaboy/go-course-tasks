package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
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

type rateBucket struct {
	count   int
	resetAt time.Time
}

func RateLimitMiddleware(maxRequests int, window time.Duration) Middleware {
	var mu sync.Mutex
	buckets := map[string]*rateBucket{}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)

			mu.Lock()
			b, ok := buckets[host]
			if !ok || time.Now().After(b.resetAt) {
				b = &rateBucket{resetAt: time.Now().Add(window)}
				buckets[host] = b
			}
			b.count++
			count := b.count
			mu.Unlock()

			if count > maxRequests {
				writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "rate limit exceeded"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "pong"})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /ping", Chain(
		http.HandlerFunc(pingHandler),
		RateLimitMiddleware(3, 10*time.Second),
	))

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
