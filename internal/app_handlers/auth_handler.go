package app_handlers

import (
	"coding_exercise/internal/lib"
	"errors"
	"log"
	"strings"
)

var (
	ErrAuthValidationError      = errors.New("username or password is empty")
	ErrAuthTokenGenerationError = errors.New("error generating token")
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	oidcProvider lib.OidcProvider
}

func NewAuthHandler(oidcProvider lib.OidcProvider) AppHandler[AuthRequest, AuthResponse] {
	return &AuthHandler{
		oidcProvider: oidcProvider,
	}
}

func (h *AuthHandler) Handle(request AuthRequest) (*AuthResponse, error) {
	request.Username = strings.TrimSpace(request.Username)

	if request.Username == "" || request.Password == "" {
		return nil, ErrAuthValidationError
	}

	// todo: validate credentials

	token, err := h.oidcProvider.GenerateToken(request.Username)
	if err != nil {
		log.Printf("error while generating token for %s: %s", request.Username, err)
		return nil, ErrAuthTokenGenerationError
	}

	response := &AuthResponse{
		Token: token,
	}

	return response, nil
}
