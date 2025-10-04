package flow

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserUnauthorized  = errors.New("user unauthorized")

	ErrTokenInvalid = errors.New("token is invalid")
)
