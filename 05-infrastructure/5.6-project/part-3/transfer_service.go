package main

import (
	"context"
	"errors"
	"fmt"

	"5.6/part-3/account"
	"5.6/part-3/transfer"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotEnoughFunds = errors.New("not enough funds on the account")

type TransferService struct {
	pool      *pgxpool.Pool
	accounts  account.AccountRepository
	transfers transfer.TransferRepository
}

func (ts *TransferService) Transfer(ctx context.Context, sourceID, targetID int64, sum int) error {
	tx, err := ts.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	accounts := ts.accounts.WithTx(tx)
	transfers := ts.transfers.WithTx(tx)

	balance, err := accounts.GetBalanceByUserID(ctx, sourceID)
	if err != nil {
		return err
	}
	if balance < sum {
		return ErrNotEnoughFunds
	}

	_, err = accounts.WriteOff(ctx, sourceID, sum)
	if err != nil {
		return fmt.Errorf("failed to write off: %w", err)
	}

	_, err = accounts.CreditTo(ctx, targetID, sum)
	if err != nil {
		return fmt.Errorf("failed to credit to: %w", err)
	}

	transfer := transfer.Transfer{SourceID: sourceID, TargetID: targetID, Amount: sum}
	_, err = transfers.Create(ctx, transfer)
	if err != nil {
		return fmt.Errorf("failed to create the record: %w", err)
	}

	return tx.Commit(ctx)
}
