package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"5.6/task-6/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

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

	mux.Handle("GET /auth", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		respStatus := http.StatusOK

		resp, err := client.IssueToken(ctx, &pb.IssueTokenRequest{Login: "specificUser"})
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
				default:
					respStatus = http.StatusInternalServerError
				}
			}
			w.WriteHeader(respStatus)
			_, _ = w.Write([]byte(strconv.Itoa(respStatus)))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(respStatus)
		_ = json.NewEncoder(w).Encode(resp)
	}))

	_ = http.ListenAndServe(":8080", mux)
}
