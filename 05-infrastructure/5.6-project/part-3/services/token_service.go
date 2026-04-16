package main

import (
	"context"
	"time"

	"5.6/part-3/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//var ErrNotEnoughFunds = errors.New("not enough funds on the account")

type TokenService struct {
	pool           *pgxpool.Pool
	refresh_tokens db.Queries
	users          db.Queries
}

func NewTokenService(pool *pgxpool.Pool, tokens db.Queries, users db.Queries) *TokenService {
	return &TokenService{pool: pool, refresh_tokens: tokens, users: users}
}

func (ts *TokenService) Login(ctx context.Context, email string, password string) error {
	tx, err := ts.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	refresh_tokens := ts.refresh_tokens.WithTx(tx)
	users := ts.users.WithTx(tx)

	user, err := users.GetUser(ctx, email)
	if err != nil {
		return err
	}

	// check user.password_hash == hash(password)
	// if not then 401
	// else:

	//gen token hash
	tknHash := "tokenhash"
	timeToExpire := time.Now().AddDate(0, 1, 0)
	refresh_tokens.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: tknHash,
		ExpiresAt: pgtype.Timestamptz{Time: timeToExpire},
	})

	return tx.Commit(ctx)
}
