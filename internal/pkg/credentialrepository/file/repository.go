package file

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/neatflowcv/key-stone/internal/pkg/credentialrepository"
	"github.com/neatflowcv/key-stone/internal/pkg/domain"
)

var _ credentialrepository.Repository = (*Repository)(nil)

type Repository struct {
	path string
}

func NewRepository(path string) (*Repository, error) {
	const createPerm = 0750

	err := os.MkdirAll(path, createPerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return &Repository{path: path}, nil
}

func (r *Repository) CreateCredential(ctx context.Context, credential *domain.Credential) error {
	const openPerm = 0600

	filePath := filepath.Join(r.path, credential.Username())

	file, err := os.OpenFile(filepath.Clean(filePath), os.O_CREATE|os.O_WRONLY|os.O_EXCL, openPerm)
	if err != nil {
		if os.IsExist(err) {
			return credentialrepository.ErrCredentialAlreadyExists
		}

		return fmt.Errorf("failed to create credential: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	_, err = file.WriteString(credential.Password())
	if err != nil {
		return fmt.Errorf("failed to write credential: %w", err)
	}

	return nil
}

func (r *Repository) DeleteCredential(ctx context.Context, credential *domain.Credential) error {
	filePath := filepath.Join(r.path, credential.Username())

	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return credentialrepository.ErrCredentialNotFound
		}

		return fmt.Errorf("failed to delete credential: %w", err)
	}

	return nil
}

func (r *Repository) GetCredential(ctx context.Context, username string) (*domain.Credential, error) {
	filePath := filepath.Join(r.path, username)

	password, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, credentialrepository.ErrCredentialNotFound
		}

		return nil, fmt.Errorf("failed to read credential: %w", err)
	}

	return domain.NewCredential(username, string(password)), nil
}
