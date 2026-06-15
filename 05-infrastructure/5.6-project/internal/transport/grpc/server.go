package grpc

import (
	"context"
	"errors"
	"strconv"

	"github.com/course/token-service/internal/service"
	"github.com/course/token-service/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenServiceServer struct {
	tokenService service.TokenService
	pb.UnimplementedTokenServiceServer
}

func NewTokenServiceServer(tokenService service.TokenService) *TokenServiceServer {
	return &TokenServiceServer{tokenService: tokenService}
}

func (s *TokenServiceServer) IssueToken(ctx context.Context, in *pb.IssueTokenRequest) (*pb.IssueTokenResponse, error) {
	userIDstr := in.GetUserId()

	if userIDstr == "" {
		return nil, status.Error(codes.InvalidArgument, "id is not provided")
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong id format")
	}

	token, err := s.tokenService.Issue(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create the token")
	}

	return &pb.IssueTokenResponse{Token: token, UserId: userIDstr}, nil
}

func (s *TokenServiceServer) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token := in.GetToken()

	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is not provided")
	}

	userID, err := s.tokenService.Validate(ctx, token)
	if err != nil {
		return &pb.ValidateTokenResponse{UserId: "", IsValid: false}, status.Error(codes.Unauthenticated, "invalid token")
	}

	userIDstr := strconv.FormatInt(userID, 10)

	return &pb.ValidateTokenResponse{UserId: userIDstr, IsValid: true}, nil
}

func (s *TokenServiceServer) RevokeToken(ctx context.Context, in *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	token := in.GetToken()

	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is not provided")
	}

	err := s.tokenService.Revoke(ctx, token)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			return &pb.RevokeTokenResponse{IsRevoked: false}, status.Error(codes.Unauthenticated, "invalid token")
		}
		return &pb.RevokeTokenResponse{IsRevoked: false}, status.Error(codes.Internal, "internal server error")
	}

	return &pb.RevokeTokenResponse{IsRevoked: true}, nil
}
