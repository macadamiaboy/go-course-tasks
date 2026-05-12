package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"5.6/task-6/db"
	"5.6/task-6/pkg/pb"
	"5.6/task-6/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

// importing this data from env
const username string = "postgres"
const password string = "password"
const host string = "localhost"
const port int = 5432
const dbName string = "test_db"

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		dbName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Printf("database connection failed: %v\n", err)
	}
	defer pool.Close()

	userRepo := db.NewUserRepository(pool)
	tokenRepo := db.NewRefreshTokenRepository(pool)
	tokenService := services.NewTokenService(*tokenRepo, *userRepo)

	s := grpc.NewServer()

	pb.RegisterTokenServiceServer(s, tokenService)

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
