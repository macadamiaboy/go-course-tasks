package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// importing this data from env
const username string = "postgres"
const password string = "password"
const host string = "localhost"
const port int = 5432
const dbName string = "test_db"

func main() {
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
		log.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	// call methods
	// one of methods use from sqlc (db package)

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	log.Println("db is reachable")
}
