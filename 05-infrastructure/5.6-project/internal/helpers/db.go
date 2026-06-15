package helpers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dsn = struct {
	username string
	password string
	host     string
	port     int
	dbName   string
}{username: "postgres", password: "password", host: "localhost", port: 5432, dbName: "token_service_db"}

func InitDB(baseCtx context.Context, logger *slog.Logger) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		dsn.username,
		dsn.password,
		dsn.host,
		dsn.port,
		dsn.dbName,
	)

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
