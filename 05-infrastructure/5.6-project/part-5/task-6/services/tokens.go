package services

import (
	"context"
	"errors"
	"log"
	"time"

	"5.6/task-6/db"
	"5.6/task-6/helpers"
	"5.6/task-6/pkg/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var LoginError = errors.New("failed to login the user")
var ErrFailedToValidate = errors.New("validation faliure")

type TokenService struct {
	tokenRepo db.RefreshTokenRepository
	userRepo  db.UserRepository
	pb.UnimplementedTokenServiceServer
}

func NewTokenService(tokenRepo db.RefreshTokenRepository, userRepo db.UserRepository) *TokenService {
	return &TokenService{
		tokenRepo: tokenRepo,
		userRepo:  userRepo,
	}
}

func (s *TokenService) IssueToken(ctx context.Context, req *pb.IssueTokenRequest) (*pb.IssueTokenResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Login)
	if err != nil {
		log.Println("failed to get")
		return nil, status.Error(codes.NotFound, "couldn't find the user with the provided login")
	}
	/*
		if !helpers.CheckPasswordHash(req.Password, user.PasswordHash) {
			return nil, LoginError
		}
	*/
	accessToken, err := helpers.GenToken(user.ID)
	if err != nil {
		log.Println("failed to gen acc")
		return nil, status.Error(codes.Internal, "failed to gen the access token")
	}

	refreshToken, err := helpers.GenRefreshToken()
	if err != nil {
		log.Println("failed to gen ref")
		return nil, status.Error(codes.Internal, "failed to gen the refresh token")
	}

	tokenHash := helpers.HashToken(refreshToken)
	expiresAt := time.Now().Add(15 * time.Minute)

	_, err = s.tokenRepo.Create(ctx, db.RefreshToken{UserID: user.ID, TokenHash: tokenHash, ExpiresAt: expiresAt})
	if err != nil {
		log.Println("failed to save")
		return nil, status.Error(codes.InvalidArgument, "failed to save: invalid data provided")
	}

	return &pb.IssueTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *TokenService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	badResult := &pb.ValidateTokenResponse{IsValid: false}

	hashedIncomingToken := helpers.HashToken(req.Token)
	tokenData, err := s.tokenRepo.GetByHash(ctx, db.RefreshToken{TokenHash: hashedIncomingToken})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return badResult, status.Error(codes.NotFound, "there's no such token in the db")
		}
		return badResult, status.Error(codes.InvalidArgument, "failed to get: invalid data provided")
	}

	if tokenData.IsRevoked {
		return badResult, status.Error(codes.Unauthenticated, "token in not valid")
	}

	if tokenData.ExpiresAt.Before(time.Now()) {
		return badResult, status.Error(codes.Unauthenticated, "token in not valid")
	}

	return &pb.ValidateTokenResponse{IsValid: true}, nil
}

func (s *TokenService) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*empty.Empty, error) {
	token, err := s.tokenRepo.FindActive(ctx, db.RefreshToken{UserID: req.UserID})
	if err != nil {
		return &empty.Empty{}, status.Error(codes.InvalidArgument, "there's no such active token")
	}

	err = s.tokenRepo.Revoke(ctx, token)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.Internal, "failed to revoke")
	}

	return &empty.Empty{}, nil
}
