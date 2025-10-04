package main

import (
	"context"
	"errors"

	"github.com/neatflowcv/key-stone/gen/token"
	"github.com/neatflowcv/key-stone/internal/app/flow"
)

var _ token.Service = (*TokenHandler)(nil)

type TokenHandler struct {
	service *flow.Service
}

func NewTokenHandler(
	service *flow.Service,
) *TokenHandler {
	return &TokenHandler{
		service: service,
	}
}

func (h *TokenHandler) Issue(ctx context.Context, payload *token.IssuePayload) (*token.TokenDetail, error) {
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

func (h *TokenHandler) Refresh(ctx context.Context, payload *token.RefreshPayload) (*token.TokenDetail, error) {
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
