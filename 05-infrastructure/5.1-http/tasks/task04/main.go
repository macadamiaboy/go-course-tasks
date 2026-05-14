// Задание: Базовые таймауты HTTP-сервера
//
// Перепиши запуск сервера с http.Server вместо http.ListenAndServe.
// Задай все ключевые таймауты и объясни их в комментариях.
//
// Ожидаемый результат:
//   $ go run main.go
//   server started on :8080
//
//   $ curl http://localhost:8080/health
//   ok

package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// TODO: создай http.Server с адресом ":8080" и полем Handler: mux
	// Задай следующие таймауты (и объясни в комментарии, почему такие значения):
	//
	//   ReadHeaderTimeout: ... // время на чтение заголовков запроса
	//   ReadTimeout:       ... // время на чтение всего запроса (заголовки + тело)
	//   WriteTimeout:      ... // время на запись ответа клиенту
	//   IdleTimeout:       ... // время простоя keep-alive соединения
	//
	// Рекомендуемые значения для старта: 5s / 10s / 10s / 60s

	// TODO: запусти сервер через srv.ListenAndServe()
	// и обработай ошибку

	_ = time.Second // убери после реализации

	log.Println("server started on :8080")
}
