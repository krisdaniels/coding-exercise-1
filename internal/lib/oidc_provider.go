package lib

type OidcProvider interface {
	GenerateToken(username string) (string, error)
	ValidateToken(token string) (map[string]interface{}, error)
}
