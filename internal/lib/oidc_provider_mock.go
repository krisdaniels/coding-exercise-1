package lib

type OidcProviderMock struct {
	GenerateTokenCalled bool
	ValidateTokenCalled bool

	LastUsername string
	LastToken    string

	NextGenerateTokenResult string
	NextGenerateTokenError  error

	NextValidateTokenResult map[string]interface{}
	NextValidateTokenError  error
}

func (m *OidcProviderMock) GenerateToken(username string) (string, error) {
	m.GenerateTokenCalled = true
	m.LastUsername = username
	return m.NextGenerateTokenResult, m.NextGenerateTokenError
}

func (m *OidcProviderMock) ValidateToken(token string) (map[string]interface{}, error) {
	m.ValidateTokenCalled = true
	m.LastToken = token
	return m.NextValidateTokenResult, m.NextValidateTokenError
}
