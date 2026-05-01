package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrCannotCreate = errors.New("failed to create the record")
var ErrCannotUpdate = errors.New("failed to update the record")
var ErrCannotFind = errors.New("failed to find the record")
var ErrNoRows = errors.New("there's no such record")

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	IsRevoked bool
	ExpiresAt time.Time
}
