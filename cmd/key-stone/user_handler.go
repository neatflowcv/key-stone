package main

import (
	"context"
	"errors"
	"strings"

	"github.com/neatflowcv/key-stone/gen/user"
	"github.com/neatflowcv/key-stone/internal/app/flow"
)

var _ user.Service = (*UserHandler)(nil)

type UserHandler struct {
	service *flow.Service
}

func NewUserHandler(
	service *flow.Service,
) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Create(ctx context.Context, payload *user.CreatePayload) error {
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

func (h *UserHandler) Delete(ctx context.Context, payload *user.DeleteUserPayload) error {
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
