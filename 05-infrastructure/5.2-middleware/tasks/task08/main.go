// Задание: Rate limiting middleware
//
// Реализуй RateLimitMiddleware — ограничение числа запросов с одного IP.
//
// Алгоритм скользящего окна (fixed window):
//   - Для каждого IP-адреса храни счётчик и время сброса (resetAt).
//   - Если time.Now().After(resetAt) — сброси счётчик и обнови resetAt.
//   - Если счётчик >= maxRequests — верни 429 Too Many Requests.
//   - Иначе — увеличь счётчик и пропусти запрос дальше.
//
// Структуры:
//   type rateBucket struct {
//       count   int
//       resetAt time.Time
//   }
//   Используй map[string]*rateBucket + sync.Mutex.
//
// IP-адрес: используй r.RemoteAddr ("host:port") и net.SplitHostPort для извлечения host.
//
// Ожидаемый результат (лимит 3 запроса за 10 секунд):
//   $ go run main.go &
//   server started on :8080
//
//   $ for i in {1..5}; do curl -s http://localhost:8080/ping; echo; done
//   {"status":"pong"}
//   {"status":"pong"}
//   {"status":"pong"}
//   {"error":"rate limit exceeded"}
//   {"error":"rate limit exceeded"}

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

// TODO: объяви структуру rateBucket
// type rateBucket struct {
//     count   int
//     resetAt time.Time
// }

// TODO: реализуй RateLimitMiddleware(maxRequests int, window time.Duration) Middleware
//
// Подсказка:
//   func RateLimitMiddleware(maxRequests int, window time.Duration) Middleware {
//       var mu sync.Mutex
//       buckets := map[string]*rateBucket{}
//
//       return func(next http.Handler) http.Handler {
//           return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//               host, _, _ := net.SplitHostPort(r.RemoteAddr)
//               mu.Lock()
//               // ... проверяй и обновляй bucket ...
//               mu.Unlock()
//
//               if превышен лимит {
//                   writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "rate limit exceeded"})
//                   return
//               }
//               next.ServeHTTP(w, r)
//           })
//       }
//   }

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "pong"})
}

func main() {
	_ = sync.Mutex{}        // подсказка: защищает map
	_ = net.SplitHostPort   // подсказка: извлекает host из RemoteAddr
	_ = time.Now            // подсказка: для проверки resetAt

	mux := http.NewServeMux()

	// TODO: оберни pingHandler в RateLimitMiddleware(3, 10*time.Second)
	// mux.Handle("GET /ping", Chain(http.HandlerFunc(pingHandler), RateLimitMiddleware(3, 10*time.Second)))

	mux.Handle("GET /ping", http.HandlerFunc(pingHandler)) // убери после реализации

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
