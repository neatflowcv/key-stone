package flow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neatflowcv/key-stone/internal/pkg/credentialrepository"
	"github.com/neatflowcv/key-stone/internal/pkg/domain"
	"github.com/neatflowcv/key-stone/internal/pkg/hasher"
	"github.com/neatflowcv/key-stone/internal/pkg/tokengenerator"
)

type Service struct {
	repo   credentialrepository.Repository
	hasher hasher.Hasher
	pubGen tokengenerator.Generator
	priGen tokengenerator.Generator
}

func NewService(
	repo credentialrepository.Repository,
	hasher hasher.Hasher,
	pubGen tokengenerator.Generator,
	priGen tokengenerator.Generator,
) *Service {
	return &Service{
		repo:   repo,
		hasher: hasher,
		pubGen: pubGen,
		priGen: priGen,
	}
}

// CreateUser creates a new user
// Returns:
//   - ErrUserAlreadyExists if the user already exists
func (s *Service) CreateUser(ctx context.Context, credential *Credential) error {
	hashedPassword, err := s.hasher.Hash(credential.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	cred := domain.NewCredential(credential.Username, hashedPassword)

	err = s.repo.CreateCredential(ctx, cred)
	if err != nil {
		return casError(err, credentialrepository.ErrCredentialAlreadyExists, ErrUserAlreadyExists)
	}

	return nil
}

// DeleteUser deletes the user
// Returns:
//   - ErrTokenInvalid if the token is invalid
//   - ErrUserNotFound if the user does not exist
func (s *Service) DeleteUser(ctx context.Context, token string) error {
	username, err := s.pubGen.ParseToken(token, time.Now())
	if err != nil {
		return casError(err, tokengenerator.ErrTokenInvalid, ErrTokenInvalid)
	}

	cred, err := s.repo.GetCredential(ctx, username)
	if err != nil {
		return casError(err, credentialrepository.ErrCredentialNotFound, ErrUserNotFound)
	}

	err = s.repo.DeleteCredential(ctx, cred)
	if err != nil {
		return casError(err, credentialrepository.ErrCredentialNotFound, ErrUserNotFound)
	}

	return nil
}

// CreateToken creates a new token
// Returns:
//   - ErrUserNotFound if the user does not exist
//   - ErrUserUnauthorized if the user is unauthorized
func (s *Service) CreateToken(ctx context.Context, credential *Credential) (*TokenSetOutput, error) {
	cred, err := s.repo.GetCredential(ctx, credential.Username)
	if err != nil {
		return nil, casError(err, credentialrepository.ErrCredentialNotFound, ErrUserNotFound)
	}

	err = s.hasher.Compare(credential.Password, cred.Password())
	if err != nil {
		return nil, casError(err, hasher.ErrMismatched, ErrUserUnauthorized)
	}

	return s.createTokenSet(cred.Username()), nil
}

// RefreshToken refreshes a token
// Returns:
//   - ErrTokenInvalid if the token is invalid
//   - ErrUserNotFound if the user does not exist
func (s *Service) RefreshToken(ctx context.Context, tokenSet *TokenSetInput) (*TokenSetOutput, error) {
	username, err := s.extractSubject(tokenSet)
	if err != nil {
		return nil, casError(err, ErrTokenInvalid, ErrTokenInvalid)
	}

	cred, err := s.repo.GetCredential(ctx, username)
	if err != nil {
		return nil, casError(err, credentialrepository.ErrCredentialNotFound, ErrUserNotFound)
	}

	return s.createTokenSet(cred.Username()), nil
}

func (s *Service) createTokenSet(username string) *TokenSetOutput {
	policy := domain.NewTokenPolicy()

	accessToken := s.pubGen.GenerateToken(username, time.Now(), policy.AccessTokenDuration())
	refreshToken := s.priGen.GenerateToken(username, time.Now(), policy.RefreshTokenDuration())
	expiresIn := int(policy.AccessTokenDuration().Seconds())

	return &TokenSetOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}

func (s *Service) extractSubject(tokenSet *TokenSetInput) (string, error) {
	subject, err := s.pubGen.ParseToken(tokenSet.AccessToken, time.Now())
	if err == nil {
		return subject, nil
	}

	subject, err = s.priGen.ParseToken(tokenSet.RefreshToken, time.Now())
	if err == nil {
		return subject, nil
	}

	return "", ErrTokenInvalid
}

func casError(err error, expected error, fresh error) error {
	if errors.Is(err, expected) {
		return fresh
	}

	return err
}
