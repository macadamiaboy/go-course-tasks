package transfer

import (
	"context"

	"5.6/part-3/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transfer struct {
	ID       int64
	SourceID int64
	TargetID int64
	Amount   int
}

type TransferRepository struct {
	db db.DBTX
}

func NewTransferRepository(pool *pgxpool.Pool) *TransferRepository {
	return &TransferRepository{db: pool}
}

func (tyr *TransferRepository) WithTx(tx pgx.Tx) *TransferRepository {
	return &TransferRepository{
		db: tx,
	}
}

func (tr *TransferRepository) Create(ctx context.Context, transferParams Transfer) (Transfer, error) {
	query := `INSERT INTO transfers(source, target, amount) VALUES ($1, $2, $3) RETURNING id, source, target, amount;`
	var transfer Transfer

	err := tr.db.QueryRow(ctx, query, transferParams.SourceID, transferParams.TargetID, transferParams.Amount).Scan(
		&transfer.ID,
		&transfer.SourceID,
		&transfer.TargetID,
		&transfer.Amount,
	)
	if err != nil {
		return Transfer{}, err
	}
	return transfer, nil
}
