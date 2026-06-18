package grpc

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/course/token-service/internal/service"
	"github.com/course/token-service/pkg/pb"
	"go.opentelemetry.io/otel"
	otelcodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tracer = otel.Tracer("/grpc/server")

type TokenServiceServer struct {
	tokenService service.TokenService
	logger       *slog.Logger
	pb.UnimplementedTokenServiceServer
}

func NewTokenServiceServer(tokenService service.TokenService, logger *slog.Logger) *TokenServiceServer {
	return &TokenServiceServer{tokenService: tokenService, logger: logger}
}

func (s *TokenServiceServer) IssueToken(ctx context.Context, in *pb.IssueTokenRequest) (*pb.IssueTokenResponse, error) {
	ctx, span := tracer.Start(ctx, "IssueToken")
	defer span.End()

	userIDstr := in.GetUserId()

	if userIDstr == "" {
		errString := "id is not provided"
		span.SetStatus(otelcodes.Error, errString)
		span.RecordError(errors.New(errString))
		s.logger.ErrorContext(ctx, "IssueToken: invalid argument", "error", errString)

		return nil, status.Error(codes.InvalidArgument, errString)
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		errString := "wrong id format"
		span.SetStatus(otelcodes.Error, err.Error())
		span.RecordError(err)
		s.logger.ErrorContext(ctx, errString, "error", err.Error())

		return nil, status.Error(codes.InvalidArgument, errString)
	}

	token, err := s.tokenService.Issue(ctx, userID)
	if err != nil {
		errString := "failed to create the token"
		span.SetStatus(otelcodes.Error, err.Error())
		span.RecordError(err)
		s.logger.ErrorContext(ctx, errString, "error", err.Error())

		return nil, status.Error(codes.Internal, errString)
	}

	return &pb.IssueTokenResponse{Token: token, UserId: userIDstr}, nil
}

func (s *TokenServiceServer) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	ctx, span := tracer.Start(ctx, "ValidateToken")
	defer span.End()

	token := in.GetToken()

	if token == "" {
		errString := "token is not provided"
		span.SetStatus(otelcodes.Error, errString)
		span.RecordError(errors.New(errString))
		s.logger.ErrorContext(ctx, "ValidateToken: invalid argument", "error", errString)

		return nil, status.Error(codes.InvalidArgument, errString)
	}

	userID, err := s.tokenService.Validate(ctx, token)
	if err != nil {
		errString := "invalid token"
		span.SetStatus(otelcodes.Error, err.Error())
		span.RecordError(err)
		s.logger.ErrorContext(ctx, errString, "error", err.Error())

		return &pb.ValidateTokenResponse{UserId: "", IsValid: false}, status.Error(codes.Unauthenticated, errString)
	}

	userIDstr := strconv.FormatInt(userID, 10)

	return &pb.ValidateTokenResponse{UserId: userIDstr, IsValid: true}, nil
}

func (s *TokenServiceServer) RevokeToken(ctx context.Context, in *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	ctx, span := tracer.Start(ctx, "RevokeToken")
	defer span.End()

	token := in.GetToken()

	if token == "" {
		errString := "token is not provided"
		span.SetStatus(otelcodes.Error, errString)
		span.RecordError(errors.New(errString))
		s.logger.ErrorContext(ctx, "RevokeToken: invalid argument", "error", errString)

		return nil, status.Error(codes.InvalidArgument, errString)
	}

	err := s.tokenService.Revoke(ctx, token)
	if err != nil {
		span.SetStatus(otelcodes.Error, err.Error())
		span.RecordError(err)

		if errors.Is(err, service.ErrUnauthorized) {
			errString := "invalid token"
			s.logger.ErrorContext(ctx, errString, "error", err.Error())

			return &pb.RevokeTokenResponse{IsRevoked: false}, status.Error(codes.Unauthenticated, errString)
		}
		errString := "internal server error"
		s.logger.ErrorContext(ctx, errString, "error", err.Error())
		return &pb.RevokeTokenResponse{IsRevoked: false}, status.Error(codes.Internal, errString)
	}

	return &pb.RevokeTokenResponse{IsRevoked: true}, nil
}
