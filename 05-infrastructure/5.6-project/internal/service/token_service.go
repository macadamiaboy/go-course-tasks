package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/course/token-service/internal/db"
	"github.com/course/token-service/internal/helpers"
	"github.com/course/token-service/internal/observability"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var ErrBadRequest = errors.New("bad request")       //400
var ErrUnauthorized = errors.New("failed to login") //401
var ErrNotFound = errors.New("not found")           //404
var RefreshExpTime = 30 * time.Minute

var tracer = otel.Tracer("5.6-project/service/token-service")

type TokenService struct {
	pool      *pgxpool.Pool
	queries   *db.Queries
	appConfig helpers.AppConfig
	metrics   *observability.AuthMetrics
}

func NewTokenService(pool *pgxpool.Pool, q *db.Queries, ac helpers.AppConfig, m *observability.AuthMetrics) *TokenService {
	return &TokenService{pool: pool, queries: q, appConfig: ac, metrics: m}
}

func (ts *TokenService) Register(ctx context.Context, email, password string) (userID int64, err error) {
	ctx, span := tracer.Start(ctx, "TokenService.Register")
	defer span.End()

	failureReason := "none"

	defer func() {
		if err != nil {
			wrappedErr := fmt.Errorf("%s: %w", failureReason, err)
			span.SetStatus(codes.Error, wrappedErr.Error())
			span.RecordError(wrappedErr)
		}
	}()

	refreshToken, err := helpers.GenRefreshToken()
	if err != nil {
		failureReason = "failed_to_gen"
		return 0, err
	}

	hashed := helpers.HashToken(refreshToken)

	pwdHash, err := helpers.HashPassword(password)
	if err != nil {
		failureReason = "failed_to_hash"
		return 0, err
	}

	userID, err = ts.registerWithToken(ctx, email, pwdHash, hashed)
	if err != nil {
		failureReason = "failed_to_save"
		return 0, err
	}

	return userID, nil
}

func (ts *TokenService) registerWithToken(ctx context.Context, email, pwdHash, tokenHash string) (int64, error) {
	ctx, span := tracer.Start(ctx, "TokenService.registerWithToken")
	defer span.End()

	tx, err := ts.pool.Begin(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return 0, errors.New("failed to start the tx")
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	qtx := ts.queries.WithTx(tx)

	createUserParams := db.CreateUserParams{Email: email, PasswordHash: pwdHash}
	user, err := qtx.CreateUser(ctx, createUserParams)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return 0, ErrBadRequest
	}

	expAt := pgtype.Timestamptz{Time: time.Now().Add(RefreshExpTime), Valid: true}
	createTokenParams := db.CreateRefreshTokenParams{UserID: user.ID, TokenHash: tokenHash, ExpiresAt: expAt}
	_, err = qtx.CreateRefreshToken(ctx, createTokenParams)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return 0, ErrBadRequest
	}

	err = tx.Commit(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return 0, errors.New("failed to commit the tx")
	}

	return user.ID, nil
}

func (ts *TokenService) Login(ctx context.Context, email, password string) (refreshToken, accessToken string, err error) {
	ctx, span := tracer.Start(ctx, "TokenService.Login")
	defer span.End()

	failureReason := "none"

	defer func() {
		if err != nil {
			wrappedErr := fmt.Errorf("%s: %w", failureReason, err)
			span.SetStatus(codes.Error, wrappedErr.Error())
			span.RecordError(wrappedErr)

			ts.metrics.TokensIssued.WithLabelValues("refresh", "fail", failureReason).Inc()
		} else {
			ts.metrics.TokensIssued.WithLabelValues("refresh", "success", "none").Inc()
			ts.metrics.TokensIssued.WithLabelValues("access", "success", "none").Inc()
		}
	}()

	user, err := ts.queries.GetUserByEmail(ctx, email)
	if err != nil {
		failureReason = "user_not_found"
		return "", "", ErrUnauthorized
	}

	if !helpers.CheckPasswordHash(password, user.PasswordHash) {
		failureReason = "unauthorized"
		return "", "", ErrUnauthorized
	}

	accessToken, err = helpers.GenToken(user.ID, ts.appConfig)
	if err != nil {
		ts.metrics.TokensIssued.WithLabelValues("access", "fail", "failed_to_gen").Inc()
		return "", "", err
	}

	refreshToken, err = helpers.GenRefreshToken()
	if err != nil {
		failureReason = "failed_to_gen"
		return "", "", err
	}

	hashed := helpers.HashToken(refreshToken)

	expAt := pgtype.Timestamptz{Time: time.Now().Add(RefreshExpTime), Valid: true}
	createTokenParams := db.CreateRefreshTokenParams{UserID: user.ID, TokenHash: hashed, ExpiresAt: expAt}
	_, err = ts.queries.CreateRefreshToken(ctx, createTokenParams)
	if err != nil {
		failureReason = "failed_to_save"
		return "", "", ErrBadRequest
	}

	return refreshToken, accessToken, nil
}

func (ts *TokenService) Refresh(ctx context.Context, oldRefresh string) (newRefresh, accessToken string, err error) {
	ctx, span := tracer.Start(ctx, "TokenService.Refresh")
	defer span.End()

	failureReason := "none"

	defer func() {
		if err != nil {
			wrappedErr := fmt.Errorf("%s: %w", failureReason, err)
			span.SetStatus(codes.Error, wrappedErr.Error())
			span.RecordError(wrappedErr)

			ts.metrics.TokensIssued.WithLabelValues("refresh", "fail", failureReason).Inc()
		} else {
			ts.metrics.TokensIssued.WithLabelValues("refresh", "success", "none").Inc()
			ts.metrics.TokensIssued.WithLabelValues("access", "success", "none").Inc()
		}
	}()

	tx, err := ts.pool.Begin(ctx)
	if err != nil {
		failureReason = "transaction_failed"
		return "", "", errors.New("failed to start the tx")
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	qtx := ts.queries.WithTx(tx)

	oldHashed := helpers.HashToken(oldRefresh)

	active, err := qtx.FindActiveRefreshToken(ctx, oldHashed)
	if err != nil {
		failureReason = "active_token_not_found"
		return "", "", ErrUnauthorized
	}

	err = qtx.RevokeRefreshToken(ctx, oldHashed)
	if err != nil {
		failureReason = "failed_to_revoke"
		return "", "", errors.New("failed to revoke the token")
	}

	accessToken, err = helpers.GenToken(active.UserID, ts.appConfig)
	if err != nil {
		ts.metrics.TokensIssued.WithLabelValues("access", "fail", "failed_to_gen").Inc()
		return "", "", err
	}

	newRefresh, err = helpers.GenRefreshToken()
	if err != nil {
		failureReason = "failed_to_gen"
		return "", "", err
	}

	newHashed := helpers.HashToken(newRefresh)

	expAt := pgtype.Timestamptz{Time: time.Now().Add(RefreshExpTime), Valid: true}
	createTokenParams := db.CreateRefreshTokenParams{UserID: active.UserID, TokenHash: newHashed, ExpiresAt: expAt}
	_, err = qtx.CreateRefreshToken(ctx, createTokenParams)
	if err != nil {
		failureReason = "failed_to_save"
		return "", "", errors.New("failed to revoke the token")
	}

	err = tx.Commit(ctx)
	if err != nil {
		failureReason = "failed_to_commit"
		return "", "", errors.New("failed to commit the tx")
	}

	return newRefresh, accessToken, nil
}

func (ts *TokenService) Issue(ctx context.Context, userID int64) (token string, err error) {
	ctx, span := tracer.Start(ctx, "TokenService.Issue")
	defer span.End()

	failureReason := "none"

	defer func() {
		if err != nil {
			wrappedErr := fmt.Errorf("%s: %w", failureReason, err)
			span.SetStatus(codes.Error, wrappedErr.Error())
			span.RecordError(wrappedErr)

			ts.metrics.TokensIssued.WithLabelValues("refresh", "fail", failureReason).Inc()
		} else {
			ts.metrics.TokensIssued.WithLabelValues("refresh", "success", "none").Inc()
		}
	}()

	refreshToken, err := helpers.GenRefreshToken()
	if err != nil {
		failureReason = "failed_to_gen"
		return "", err
	}

	hashed := helpers.HashToken(refreshToken)

	expAt := pgtype.Timestamptz{Time: time.Now().Add(RefreshExpTime), Valid: true}
	createTokenParams := db.CreateRefreshTokenParams{UserID: userID, TokenHash: hashed, ExpiresAt: expAt}
	_, err = ts.queries.CreateRefreshToken(ctx, createTokenParams)
	if err != nil {
		failureReason = "failed_to_save"
		return "", errors.New("failed to create")
	}

	return refreshToken, nil
}

func (ts *TokenService) Validate(ctx context.Context, token string) (int64, error) {
	ctx, span := tracer.Start(ctx, "TokenService.Validate")
	defer span.End()

	hashed := helpers.HashToken(token)

	active, err := ts.queries.FindActiveRefreshToken(ctx, hashed)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return 0, ErrUnauthorized
	}

	return active.UserID, nil
}

func (ts *TokenService) Revoke(ctx context.Context, tokenToRevoke string) error {
	ctx, span := tracer.Start(ctx, "TokenService.Revoke")
	defer span.End()

	tx, err := ts.pool.Begin(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return errors.New("failed to start the tx")
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	qtx := ts.queries.WithTx(tx)

	hashed := helpers.HashToken(tokenToRevoke)

	_, err = qtx.FindActiveRefreshToken(ctx, hashed)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return ErrUnauthorized
	}

	err = qtx.RevokeRefreshToken(ctx, hashed)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return errors.New("failed to revoke the token")
	}

	err = tx.Commit(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return errors.New("failed to commit the tx")
	}

	return nil
}

func (ts *TokenService) GetUser(ctx context.Context, id int64) (string, error) {
	ctx, span := tracer.Start(ctx, "TokenService.GetUser")
	defer span.End()

	user, err := ts.queries.GetUserByID(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return "", ErrNotFound
	}

	return user.Email, nil
}
