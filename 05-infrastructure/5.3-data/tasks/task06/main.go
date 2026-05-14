// Задание: Интеграционная мини-задача — data-слой Token Service
//
// Спроектируй и реализуй data-слой для Token Service.
//
// Таблицы:
//   users         (id, email, password_hash, created_at)
//   refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at)
//
// Методы репозитория:
//   UserRepository:
//     CreateUser(ctx, email, passwordHash string) (User, error)
//     GetUserByEmail(ctx, email string) (User, error)
//
//   TokenRepository:
//     CreateRefreshToken(ctx, userID int64, tokenHash string, expiresAt time.Time) (RefreshToken, error)
//     GetActiveToken(ctx, tokenHash string) (RefreshToken, error)  — только не отозванные и не истёкшие
//     RevokeToken(ctx, tokenHash string) error
//
// Требования:
//   - PostgreSQL + pgxpool
//   - type-safe слой через sqlc (или ручной pgx)
//   - миграции в sql/migrations/
//   - транзакция для сценария "создание пользователя + выдача refresh-token"
//   - краткий README с описанием схемы
//
// Структура директории:
//   task06/
//     main.go              ← этот файл
//     go.mod
//     sql/
//       migrations/
//         001_users.up.sql
//         001_users.down.sql
//         002_tokens.up.sql
//         002_tokens.down.sql
//       queries/
//         users.sql
//         tokens.sql
//     db/                  ← сгенерированный sqlc-код (или ручные репозитории)
//     README.md

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenRevoked  = errors.New("token revoked")
	ErrTokenExpired  = errors.New("token expired")
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// TODO: реализуй UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

// TODO: реализуй TokenRepository
type TokenRepository interface {
	CreateRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) (RefreshToken, error)
	GetActiveToken(ctx context.Context, tokenHash string) (RefreshToken, error)
	RevokeToken(ctx context.Context, tokenHash string) error
}

// TODO: реализуй RegisterWithToken(ctx, pool, userRepo, tokenRepo, email, pwdHash, tokenHash string) error
// Выполни в одной транзакции:
//   1. CreateUser
//   2. CreateRefreshToken (expires через 30 дней)
// При ошибке — rollback

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://dev:dev@localhost:5432/devdb"
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Println("connect error:", err)
		return
	}
	defer pool.Close()

	// TODO: создай реализации репозиториев и вызови RegisterWithToken
	fmt.Println("TODO: implement repositories and call RegisterWithToken")
}
