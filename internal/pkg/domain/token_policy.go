package domain

import "time"

type TokenPolicy struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewTokenPolicy() *TokenPolicy {
	return &TokenPolicy{
		accessTokenDuration:  15 * time.Minute, //nolint:mnd
		refreshTokenDuration: 14 * 24 * time.Hour,
	}
}

func (t *TokenPolicy) AccessTokenDuration() time.Duration {
	return t.accessTokenDuration
}

func (t *TokenPolicy) RefreshTokenDuration() time.Duration {
	return t.refreshTokenDuration
}
