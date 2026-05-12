package user

import (
	"context"
	"errors"

	"5.6/part-3/db"
	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID    int64
	Email string
}

type UserRepository struct {
	db db.DBTX
}

func NewUserRepository(db db.DBTX) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) WithTx(tx pgx.Tx) *UserRepository {
	return &UserRepository{
		db: tx,
	}
}

func (ur *UserRepository) Create(ctx context.Context, userParams User) (User, error) {
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
}

func (ur *UserRepository) GetByID(ctx context.Context, id int64) (User, error) {
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
}
