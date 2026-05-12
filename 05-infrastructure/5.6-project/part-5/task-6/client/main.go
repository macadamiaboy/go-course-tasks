package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"5.6/task-6/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type apiError struct {
	Error string `json:"error"`
}

func authHandler(client pb.TokenServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestBody pb.IssueTokenRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{"incorrect request body"})
			return
		}

		respStatus := http.StatusOK

		resp, err := client.IssueToken(ctx, &requestBody)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				switch st.Code() {
				case codes.InvalidArgument:
					respStatus = http.StatusBadRequest
				case codes.DeadlineExceeded:
					respStatus = http.StatusRequestTimeout
				case codes.NotFound:
					respStatus = http.StatusNotFound
				case codes.Unauthenticated:
					respStatus = http.StatusUnauthorized
				default:
					respStatus = http.StatusInternalServerError
				}
			}
			writeJSON(w, respStatus, apiError{st.Message()})
			return
		}

		writeJSON(w, respStatus, resp)
	}
}

func revokeHandler(client pb.TokenServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestBody pb.RevokeTokenRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{"incorrect request body"})
			return
		}

		respStatus := http.StatusOK

		resp, err := client.RevokeToken(ctx, &requestBody)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				switch st.Code() {
				case codes.InvalidArgument:
					respStatus = http.StatusBadRequest
				case codes.DeadlineExceeded:
					respStatus = http.StatusRequestTimeout
				case codes.NotFound:
					respStatus = http.StatusNotFound
				case codes.Unauthenticated:
					respStatus = http.StatusUnauthorized
				default:
					respStatus = http.StatusInternalServerError
				}
			}
			writeJSON(w, respStatus, apiError{st.Message()})
			return
		}

		writeJSON(w, respStatus, resp)
	}
}

func validateHandler(client pb.TokenServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestBody pb.ValidateTokenRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{"incorrect request body"})
			return
		}

		respStatus := http.StatusOK

		resp, err := client.ValidateToken(ctx, &requestBody)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				switch st.Code() {
				case codes.InvalidArgument:
					respStatus = http.StatusBadRequest
				case codes.DeadlineExceeded:
					respStatus = http.StatusRequestTimeout
				case codes.NotFound:
					respStatus = http.StatusNotFound
				case codes.Unauthenticated:
					respStatus = http.StatusUnauthorized
				default:
					respStatus = http.StatusInternalServerError
				}
			}
			writeJSON(w, respStatus, apiError{st.Message()})
			return
		}

		writeJSON(w, respStatus, resp)
	}
}

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewTokenServiceClient(conn)

	mux := http.NewServeMux()

	mux.Handle("POST /auth", http.HandlerFunc(authHandler(client)))
	mux.Handle("GET /validate", http.HandlerFunc(validateHandler(client)))
	mux.Handle("POST /revoke", http.HandlerFunc(revokeHandler(client)))

	_ = http.ListenAndServe(":8080", mux)
}
