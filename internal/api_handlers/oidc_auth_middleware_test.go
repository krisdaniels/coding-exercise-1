package api_handlers

import (
	"coding_exercise/internal/lib"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OidcAuthMiddleware_returns_401_when_missing_authorization_header(t *testing.T) {
	// Arrange
	called_next := false
	next_func := func(w http.ResponseWriter, r *http.Request) {
		called_next = true
	}

	oidc_provider_mock := &lib.OidcProviderMock{}
	middleware := NewOidcAuthMiddleware(oidc_provider_mock)
	handler := http.HandlerFunc(next_func)

	sut := middleware.GetHandler(handler)
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	// Act
	sut.ServeHTTP(recorder, req)

	// Assert
	assert.False(t, called_next)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func Test_OidcAuthMiddleware_returns_401_when_missing_bearer_in_authorization_header(t *testing.T) {
	called_next := false
	next_func := func(w http.ResponseWriter, r *http.Request) {
		called_next = true
	}

	oidc_provider_mock := &lib.OidcProviderMock{}
	middleware := NewOidcAuthMiddleware(oidc_provider_mock)
	handler := http.HandlerFunc(next_func)

	sut := middleware.GetHandler(handler)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "something invalid")
	recorder := httptest.NewRecorder()

	// Act
	sut.ServeHTTP(recorder, req)

	// Assert
	assert.False(t, called_next)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func Test_OidcAuthMiddleware_returns_401_when_missing_token_authorization_header(t *testing.T) {
	called_next := false
	next_func := func(w http.ResponseWriter, r *http.Request) {
		called_next = true
	}

	oidc_provider_mock := &lib.OidcProviderMock{}
	middleware := NewOidcAuthMiddleware(oidc_provider_mock)
	handler := http.HandlerFunc(next_func)

	sut := middleware.GetHandler(handler)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer ")
	recorder := httptest.NewRecorder()

	// Act
	sut.ServeHTTP(recorder, req)

	// Assert
	assert.False(t, called_next)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func Test_OidcAuthMiddleware_returns_401_when_invalid_token_in_authorization_header(t *testing.T) {
	called_next := false
	next_func := func(w http.ResponseWriter, r *http.Request) {
		called_next = true
	}

	oidc_provider_mock := &lib.OidcProviderMock{
		NextValidateTokenError: errors.New("some error"),
	}
	middleware := NewOidcAuthMiddleware(oidc_provider_mock)
	handler := http.HandlerFunc(next_func)

	sut := middleware.GetHandler(handler)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer invalid-token")
	recorder := httptest.NewRecorder()

	// Act
	sut.ServeHTTP(recorder, req)

	// Assert
	assert.False(t, called_next)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func Test_OidcAuthMiddleware_returns_calls_next_when_valid_token_in_authorization_header(t *testing.T) {
	called_next := false
	next_func := func(w http.ResponseWriter, r *http.Request) {
		called_next = true
	}

	oidc_provider_mock := &lib.OidcProviderMock{
		NextValidateTokenError: nil,
	}
	middleware := NewOidcAuthMiddleware(oidc_provider_mock)
	handler := http.HandlerFunc(next_func)

	sut := middleware.GetHandler(handler)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer valid-token")
	recorder := httptest.NewRecorder()

	// Act
	sut.ServeHTTP(recorder, req)

	// Assert
	assert.True(t, called_next)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
