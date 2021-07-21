package model

type AuthData struct {
	Profile         map[string]string
	IsAuthenticated bool
	BaseUrl         string
	RedirectUri     string
	ClientId        string
	Issuer          string
	State           string
	Nonce           string
}
