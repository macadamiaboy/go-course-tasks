package db

import (
	"context"
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID    int64
	Email string
}

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(ctx context.Context, userParams User) (User, error) {
	/*
		query := `INSERT INTO users(email) VALUES ($1) RETURNING id, email;`
		var user User

		err := ur.db.QueryRow(ctx, query, userParams.Email).Scan(
			&user.ID,
			&user.Email,
		)
		if err != nil {
			return User{}, err
		}
		return user, nil
	*/
	time.Sleep(100 * time.Millisecond)
	return User{}, nil
}

func (ur *UserRepository) GetByLogin(ctx context.Context, login string) (User, error) {
	/*
		query := `SELECT id, email FROM users WHERE id = $1`
		var user User

		err := ur.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return User{}, ErrUserNotFound
			}
			return User{}, err
		}
		return user, nil
	*/
	time.Sleep(100 * time.Millisecond)
	if login == "unknown" {
		return User{}, errors.New("unknown user")
	}
	return User{ID: 14, Email: login}, nil
}
