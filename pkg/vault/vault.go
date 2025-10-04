package vault

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Vault struct {
	issuer    string
	secretKey []byte
}

func NewVault(issuer string, secretKey []byte) *Vault {
	return &Vault{
		issuer:    issuer,
		secretKey: secretKey,
	}
}

func (v *Vault) Encrypt(subject string, now time.Time, duration time.Duration) string {
	expiresAt := now.Add(duration)
	claim := &jwt.RegisteredClaims{
		Issuer:    v.issuer,
		Subject:   subject,
		Audience:  jwt.ClaimStrings{},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        "",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	token, err := accessToken.SignedString(v.secretKey)
	if err != nil {
		panic(err) // 어떻게 발생하지? 가능한가?
	}

	return token
}

func (v *Vault) Decrypt(token string, now time.Time) (string, error) {
	var claims jwt.RegisteredClaims

	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %w", ErrInvalidToken)
		}

		return v.secretKey, nil
	}, jwt.WithTimeFunc(func() time.Time { return now }))
	if err != nil {
		return "", errors.Join(ErrInvalidToken, err)
	}

	if claims.Issuer != v.issuer {
		return "", fmt.Errorf("invalid issuer: %w", ErrInvalidToken)
	}

	return claims.Subject, nil
}
