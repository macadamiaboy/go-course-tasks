package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/course/token-service/internal/transport/util"
)

func LoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sRecorder := util.StatusRecorder{ResponseWriter: w, Status: http.StatusOK}
			start := time.Now()

			next.ServeHTTP(&sRecorder, r)

			logger.Info("http request",
				"method", r.Method,
				"path", util.GetRoutePattern(r),
				"status", sRecorder.Status,
				"duration", time.Since(start),
			)
		})
	}
}
