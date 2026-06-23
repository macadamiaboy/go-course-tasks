package util

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	Error string `json:"error"`
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func GetRoutePattern(r *http.Request) string {
	if r.Pattern != "" {
		return r.Pattern
	}

	return "unknown"
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, apiError{msg})
}
