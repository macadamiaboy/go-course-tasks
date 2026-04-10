package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenRequest struct {
	User_id string `json:"user_id"`
}

type apiError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

var secretSalt = []byte("super_secret_salt")

func genToken(id string) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		Issuer:    "MerchShop",
		Subject:   id,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(secretSalt)
	return s, err
}

func main() {
	mux := http.NewServeMux()

	//5.1.1
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	//5.1.2
	mux.HandleFunc("GET /api/v1/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			User_id string `json:"user_id"`
		}{User_id: r.PathValue("id")}

		writeJSON(w, http.StatusOK, response)
	})

	//5.1.3
	mux.HandleFunc("POST /api/v1/tokens", func(w http.ResponseWriter, r *http.Request) {
		var request tokenRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{"invalid json"})
			return
		}

		token, err := genToken(request.User_id)
		if err != nil {
			log.Printf("failed to create the token")
			writeJSON(w, http.StatusInternalServerError, apiError{"cannot process"})
			return
		}

		response := struct {
			Access_token string `json:"access_token"`
		}{Access_token: token}

		writeJSON(w, http.StatusCreated, response)
	})

	//5.1.6
	mux.HandleFunc("GET /api/v1/tokens/{id}", func(w http.ResponseWriter, r *http.Request) {
		tokenMock := struct {
			ID        string    `json:"id"`
			Name      string    `json:"name"`
			CreatedAt time.Time `json:"created_at"`
			Status    string    `json:"status"`
		}{
			ID:        r.PathValue("id"),
			Name:      "Mock-Token",
			CreatedAt: time.Now(),
			Status:    "active",
		}

		writeJSON(w, http.StatusOK, tokenMock)
	})

	//5.1.4
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		// Читает заголовки запроса. Позволяет оперативно отсечь висящие соединения
		ReadHeaderTimeout: 2 * time.Second,
		// Читает весь запрос. Дает возможность защититься от медленных клиентов
		ReadTimeout: 5 * time.Second,
		// Время на запись ответа. От клиентов, медленно читающих ответ
		WriteTimeout: 10 * time.Second,
		// Время ожидания между двумя запросами в режиме keep-alive. Экономиит время на TCP-рукопожатии
		IdleTimeout: 30 * time.Second,
	}

	//5.1.5
	go func() {
		log.Println("server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
