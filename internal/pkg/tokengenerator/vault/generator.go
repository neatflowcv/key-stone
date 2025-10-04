package vault

import (
	"errors"
	"time"

	"github.com/neatflowcv/key-stone/internal/pkg/tokengenerator"
	"github.com/neatflowcv/key-stone/pkg/vault"
)

var _ tokengenerator.Generator = (*Generator)(nil)

type Generator struct {
	vault *vault.Vault
}

func NewGenerator(issuer string, secretKey []byte) *Generator {
	return &Generator{vault: vault.NewVault(issuer, secretKey)}
}

func (g *Generator) GenerateToken(subject string, now time.Time, duration time.Duration) string {
	return g.vault.Encrypt(subject, now, duration)
}

func (g *Generator) ParseToken(encrypted string, now time.Time) (string, error) {
	ret, err := g.vault.Decrypt(encrypted, now)
	if err != nil {
		return "", errors.Join(err, tokengenerator.ErrTokenInvalid)
	}

	return ret, nil
}
