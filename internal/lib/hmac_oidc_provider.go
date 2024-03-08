package lib

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrHmacOidcProviderValidationError = errors.New("username cannot be empty")
	ErrHmacOidcProviderInvalidToken    = errors.New("invalid token")
)

type HmacOidcProvider struct {
	secret []byte
	issuer string
	now    func() time.Time
}

func NewHmacOidcProvider(secret string, issuer string) OidcProvider {
	return &HmacOidcProvider{
		secret: []byte(secret),
		issuer: issuer,
		now:    time.Now,
	}
}

func (p *HmacOidcProvider) GenerateToken(username string) (string, error) {
	if username == "" {
		return "", ErrHmacOidcProviderValidationError
	}

	jwt.TimeFunc = p.now

	now := p.now().Unix()
	claims := &jwt.StandardClaims{
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: now + 3600,
		Issuer:    p.issuer,
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(p.secret)
}

func (p *HmacOidcProvider) ValidateToken(tokenString string) (map[string]interface{}, error) {
	jwt.TimeFunc = p.now

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}

		if token.Claims.(jwt.MapClaims)["iss"] != p.issuer {
			return nil, fmt.Errorf("unknown issuer: %v", token.Claims.(jwt.MapClaims)["iss"])
		}

		return p.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, ErrHmacOidcProviderInvalidToken
	}
}
