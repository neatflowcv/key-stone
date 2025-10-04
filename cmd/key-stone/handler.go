package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/neatflowcv/key-stone/gen/token"
	"github.com/neatflowcv/key-stone/gen/user"
	"github.com/neatflowcv/key-stone/pkg/vault"
)

var _ user.Service = (*Handler)(nil)
var _ token.Service = (*Handler)(nil)

type Handler struct {
	users    map[string]user.UserInput
	pubVault *vault.Vault
	priVault *vault.Vault
}

func NewHandler() *Handler {
	return &Handler{
		users:    make(map[string]user.UserInput),
		pubVault: vault.NewVault("key-stone", []byte("public-key")),
		priVault: vault.NewVault("key-stone", []byte("private-key")),
	}
}

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func (h *Handler) Issue(ctx context.Context, payload *token.IssuePayload) (*token.TokenDetail, error) {
	user, ok := h.users[payload.User.Username]
	if !ok {
		return nil, token.MakeUnauthorized(ErrUnauthorized)
	}

	if user.Password != payload.User.Password {
		return nil, token.MakeUnauthorized(ErrUnauthorized)
	}

	now := time.Now()

	return h.generate(now, user.Username)
}

func (h *Handler) Refresh(ctx context.Context, payload *token.RefreshPayload) (*token.TokenDetail, error) {
	subject, err := h.getSubject(payload)
	if err != nil {
		return nil, err
	}

	user, ok := h.users[subject]
	if !ok {
		return nil, token.MakeUnauthorized(ErrUnauthorized)
	}

	now := time.Now()

	return h.generate(now, user.Username)
}

func (h *Handler) Create(ctx context.Context, payload *user.CreatePayload) error {
	log.Printf("payload %+v", payload)
	h.users[payload.User.Username] = *payload.User

	return nil
}

func (h *Handler) Delete(ctx context.Context, payload *user.DeleteUserPayload) error {
	now := time.Now()
	token := strings.TrimPrefix(payload.Authorization, "Bearer ")

	subject, err := h.pubVault.Decrypt(token, now)
	if err == nil {
		return nil
	}

	delete(h.users, subject)

	return nil
}

func (h *Handler) getSubject(payload *token.RefreshPayload) (string, error) {
	now := time.Now()

	subject, err := h.pubVault.Decrypt(payload.Token.AccessToken, now)
	if err == nil {
		return subject, nil
	}

	subject, err = h.priVault.Decrypt(payload.Token.RefreshToken, now)
	if err == nil {
		return subject, nil
	}

	return "", token.MakeUnauthorized(ErrUnauthorized)
}

func (h *Handler) generate(now time.Time, subject string) (*token.TokenDetail, error) {
	const (
		accessTokenDuration  = time.Minute * 15
		refreshTokenDuration = time.Hour * 24 * 14
	)

	accessToken := h.pubVault.Encrypt(subject, now, accessTokenDuration)
	refreshToken := h.priVault.Encrypt(subject, now, refreshTokenDuration)
	tokenType := "Bearer"
	expiresIn := int(accessTokenDuration.Seconds())

	return &token.TokenDetail{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
	}, nil
}
