package handlers

import (
	"complex-task/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("complex-task/handler")

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
	ctx, span := tracer.Start(r.Context(), "RefreshHandler")
	defer span.End()

	var requestBody service.UserData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		errString := "incorrect request body"
		th.logger.ErrorContext(ctx, errString, "error", err.Error())
		writeJSON(w, http.StatusBadRequest, apiError{errString})
		return
	}

	tokens, err := th.service.IssueToken(ctx, requestBody)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		if errors.Is(err, service.LoginError) {
			th.logger.ErrorContext(ctx, "incorrect login or password", "error", err.Error())
			writeJSON(w, http.StatusUnauthorized, apiError{"incorrect login or password"})
			return
		}
		th.logger.ErrorContext(ctx, "internal error", "error", err.Error())
		writeJSON(w, http.StatusInternalServerError, apiError{"internal error"})
		return
	}

	th.logger.InfoContext(ctx, "successfully executed")
	writeJSON(w, http.StatusCreated, service.PairOfTokens{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
}

func (th *TokenHandler) Validate(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "ValidateHandler")
	defer span.End()

	var requestBody service.ValidationRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		errString := "incorrect request body"
		th.logger.ErrorContext(ctx, errString, "error", err.Error())
		writeJSON(w, http.StatusBadRequest, apiError{errString})
		return
	}

	status := http.StatusOK
	validResult, err := th.service.ValidateToken(r.Context(), requestBody)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		th.logger.ErrorContext(ctx, "invalid", "error", err.Error())
		status = http.StatusUnauthorized
	}

	writeJSON(w, status, service.ValidationResponse{IsValid: validResult.IsValid})
}

func (th *TokenHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "RevokeHandler")
	defer span.End()

	var requestBody service.RevokeTokenRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		errString := "incorrect request body"
		th.logger.ErrorContext(ctx, errString, "error", err.Error())
		writeJSON(w, http.StatusBadRequest, apiError{errString})
		return
	}

	err := th.service.RevokeToken(r.Context(), requestBody)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		th.logger.ErrorContext(ctx, "failed to revoke", "error", err.Error())
		writeJSON(w, http.StatusInternalServerError, apiError{"failed to revoke"})
		return
	}

	th.logger.InfoContext(ctx, "successfully executed")
	w.WriteHeader(http.StatusOK)
}
