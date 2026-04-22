package db

import (
	"context"
	"time"
)

type RefreshTokenRepository struct {
	db DBTX
}

func NewRefreshTokenRepository(db DBTX) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	IsRevoked bool
	ExpiresAt time.Time
}

func (rtr *RefreshTokenRepository) Create(ctx context.Context, rtParams RefreshToken) (RefreshToken, error) {
	/*
		query := `INSERT INTO refresh_tokens(user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id, user_id, token_hash, is_revoked, expires_at;`
		var refToken RefreshToken

		err := rtr.db.QueryRow(ctx, query, rtParams.UserID, rtParams.TokenHash, rtParams.ExpiresAt).Scan(
			&refToken.ID,
			&refToken.UserID,
			&refToken.TokenHash,
			&refToken.IsRevoked,
			&refToken.ExpiresAt,
		)
		if err != nil {
			return RefreshToken{}, err
		}
		return refToken, nil
	*/
	time.Sleep(100 * time.Millisecond)
	return rtParams, nil
}

func (rtr *RefreshTokenRepository) Revoke(ctx context.Context, rtParams RefreshToken) error {
	/*
		query := `UPDATE refresh_tokens SET is_revoked = true WHERE id = $1;`

		_, err := rtr.db.Exec(ctx, query, rtParams.ID)
		if err != nil {
			return err
		}
	*/
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (rtr *RefreshTokenRepository) FindActive(ctx context.Context, rtParams RefreshToken) (RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, is_revoked, expires_at FROM refresh_tokens WHERE is_revoked = false AND user_id = $1 LIMIT 1;`

	var refToken RefreshToken

	err := rtr.db.QueryRow(ctx, query, rtParams.UserID).Scan(
		&refToken.ID,
		&refToken.UserID,
		&refToken.TokenHash,
		&refToken.IsRevoked,
		&refToken.ExpiresAt,
	)
	if err != nil {
		return RefreshToken{}, err
	}

	return refToken, nil
}
