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

func (v *Vault) Encrypt(issuedAt time.Time, duration time.Duration, subject string) string {
	expiresAt := issuedAt.Add(duration)
	claim := &jwt.RegisteredClaims{
		Issuer:    v.issuer,
		Subject:   subject,
		Audience:  jwt.ClaimStrings{},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(issuedAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ID:        "",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	token, err := accessToken.SignedString(v.secretKey)
	if err != nil {
		panic(err) // 어떻게 발생하지? 가능한가?
	}

	return token
}

var (
	ErrInvalidMethod = errors.New("unexpected signing method")
)

func (v *Vault) Decrypt(now time.Time, encryptedValue string) (string, error) {
	var claims jwt.RegisteredClaims

	_, err := jwt.ParseWithClaims(encryptedValue, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidMethod
		}

		return v.secretKey, nil
	}, jwt.WithTimeFunc(func() time.Time { return now }))
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return "", ErrInvalidMethod
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return "", ErrInvalidMethod
		case errors.Is(err, jwt.ErrTokenExpired):
			return "", ErrInvalidMethod
		case errors.Is(err, ErrInvalidMethod):
			return "", ErrInvalidMethod
		default:
			return "", fmt.Errorf("failed to parse token: %w", err)
		}
	}

	if claims.Issuer != v.issuer {
		return "", ErrInvalidMethod
	}

	return claims.Subject, nil
}
