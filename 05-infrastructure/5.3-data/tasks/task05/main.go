// Задание: Транзакционная операция — Transfer
//
// Реализуй Transfer(ctx, fromID, toID, amount) в одной транзакции.
// При ошибке выполняй rollback; возвращай контекстные ошибки через %w.
//
// Ожидаемый результат:
//   transfer 300 from Alice to Bob: ok
//   Alice after: 700
//   Bob after:   800
//   transfer 2000 (too much): insufficient funds
//   transfer from unknown account: account not found

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type Account struct {
	ID      int64
	Owner   string
	Balance int64
}

// TODO: реализуй Transfer(ctx, pool, fromID, toID, amount int64) error
// Шаги:
//   1. Начать транзакцию: tx, err := pool.Begin(ctx)
//   2. Отложить rollback: defer tx.Rollback(ctx)
//   3. Получить from-аккаунт через SELECT ... FOR UPDATE (блокировка строки)
//   4. Проверить баланс >= amount, иначе вернуть ErrInsufficientFunds
//   5. Получить to-аккаунт
//   6. Обновить оба баланса через UPDATE
//   7. Выполнить tx.Commit(ctx)
// Подсказка: pgx.ErrNoRows → ErrAccountNotFound через errors.Is

func Transfer(ctx context.Context, pool *pgxpool.Pool, fromID, toID, amount int64) error {
	// TODO: implement
	_ = pgx.ErrNoRows // подсказка
	return nil
}

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

	// Предполагается, что в БД уже есть таблица accounts и тестовые данные
	// CREATE TABLE accounts (id bigserial primary key, owner text, balance bigint);
	// INSERT INTO accounts (owner, balance) VALUES ('Alice', 1000), ('Bob', 500);

	if err = Transfer(ctx, pool, 1, 2, 300); err != nil {
		fmt.Println("transfer 300 from Alice to Bob:", err)
	} else {
		fmt.Println("transfer 300 from Alice to Bob: ok")
	}

	if err = Transfer(ctx, pool, 1, 2, 2000); err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			fmt.Println("transfer 2000 (too much): insufficient funds")
		} else {
			fmt.Println("unexpected error:", err)
		}
	}

	if err = Transfer(ctx, pool, 9999, 2, 100); err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			fmt.Println("transfer from unknown account: account not found")
		} else {
			fmt.Println("unexpected error:", err)
		}
	}
}
