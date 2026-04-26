package service

import (
	"complex-task/db"
	"complex-task/helpers"
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
)

var LoginError = errors.New("failed to login the user")
var ErrFailedToValidate = errors.New("validation faliure")

type TokenService struct {
	tokenRepo db.RefreshTokenRepository
	userRepo  db.UserRepository
}

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PairOfTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ValidationRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ValidationResponse struct {
	IsValid bool `json:"is_valid"`
}

type RevokeTokenRequest struct {
	UserId int64 `json:"user_id"`
}

func NewTokenService(tokenRepo db.RefreshTokenRepository, userRepo db.UserRepository) *TokenService {
	return &TokenService{
		tokenRepo: tokenRepo,
		userRepo:  userRepo,
	}
}

func (s *TokenService) IssueToken(ctx context.Context, req UserData) (*PairOfTokens, error) {
	_, span := otel.Tracer("token-service").Start(ctx, "Service.IssueToken")
	defer span.End()

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, LoginError
	}
	/*
		if !helpers.CheckPasswordHash(req.Password, user.PasswordHash) {
			return nil, LoginError
		}
	*/
	accessToken, err := helpers.GenToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("issueToken, genToken err: %w", err)
	}

	refreshToken, err := helpers.GenRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("issueToken, genRefreshToken err: %w", err)
	}

	tokenHash := helpers.HashToken(refreshToken)
	expiresAt := time.Now().Add(15 * time.Minute)

	_, err = s.tokenRepo.Create(ctx, db.RefreshToken{UserID: user.ID, TokenHash: tokenHash, ExpiresAt: expiresAt})
	if err != nil {
		return nil, fmt.Errorf("issueToken, cannot save the token, err: %w", err)
	}

	return &PairOfTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *TokenService) ValidateToken(ctx context.Context, tokenReq ValidationRequest) (*ValidationResponse, error) {
	badResult := &ValidationResponse{IsValid: false}

	hashedIncomingToken := helpers.HashToken(tokenReq.RefreshToken)
	tokenData, err := s.tokenRepo.GetByHash(ctx, db.RefreshToken{TokenHash: hashedIncomingToken})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return badResult, fmt.Errorf("%w, not found", ErrFailedToValidate)
		}
		return badResult, fmt.Errorf("validateToken, failed to find the token, err: %w", err)
	}

	if tokenData.IsRevoked {
		return badResult, fmt.Errorf("%w, is revoked", ErrFailedToValidate)
	}

	if tokenData.ExpiresAt.Before(time.Now()) {
		return badResult, fmt.Errorf("%w, has been expired", ErrFailedToValidate)
	}

	return &ValidationResponse{IsValid: true}, nil
}

func (s *TokenService) RevokeToken(ctx context.Context, userIDReq RevokeTokenRequest) error {
	token, err := s.tokenRepo.FindActive(ctx, db.RefreshToken{UserID: userIDReq.UserId})
	if err != nil {
		return fmt.Errorf("revokeToken, couldn't find the active token, err: %w", err)
	}

	err = s.tokenRepo.Revoke(ctx, token)
	if err != nil {
		return fmt.Errorf("revokeToken, couldn't revoke the token, err: %w", err)
	}

	return nil
}
