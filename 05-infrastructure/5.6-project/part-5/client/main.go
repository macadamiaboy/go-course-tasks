package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"5.6/part-5/pkg/pb"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	resp, err := client.IssueToken(ctx, &pb.IssueTokenRequest{Login: "specificUser"})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				fmt.Printf("Illegal args provided: %s\n", st.Message())
			case codes.DeadlineExceeded:
				fmt.Println("Timeout exceeded")
			case codes.NotFound:
				fmt.Println("The user is not found")
			default:
				fmt.Printf("Some other grpc error: %v - %s\n", st.Code(), st.Message())
			}
		}

		return
	}

	log.Printf("Received data: %s", resp.RefreshToken)
}
