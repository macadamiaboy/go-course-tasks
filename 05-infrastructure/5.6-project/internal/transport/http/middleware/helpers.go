package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func getRoutePattern(r *http.Request) string {
	if r.Pattern != "" {
		return r.Pattern
	}

	return "unknown"
}
