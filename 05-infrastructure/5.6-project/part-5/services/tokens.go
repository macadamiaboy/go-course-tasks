package services

import (
	"context"
	"fmt"

	"5.6/part-5/db"
	"5.6/part-5/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenService struct {
	tokenRepo db.RefreshTokenRepository
	userRepo  db.UserRepository
	pb.UnimplementedTokenServiceServer
}

func (s *TokenService) IssueToken(ctx context.Context, req *pb.IssueTokenRequest) (*pb.IssueTokenResponse, error) {
	user, err := s.userRepo.GetByLogin(ctx, req.Login)
	if err != nil {
		return nil, status.Error(codes.NotFound, "couldn't find the user with the provided login")
	}

	// is hash(req.Password) == user.password?
	// gen token

	tokenHash := fmt.Sprintf("hashed_token_for_email_%s", user.Email)

	refToken, err := s.tokenRepo.Create(ctx, db.RefreshToken{UserID: user.ID, TokenHash: tokenHash, IsRevoked: false})
	accToken := ""

	return &pb.IssueTokenResponse{
		AccessToken:  accToken,
		RefreshToken: refToken.TokenHash,
	}, nil
}

func (s *TokenService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	_ = req.Token

	// check if it's valid

	return &pb.ValidateTokenResponse{
		IsValid: true,
		Subject: "",
	}, nil
}

func (s *TokenService) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	token, err := s.tokenRepo.FindActive(ctx, db.RefreshToken{UserID: req.UserID})
	if err != nil {
		return &pb.RevokeTokenResponse{Error: err.Error()}, status.Error(codes.NotFound, "couldn't find active token for the provided user")
	}

	err = s.tokenRepo.Revoke(ctx, token)
	if err != nil {
		return &pb.RevokeTokenResponse{Error: err.Error()}, status.Error(codes.Internal, "couldn't revoke the token")
	}

	return &pb.RevokeTokenResponse{
		Error: "",
	}, nil
}
