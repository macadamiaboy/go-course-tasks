// Задание: Recovery middleware
//
// Реализуй RecoveryMiddleware, которая перехватывает паники в handler'ах.
//
// Требования:
// 1. Используй defer + recover()
// 2. Если паника произошла — верни клиенту 500 {"error":"internal server error"}
// 3. Залогируй факт паники через slog (уровень Error, поле "panic")
//
// Важно: WriteHeader нужно вызвать явно — иначе клиент получит 200 с пустым телом.
//
// Ожидаемый результат:
//   $ go run main.go &
//   server started on :8080
//
//   $ curl http://localhost:8080/safe
//   {"status":"ok"}
//
//   $ curl http://localhost:8080/panic
//   {"error":"internal server error"}
//   (в логах: {"level":"ERROR","msg":"panic recovered","panic":"something went wrong"})
//
//   $ curl http://localhost:8080/safe
//   {"status":"ok"}   ← сервер продолжает работать!

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

// TODO: реализуй RecoveryMiddleware(logger *slog.Logger) Middleware
//
// Подсказка по структуре:
//   func RecoveryMiddleware(logger *slog.Logger) Middleware {
//       return func(next http.Handler) http.Handler {
//           return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//               defer func() {
//                   if p := recover(); p != nil {
//                       logger.Error("panic recovered", "panic", p)
//                       writeJSON(w, http.StatusInternalServerError, ...)
//                   }
//               }()
//               next.ServeHTTP(w, r)
//           })
//       }
//   }

func safeHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func panicHandler(_ http.ResponseWriter, _ *http.Request) {
	panic("something went wrong")
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	_ = logger // убери после реализации

	mux := http.NewServeMux()

	// TODO: оберни оба handler'а в RecoveryMiddleware
	// mux.Handle("GET /safe", Chain(http.HandlerFunc(safeHandler), RecoveryMiddleware(logger)))
	// mux.Handle("GET /panic", Chain(http.HandlerFunc(panicHandler), RecoveryMiddleware(logger)))

	mux.Handle("GET /safe", http.HandlerFunc(safeHandler))   // убери после реализации
	mux.Handle("GET /panic", http.HandlerFunc(panicHandler)) // убери после реализации

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
