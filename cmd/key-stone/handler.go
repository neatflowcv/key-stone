package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/neatflowcv/key-stone/gen/token"
	"github.com/neatflowcv/key-stone/gen/user"
	"github.com/neatflowcv/key-stone/internal/pkg/credentialrepository"
	"github.com/neatflowcv/key-stone/internal/pkg/domain"
	"github.com/neatflowcv/key-stone/internal/pkg/tokengenerator"
)

var _ user.Service = (*Handler)(nil)
var _ token.Service = (*Handler)(nil)

type Handler struct {
	repo   credentialrepository.Repository
	pubGen tokengenerator.Generator
	priGen tokengenerator.Generator
}

func NewHandler(
	repo credentialrepository.Repository,
	pubGen tokengenerator.Generator,
	priGen tokengenerator.Generator,
) *Handler {
	return &Handler{
		repo:   repo,
		pubGen: pubGen,
		priGen: priGen,
	}
}

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func (h *Handler) Issue(ctx context.Context, payload *token.IssuePayload) (*token.TokenDetail, error) {
	cred, err := h.repo.GetCredential(ctx, payload.User.Username)
	if err != nil {
		if errors.Is(err, credentialrepository.ErrCredentialNotFound) {
			return nil, token.MakeUnauthorized(ErrUnauthorized)
		}

		return nil, token.MakeInternalServerError(err)
	}

	if cred.Password() != payload.User.Password {
		return nil, token.MakeUnauthorized(ErrUnauthorized)
	}

	now := time.Now()

	return h.generate(now, cred.Username())
}

func (h *Handler) Refresh(ctx context.Context, payload *token.RefreshPayload) (*token.TokenDetail, error) {
	subject, err := h.getSubject(payload)
	if err != nil {
		return nil, err
	}

	user, err := h.repo.GetCredential(ctx, subject)
	if err != nil {
		if errors.Is(err, credentialrepository.ErrCredentialNotFound) {
			return nil, token.MakeUnauthorized(ErrUnauthorized)
		}

		return nil, token.MakeInternalServerError(err)
	}

	now := time.Now()

	return h.generate(now, user.Username())
}

func (h *Handler) Create(ctx context.Context, payload *user.CreatePayload) error {
	cred := domain.NewCredential(payload.User.Username, payload.User.Password)

	err := h.repo.CreateCredential(ctx, cred)
	if err != nil {
		return user.MakeInternalServerError(err)
	}

	return nil
}

func (h *Handler) Delete(ctx context.Context, payload *user.DeleteUserPayload) error {
	now := time.Now()
	token := strings.TrimPrefix(payload.Authorization, "Bearer ")

	subject, err := h.pubGen.ParseToken(token, now)
	if err == nil {
		return nil
	}

	cred, err := h.repo.GetCredential(ctx, subject)
	if err != nil {
		if errors.Is(err, credentialrepository.ErrCredentialNotFound) {
			return user.MakeUnauthorized(ErrUnauthorized)
		}

		return user.MakeInternalServerError(err)
	}

	err = h.repo.DeleteCredential(ctx, cred)
	if err != nil {
		return user.MakeUnauthorized(err)
	}

	return nil
}

func (h *Handler) getSubject(payload *token.RefreshPayload) (string, error) {
	now := time.Now()

	subject, err := h.pubGen.ParseToken(payload.Token.AccessToken, now)
	if err == nil {
		return subject, nil
	}

	subject, err = h.priGen.ParseToken(payload.Token.RefreshToken, now)
	if err == nil {
		return subject, nil
	}

	return "", token.MakeUnauthorized(ErrUnauthorized)
}

func (h *Handler) generate(now time.Time, subject string) (*token.TokenDetail, error) {
	policy := domain.NewTokenPolicy()

	accessToken := h.pubGen.GenerateToken(subject, now, policy.AccessTokenDuration())
	refreshToken := h.priGen.GenerateToken(subject, now, policy.RefreshTokenDuration())
	tokenType := "Bearer"
	expiresIn := int(policy.AccessTokenDuration().Seconds())

	return &token.TokenDetail{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
	}, nil
}
