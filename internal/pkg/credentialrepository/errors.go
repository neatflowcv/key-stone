package credentialrepository

import "errors"

var (
	ErrCredentialNotFound      = errors.New("credential not found")
	ErrCredentialAlreadyExists = errors.New("credential already exists")
)
