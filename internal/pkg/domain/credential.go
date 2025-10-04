package domain

type Credential struct {
	username string
	password string
}

func NewCredential(username, password string) *Credential {
	return &Credential{
		username: username,
		password: password,
	}
}

func (c *Credential) Username() string {
	return c.username
}

func (c *Credential) Password() string {
	return c.password
}
