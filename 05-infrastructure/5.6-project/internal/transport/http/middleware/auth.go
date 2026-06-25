package middleware

import (
	"log/slog"
	"net/http"

	"github.com/course/token-service/internal/helpers"
	"github.com/course/token-service/internal/transport/util"
	"github.com/course/token-service/internal/transport/util/auth"
)

func AuthMiddleware(logger *slog.Logger, ac helpers.AppConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := helpers.ExtractBearerToken(r.Header.Get("Authorization"))
			if err != nil {
				logger.Error("failed to get the token", "err", err)
				util.WriteError(w, http.StatusBadRequest, "failed to get the token")
				return
			}

			userID, err := helpers.GetIdFromToken(token, ac)
			if err != nil {
				logger.Error("failed to verify the token", "err", err)
				util.WriteError(w, http.StatusUnauthorized, "failed to verify the token")
				return
			}

			ctxWithUserID := auth.WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctxWithUserID))
		})
	}
}
