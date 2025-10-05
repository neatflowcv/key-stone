package bcrypt

import (
	"errors"
	"fmt"

	"github.com/neatflowcv/key-stone/internal/pkg/hasher"
	"golang.org/x/crypto/bcrypt"
)

var _ hasher.Hasher = (*Hasher)(nil)

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func (h *Hasher) Compare(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return hasher.ErrMismatched
		}

		return fmt.Errorf("failed to compare hash and password: %w", err)
	}

	return nil
}
