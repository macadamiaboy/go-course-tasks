package interceptor

import (
	"context"

	"github.com/course/token-service/internal/helpers"
	"github.com/course/token-service/internal/transport/util/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthUnaryInterceptor(ac helpers.AppConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "the token is missing")
		}

		authHeader := md.Get("authorization")

		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "the token is missing")
		}

		token, err := helpers.ExtractBearerToken(authHeader[0])
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "wrong token format")
		}

		userID, err := helpers.GetIdFromToken(token, ac)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		newCtx := auth.WithUserID(ctx, userID)

		return handler(newCtx, req)
	}
}
