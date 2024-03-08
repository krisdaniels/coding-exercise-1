package app_handlers

import (
	"coding_exercise/internal/lib"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AuthHandler_Handle_returns_error_on_empty_username(t *testing.T) {
	// Arrange
	oidc_provider_mock := lib.OidcProviderMock{}
	sut := NewAuthHandler(&oidc_provider_mock)
	req := AuthRequest{
		Username: "",
		Password: "some-password",
	}

	// Act
	res, err := sut.Handle(req)

	// Assert
	assert.Nil(t, res)
	require.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrAuthValidationError))
}

func Test_AuthHandler_Handle_doesnt_allow_spaces_as_username(t *testing.T) {
	// Arrange
	oidc_provider_mock := lib.OidcProviderMock{}
	sut := NewAuthHandler(&oidc_provider_mock)
	req := AuthRequest{
		Username: "  ",
		Password: "some-password  ",
	}

	// Act
	res, err := sut.Handle(req)

	// Assert
	assert.Nil(t, res)
	require.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrAuthValidationError))
}

func Test_AuthHandler_Handle_returns_error_on_empty_password(t *testing.T) {
	// Arrange
	oidc_provider_mock := lib.OidcProviderMock{}
	sut := NewAuthHandler(&oidc_provider_mock)
	req := AuthRequest{
		Username: "some-username",
		Password: "",
	}

	// Act
	res, err := sut.Handle(req)

	// Assert
	assert.Nil(t, res)
	require.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrAuthValidationError))
}

func Test_AuthHandler_Handle_allows_spaces_as_password(t *testing.T) {
	// Arrange
	oidc_provider_mock := lib.OidcProviderMock{}
	sut := NewAuthHandler(&oidc_provider_mock)
	req := AuthRequest{
		Username: "some-username",
		Password: "    ",
	}

	// Act
	res, err := sut.Handle(req)

	// Assert
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func Test_AuthHandler_Handle_calls_oidc_provider_with_username(t *testing.T) {
	// Arrange
	oidc_provider_mock := lib.OidcProviderMock{}
	sut := NewAuthHandler(&oidc_provider_mock)
	req := AuthRequest{
		Username: "some-username",
		Password: "some-password",
	}

	// Act
	_, err := sut.Handle(req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "some-username", oidc_provider_mock.LastUsername)
}

func Test_AuthHandler_Handle_returns_result_from_oidc_provider(t *testing.T) {
	// Arrange
	oidc_provider_mock := lib.OidcProviderMock{
		NextGenerateTokenResult: "some-token",
		NextGenerateTokenError:  nil,
	}
	sut := NewAuthHandler(&oidc_provider_mock)
	req := AuthRequest{
		Username: "some-username",
		Password: "some-password",
	}

	// Act
	res, err := sut.Handle(req)

	// Assert
	assert.Nil(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "some-token", res.Token)
}

func Test_AuthHandler_Handle_returns_error_when_token_generation_fails(t *testing.T) {
	// Arrange
	oidc_provider_mock := lib.OidcProviderMock{
		NextGenerateTokenError: errors.New("some-error"),
	}
	sut := NewAuthHandler(&oidc_provider_mock)
	req := AuthRequest{
		Username: "some-username",
		Password: "some-password",
	}

	// Act
	res, err := sut.Handle(req)

	// Assert
	require.NotNil(t, err)
	assert.Nil(t, res)
	assert.True(t, errors.Is(err, ErrAuthTokenGenerationError))
}
