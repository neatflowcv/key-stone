package main

import (
	"context"
	"errors"
	"strings"

	"github.com/neatflowcv/key-stone/gen/token"
	"github.com/neatflowcv/key-stone/gen/user"
	"github.com/neatflowcv/key-stone/internal/app/flow"
)

var _ user.Service = (*Handler)(nil)
var _ token.Service = (*Handler)(nil)

type Handler struct {
	service *flow.Service
}

func NewHandler(
	service *flow.Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Issue(ctx context.Context, payload *token.IssuePayload) (*token.TokenDetail, error) {
	tokenSet, err := h.service.CreateToken(ctx, &flow.Credential{
		Username: payload.User.Username,
		Password: payload.User.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, flow.ErrUserNotFound):
			return nil, token.MakeUnauthorized(err)
		case errors.Is(err, flow.ErrUserUnauthorized):
			return nil, token.MakeUnauthorized(err)
		default:
			return nil, token.MakeInternalServerError(err)
		}
	}

	return &token.TokenDetail{
		AccessToken:  tokenSet.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenSet.ExpiresIn,
		RefreshToken: tokenSet.RefreshToken,
	}, nil
}

func (h *Handler) Refresh(ctx context.Context, payload *token.RefreshPayload) (*token.TokenDetail, error) {
	tokenSet, err := h.service.RefreshToken(ctx, &flow.TokenSetInput{
		AccessToken:  payload.Token.AccessToken,
		RefreshToken: payload.Token.RefreshToken,
	})
	if err != nil {
		switch {
		case errors.Is(err, flow.ErrTokenInvalid):
			return nil, token.MakeUnauthorized(err)
		case errors.Is(err, flow.ErrUserNotFound):
			return nil, token.MakeUnauthorized(err)
		default:
			return nil, token.MakeInternalServerError(err)
		}
	}

	return &token.TokenDetail{
		AccessToken:  tokenSet.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenSet.ExpiresIn,
		RefreshToken: tokenSet.RefreshToken,
	}, nil
}

func (h *Handler) Create(ctx context.Context, payload *user.CreatePayload) error {
	err := h.service.CreateUser(ctx, &flow.Credential{
		Username: payload.User.Username,
		Password: payload.User.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, flow.ErrUserAlreadyExists):
			return user.MakeUserAlreadyExists(err)
		default:
			return user.MakeInternalServerError(err)
		}
	}

	return nil
}

func (h *Handler) Delete(ctx context.Context, payload *user.DeleteUserPayload) error {
	token := strings.TrimPrefix(payload.Authorization, "Bearer ")

	err := h.service.DeleteUser(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, flow.ErrTokenInvalid),
			errors.Is(err, flow.ErrUserNotFound):
			return user.MakeUnauthorized(err)
		default:
			return user.MakeInternalServerError(err)
		}
	}

	return nil
}
