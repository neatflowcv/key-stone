package credentialrepository

import (
	"context"

	"github.com/neatflowcv/key-stone/internal/pkg/domain"
)

type Repository interface {
	CreateCredential(ctx context.Context, credential *domain.Credential) error
	DeleteCredential(ctx context.Context, credential *domain.Credential) error
	GetCredential(ctx context.Context, username string) (*domain.Credential, error)
}
