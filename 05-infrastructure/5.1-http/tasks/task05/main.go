// Задание: Graceful shutdown
//
// Добавь корректное завершение сервера при получении SIGINT / SIGTERM.
// При сигнале запускай srv.Shutdown с таймаутом и логируй этапы.
//
// Ожидаемый результат:
//   $ go run main.go
//   server started on :8080
//
//   (после Ctrl+C)
//   shutting down...
//   server stopped

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// TODO: запусти сервер в горутине через srv.ListenAndServe()
	// При ошибке, отличной от http.ErrServerClosed, вызывай log.Fatal

	// TODO: создай канал для перехвата сигналов:
	//   quit := make(chan os.Signal, 1)
	//   signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// TODO: заблокируйся на чтении из quit

	// TODO: после получения сигнала:
	//   1. Залогируй "shutting down..."
	//   2. Создай контекст с таймаутом 10 секунд
	//   3. Вызови srv.Shutdown(ctx)
	//   4. Залогируй "server stopped"

	_ = srv            // убери после реализации
	_ = os.Signal(syscall.SIGINT)
	_ = signal.Notify
	_ = context.Background

	log.Println("server started on :8080")
}
