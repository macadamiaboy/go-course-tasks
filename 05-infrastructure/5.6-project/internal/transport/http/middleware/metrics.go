package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/course/token-service/internal/observability"
	"github.com/course/token-service/internal/transport/util"
)

func MetricsMiddleware(metrics *observability.HTTPMetrics) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &util.StatusRecorder{ResponseWriter: w, Status: http.StatusOK}

			next.ServeHTTP(rec, r)

			status := strconv.Itoa(rec.Status)
			route := getRoutePattern(r)

			metrics.RequestsTotal.WithLabelValues(r.Method, route, status).Inc()
			metrics.RequestsDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
		})
	}
}
