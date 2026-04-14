package account

import (
	"context"
	"errors"

	"5.6/part-3/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAccNotFound = errors.New("account not found")

type Account struct {
	ID     int64
	UserID int64
	Coins  int
}

type AccountRepository struct {
	db db.DBTX
}

func NewAccRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: pool}
}

func (ar *AccountRepository) WithTx(tx pgx.Tx) *AccountRepository {
	return &AccountRepository{
		db: tx,
	}
}

func (ar *AccountRepository) Create(ctx context.Context, accParams Account) (Account, error) {
	query := `INSERT INTO accounts(user_id, coins) VALUES ($1, $2) RETURNING id, user_id, coins;`
	var acc Account

	err := ar.db.QueryRow(ctx, query, accParams.UserID, accParams.Coins).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.Coins,
	)
	if err != nil {
		return Account{}, err
	}
	return acc, nil
}

func (ar *AccountRepository) GetBalanceByUserID(ctx context.Context, userID int64) (int, error) {
	query := `SELECT coins FROM accounts WHERE user_id = $1`
	var balance int

	err := ar.db.QueryRow(ctx, query, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrAccNotFound
		}
		return 0, err
	}
	return balance, nil
}

func (ar *AccountRepository) setBalanceByUserID(ctx context.Context, userID int64, balance int) (Account, error) {
	query := `UPDATE accounts SET coins = $2 WHERE user_id = $1 RETURNING id, user_id, coins;`
	var res Account

	err := ar.db.QueryRow(ctx, query, userID, balance).Scan(&res.ID, &res.UserID, &res.Coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Account{}, ErrAccNotFound
		}
		return Account{}, err
	}
	return res, nil
}

func (ar *AccountRepository) changeBalance(ctx context.Context, userID int64, sum int, action func(int, int) int) (Account, error) {
	balance, err := ar.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return Account{}, err
	}

	newBalance := action(balance, sum)

	updatedAcc, err := ar.setBalanceByUserID(ctx, userID, newBalance)
	if err != nil {
		return Account{}, err
	}

	return updatedAcc, nil
}

func (ar *AccountRepository) CreditTo(ctx context.Context, userId int64, sum int) (Account, error) {
	return ar.changeBalance(ctx, userId, sum, func(a, b int) int {
		return a + b
	})
}

func (ar *AccountRepository) WriteOff(ctx context.Context, userId int64, sum int) (Account, error) {
	return ar.changeBalance(ctx, userId, sum, func(a, b int) int {
		return a - b
	})
}
