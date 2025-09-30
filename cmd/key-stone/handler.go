package main

import (
	"context"
	"log"

	"github.com/neatflowcv/key-stone/gen/token"
	"github.com/neatflowcv/key-stone/gen/user"
)

var _ user.Service = (*Handler)(nil)
var _ token.Service = (*Handler)(nil)

type Handler struct {
	users map[string]user.UserInput
}

func NewHandler() *Handler {
	return &Handler{
		users: make(map[string]user.UserInput),
	}
}

func (h *Handler) Issue(ctx context.Context, payload *token.IssuePayload) (*token.TokenDetail, error) {
	panic("unimplemented")
}

func (h *Handler) Refresh(ctx context.Context, payload *token.RefreshPayload) (*token.TokenDetail, error) {
	panic("unimplemented")
}

func (h *Handler) Create(ctx context.Context, payload *user.CreatePayload) error {
	log.Printf("payload %+v", payload)
	h.users[payload.User.Name] = *payload.User

	return nil
}

func (h *Handler) Delete(ctx context.Context, payload *user.DeleteUserPayload) error {
	panic("unimplemented")
}
