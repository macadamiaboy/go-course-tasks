package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(ctx context.Context, userParams User) (User, error) {
	query := `INSERT INTO users(email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash;`
	var user User

	err := ur.db.QueryRow(ctx, query, userParams.Email, userParams.PasswordHash).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		return User{}, ErrCannotCreate
	}

	return user, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`
	var user User

	err := ur.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNoRows
		}
		return User{}, ErrCannotFind
	}

	return user, nil
}
