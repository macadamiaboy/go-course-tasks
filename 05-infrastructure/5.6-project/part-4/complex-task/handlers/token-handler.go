package handlers

import (
	"complex-task/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
)

type TokenHandler struct {
	service *service.TokenService
	logger  *slog.Logger
}

func NewTokenHandler(ts *service.TokenService, logger *slog.Logger) *TokenHandler {
	return &TokenHandler{service: ts, logger: logger}
}

type apiError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (th *TokenHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("token-service").Start(r.Context(), "RefreshHandler")
	defer span.End()

	traceID := span.SpanContext().TraceID().String()

	var requestBody service.UserData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		errString := "incorrect request body"
		th.logger.Error(errString, "trace_id", traceID, "error", err.Error())
		writeJSON(w, http.StatusBadRequest, apiError{errString})
		return
	}

	tokens, err := th.service.IssueToken(ctx, requestBody)
	if err != nil {
		if errors.Is(err, service.LoginError) {
			th.logger.Error("incorrect login or password", "trace_id", traceID, "error", err.Error())
			writeJSON(w, http.StatusUnauthorized, apiError{"incorrect login or password"})
			return
		}
		th.logger.Error("internal error", "trace_id", traceID, "error", err.Error())
		writeJSON(w, http.StatusInternalServerError, apiError{"internal error"})
		return
	}

	th.logger.Info("successfully executed", "trace_id", traceID)
	writeJSON(w, http.StatusCreated, service.PairOfTokens{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
}

func (th *TokenHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var requestBody service.ValidationRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		errString := "incorrect request body"
		th.logger.Error(errString, "error", err.Error())
		writeJSON(w, http.StatusBadRequest, apiError{errString})
		return
	}

	status := http.StatusOK
	validResult, err := th.service.ValidateToken(r.Context(), requestBody)
	if err != nil {
		th.logger.Error("invalid", "error", err.Error())
		status = http.StatusUnauthorized
	}

	writeJSON(w, status, service.ValidationResponse{IsValid: validResult.IsValid})
}

func (th *TokenHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	var requestBody service.RevokeTokenRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		errString := "incorrect request body"
		th.logger.Error(errString, "error", err.Error())
		writeJSON(w, http.StatusBadRequest, apiError{errString})
		return
	}

	err := th.service.RevokeToken(r.Context(), requestBody)
	if err != nil {
		th.logger.Error("failed to revoke", "error", err.Error())
		writeJSON(w, http.StatusInternalServerError, apiError{"failed to revoke"})
		return
	}

	th.logger.Info("successfully executed")
	w.WriteHeader(http.StatusOK)
}
