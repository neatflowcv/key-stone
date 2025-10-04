package flow

type TokenSetInput struct {
	AccessToken  string
	RefreshToken string
}

type TokenSetOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}
