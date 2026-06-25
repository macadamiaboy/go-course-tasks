package helpers

import (
	"context"
	"log/slog"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(baseCtx context.Context, logger *slog.Logger, DBConfig DBConfig) (*pgxpool.Pool, error) {
	url := DBConfig.ConnString

	ctx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		logger.Error("pgx config parsing failed", "error", err)
		return nil, err
	}

	config.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithTrimSQLInSpanName(),
	)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("pgx pool opening failed", "error", err)
		return nil, err
	}

	return pool, nil
}
