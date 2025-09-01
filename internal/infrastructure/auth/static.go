package auth

type StaticAuthenticator struct {
	Token string
}

func NewStaticAuthenticator(token string) *StaticAuthenticator {
	return &StaticAuthenticator{Token: token}
}

func (a *StaticAuthenticator) Authenticate(header string) bool {
	expected := "Bearer " + a.Token
	return header == expected
}
