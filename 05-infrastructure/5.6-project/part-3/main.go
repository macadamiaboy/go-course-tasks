package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"5.6/part-3/db"
	"5.6/part-3/db/account"
	"5.6/part-3/db/transfer"
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

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	log.Println("db is reachable")

	// call methods
	// one of methods use from sqlc (db package)
	query := db.New(pool)
	fUserToCreate := db.CreateUserParams{Email: "1@yandex.com", PasswordHash: "zdfnjviuefnf"}
	firstCreatedUser, err := query.CreateUser(context.Background(), fUserToCreate)
	if err != nil {
		log.Fatalf("failed to create the user: %v", err)
	}

	sUserToCreate := db.CreateUserParams{Email: "2@yandex.com", PasswordHash: "zdfnjviuefnl"}
	secondCreatedUser, err := query.CreateUser(context.Background(), sUserToCreate)
	if err != nil {
		log.Fatalf("failed to create the user: %v", err)
	}

	accs := account.NewAccRepository(pool)
	fAccToCreate := account.Account{UserID: firstCreatedUser.ID, Coins: 700}
	_, err = accs.Create(context.Background(), fAccToCreate)
	if err != nil {
		log.Fatalf("failed to create the acc: %v", err)
	}

	sAccToCreate := account.Account{UserID: secondCreatedUser.ID, Coins: 900}
	_, err = accs.Create(context.Background(), sAccToCreate)
	if err != nil {
		log.Fatalf("failed to create the acc: %v", err)
	}

	transfers := transfer.NewTransferRepository(pool)
	fTransferToCreate := transfer.Transfer{SourceID: secondCreatedUser.ID, TargetID: firstCreatedUser.ID, Amount: 700}
	_, err = transfers.Create(context.Background(), fTransferToCreate)
	if err != nil {
		log.Fatalf("failed to create the transfer: %v", err)
	}

	sTransferToCreate := transfer.Transfer{SourceID: firstCreatedUser.ID, TargetID: secondCreatedUser.ID, Amount: 1500}
	_, err = transfers.Create(context.Background(), sTransferToCreate)
	if err != nil {
		log.Printf("failed to create the transfer: %v", err)
	}
}
