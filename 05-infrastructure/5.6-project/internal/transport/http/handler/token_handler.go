package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/course/token-service/internal/service"
	"github.com/course/token-service/internal/transport/util"
	"github.com/course/token-service/internal/transport/util/auth"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("/http/handler/token-handler")

type userCredentialsReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userCredStruct struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email,omitempty"`
}

type tokenStringReq struct {
	Token string `json:"token"`
}

type tokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenHandler struct {
	service *service.TokenService
	logger  *slog.Logger
}

func NewTokenHandler(ts *service.TokenService, logger *slog.Logger) *TokenHandler {
	return &TokenHandler{service: ts, logger: logger}
}

func (th *TokenHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "TokenHandler.Register")
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.path", r.URL.Path),
	)

	var requestBody userCredentialsReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		errString := "incorrect request body"
		th.logger.ErrorContext(ctx, errString, "error", err.Error())
		util.WriteError(w, http.StatusBadRequest, errString)
		return
	}

	userID, err := th.service.Register(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		if errors.Is(err, service.ErrBadRequest) {
			th.logger.ErrorContext(ctx, "bad request", "error", err.Error())
			util.WriteError(w, http.StatusBadRequest, "bad request")
			return
		}
		th.logger.ErrorContext(ctx, "internal error", "error", err.Error())
		util.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	th.logger.InfoContext(ctx, "register req is successful")
	util.WriteJSON(w, http.StatusCreated, userCredStruct{UserID: userID})
}

func (th *TokenHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "TokenHandler.Login")
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.path", r.URL.Path),
	)

	var requestBody userCredentialsReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		errString := "incorrect request body"
		th.logger.ErrorContext(ctx, errString, "error", err.Error())
		util.WriteError(w, http.StatusBadRequest, errString)
		return
	}

	refToken, accToken, err := th.service.Login(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		if errors.Is(err, service.ErrBadRequest) {
			th.logger.ErrorContext(ctx, "bad request", "error", err.Error())
			util.WriteError(w, http.StatusBadRequest, "bad request")
			return
		} else if errors.Is(err, service.ErrUnauthorized) {
			th.logger.ErrorContext(ctx, "invalid credentials", "error", err.Error())
			util.WriteError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		th.logger.ErrorContext(ctx, "internal error", "error", err.Error())
		util.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	th.logger.InfoContext(ctx, "login req is successful")
	util.WriteJSON(w, http.StatusOK, tokensResponse{AccessToken: accToken, RefreshToken: refToken})
}

func (th *TokenHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "TokenHandler.Refresh")
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.path", r.URL.Path),
	)

	var requestBody tokenStringReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		errString := "incorrect request body"
		th.logger.ErrorContext(ctx, errString, "error", err.Error())
		util.WriteError(w, http.StatusBadRequest, errString)
		return
	}

	refToken, accToken, err := th.service.Refresh(ctx, requestBody.Token)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		if errors.Is(err, service.ErrUnauthorized) {
			th.logger.ErrorContext(ctx, "invalid credentials", "error", err.Error())
			util.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		th.logger.ErrorContext(ctx, "internal error", "error", err.Error())
		util.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	th.logger.InfoContext(ctx, "refresh req is successful")
	util.WriteJSON(w, http.StatusOK, tokensResponse{AccessToken: accToken, RefreshToken: refToken})
}

func (th *TokenHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "TokenHandler.Logout")
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.path", r.URL.Path),
	)

	var requestBody tokenStringReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		errString := "incorrect request body"
		th.logger.ErrorContext(ctx, errString, "error", err.Error())
		util.WriteError(w, http.StatusBadRequest, errString)
		return
	}

	err := th.service.Revoke(ctx, requestBody.Token)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		if errors.Is(err, service.ErrUnauthorized) {
			th.logger.ErrorContext(ctx, "invalid credentials", "error", err.Error())
			util.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		th.logger.ErrorContext(ctx, "internal error", "error", err.Error())
		util.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := struct {
		Revoked bool `json:"revoked"`
	}{Revoked: true}

	th.logger.InfoContext(ctx, "logout req is successful")
	util.WriteJSON(w, http.StatusOK, resp)
}

func (th *TokenHandler) Validate(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "TokenHandler.Validate")
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.path", r.URL.Path),
	)

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		err := errors.New("no token provided")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		errString := "incorrect request headers"
		th.logger.ErrorContext(ctx, errString, "error", err.Error())
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	email, err := th.service.GetUser(ctx, userID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		th.logger.ErrorContext(ctx, "user not found", "error", err.Error())
		util.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	th.logger.InfoContext(ctx, "validate req is successful")
	util.WriteJSON(w, http.StatusOK, userCredStruct{UserID: userID, Email: email})
}
