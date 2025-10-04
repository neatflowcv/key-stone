package fake

import (
	"context"

	"github.com/neatflowcv/key-stone/internal/pkg/credentialrepository"
	"github.com/neatflowcv/key-stone/internal/pkg/domain"
)

var _ credentialrepository.Repository = (*Repository)(nil)

type Repository struct {
	credentials map[string]*domain.Credential
}

func NewRepository() *Repository {
	return &Repository{
		credentials: make(map[string]*domain.Credential),
	}
}

func (r *Repository) CreateCredential(ctx context.Context, credential *domain.Credential) error {
	if _, ok := r.credentials[credential.Username()]; ok {
		return credentialrepository.ErrCredentialAlreadyExists
	}

	r.credentials[credential.Username()] = credential

	return nil
}

func (r *Repository) DeleteCredential(ctx context.Context, credential *domain.Credential) error {
	if _, ok := r.credentials[credential.Username()]; !ok {
		return credentialrepository.ErrCredentialNotFound
	}

	delete(r.credentials, credential.Username())

	return nil
}

func (r *Repository) GetCredential(ctx context.Context, username string) (*domain.Credential, error) {
	if _, ok := r.credentials[username]; !ok {
		return nil, credentialrepository.ErrCredentialNotFound
	}

	return r.credentials[username], nil
}
