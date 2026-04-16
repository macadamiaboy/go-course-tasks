package refresh_tokens

import "5.6/part-3/db"

// generally useless
// all the methods are made with sqlc so they are located at db/query.sql.go

type RefreshTokenRepository struct {
	db db.DBTX
}

func NewRefreshTokenRepository(db db.DBTX) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

/*
type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	IsRevoked bool
	ExpiresAt time.Time
}

func (rtr *RefreshTokenRepository) Create(ctx context.Context, rtParams RefreshToken) (RefreshToken, error) {
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
}

func (rtr *RefreshTokenRepository) Revoke(ctx context.Context, rtParams RefreshToken) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE is_revoked = false AND user_id = $1;`

	_, err := rtr.db.Exec(ctx, query, rtParams.UserID)
	if err != nil {
		return err
	}

	return nil
}

// i don't check the expires_at. Should I?
func (rtr *RefreshTokenRepository) FindActive(ctx context.Context, rtParams RefreshToken) error {
	query := `SELECT (id, user_id, token_hash, is_revoked, expires_at) FROM refresh_tokens WHERE is_revoked = false AND user_id = $1 LIMIT 1;`

	var refToken RefreshToken

	err := rtr.db.QueryRow(ctx, query, rtParams.UserID).Scan(
		&refToken.ID,
		&refToken.UserID,
		&refToken.TokenHash,
		&refToken.IsRevoked,
		&refToken.ExpiresAt,
	)
	if err != nil {
		return err
	}

	return nil
}
*/
