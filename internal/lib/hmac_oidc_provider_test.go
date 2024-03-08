package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HmacOidcProvider_GenerateToken_returns_error_on_empty_username(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("", "some-issuer")

	// Act
	_, err := sut.GenerateToken("")

	// Assert
	assert.NotNil(t, err)
}

func Test_HmacOidcProvider_GenerateToken_generates_a_valid_token(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("some-secret", "some-issuer")
	now, _ := time.Parse(time.RFC3339, "2000-01-02T03:04:05.00Z")
	sut.(*HmacOidcProvider).now = func() time.Time {
		return now
	}

	// Act
	token, err := sut.GenerateToken("some-user")

	// Assert
	assert.Nil(t, err)
	assert.Equal(
		t,
		"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0Njc4NTg0NSwiaWF0Ijo5NDY3ODIyNDUsImlzcyI6InNvbWUtaXNzdWVyIiwibmJmIjo5NDY3ODIyNDUsInN1YiI6InNvbWUtdXNlciJ9.d4elSuwdI1EZMdZG-CFgkKyPZkgxg1ZWpGbRBW6FQ7C0kBcyLwbZsZkRb2qMCesYmhIgsnjHgo3I6rY-EB8AYw",
		token,
	)
}

func Test_HmacOidcProvider_ValidateToken_returns_error_on_invalid_token(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("some-secret", "some-issuer")
	now, _ := time.Parse(time.RFC3339, "2000-01-02T03:04:05.00Z")
	sut.(*HmacOidcProvider).now = func() time.Time {
		return now
	}

	// replaced signature with invalid signature
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0Njc4NTg0NSwiaWF0Ijo5NDY3ODIyNDUsImlzcyI6InNvbWUtaXNzdWVyIiwibmJmIjo5NDY3ODIyNDUsInN1YiI6InNvbWUtdXNlciJ9.invalid-signature"

	// Act
	claims, err := sut.ValidateToken(token)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, claims)
}

func Test_HmacOidcProvider_ValidateToken_returns_error_on_invalid_signature_algorithm(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("some-secret", "some-issuer")
	now, _ := time.Parse(time.RFC3339, "2000-01-02T03:04:05.00Z")
	sut.(*HmacOidcProvider).now = func() time.Time {
		return now
	}

	// changed alg in header to RS256
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0Njc4NTg0NSwiaWF0Ijo5NDY3ODIyNDUsImlzcyI6InNvbWUtaXNzdWVyIiwibmJmIjo5NDY3ODIyNDUsInN1YiI6InNvbWUtdXNlciJ9.d4elSuwdI1EZMdZG-CFgkKyPZkgxg1ZWpGbRBW6FQ7C0kBcyLwbZsZkRb2qMCesYmhIgsnjHgo3I6rY-EB8AYw"

	// Act
	claims, err := sut.ValidateToken(token)

	// Assert
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid signing method")
	assert.Nil(t, claims)
}

func Test_HmacOidcProvider_ValidateToken_returns_error_on_invalid_issuer(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("some-secret", "another-issuer")
	now, _ := time.Parse(time.RFC3339, "2000-01-02T03:04:05.00Z")
	sut.(*HmacOidcProvider).now = func() time.Time {
		return now
	}

	// issuer in token some-issuer, validated against another-issuer
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0Njc4NTg0NSwiaWF0Ijo5NDY3ODIyNDUsImlzcyI6InNvbWUtaXNzdWVyIiwibmJmIjo5NDY3ODIyNDUsInN1YiI6InNvbWUtdXNlciJ9.d4elSuwdI1EZMdZG-CFgkKyPZkgxg1ZWpGbRBW6FQ7C0kBcyLwbZsZkRb2qMCesYmhIgsnjHgo3I6rY-EB8AYw"

	// Act
	claims, err := sut.ValidateToken(token)

	// Assert
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "unknown issuer")
	assert.Nil(t, claims)
}

func Test_HmacOidcProvider_ValidateToken_returns_error_on_expired_token(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("some-secret", "some-issuer")

	// added 10h to make sure token has expired on validaton
	now, _ := time.Parse(time.RFC3339, "2000-01-02T13:04:05.00Z")
	sut.(*HmacOidcProvider).now = func() time.Time {
		return now
	}

	// token generated for date 2000-01-02T03:04:05.00Z, verified against 10h later: 2000-01-02T13:04:05.00Z
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0Njc4NTg0NSwiaWF0Ijo5NDY3ODIyNDUsImlzcyI6InNvbWUtaXNzdWVyIiwibmJmIjo5NDY3ODIyNDUsInN1YiI6InNvbWUtdXNlciJ9.d4elSuwdI1EZMdZG-CFgkKyPZkgxg1ZWpGbRBW6FQ7C0kBcyLwbZsZkRb2qMCesYmhIgsnjHgo3I6rY-EB8AYw"

	// Act
	claims, err := sut.ValidateToken(token)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, claims)
}

func Test_HmacOidcProvider_ValidateToken_returns_claims_on_valid_token(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("some-secret", "some-issuer")
	now, _ := time.Parse(time.RFC3339, "2000-01-02T03:04:05.00Z")
	sut.(*HmacOidcProvider).now = func() time.Time {
		return now
	}

	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0Njc4NTg0NSwiaWF0Ijo5NDY3ODIyNDUsImlzcyI6InNvbWUtaXNzdWVyIiwibmJmIjo5NDY3ODIyNDUsInN1YiI6InNvbWUtdXNlciJ9.d4elSuwdI1EZMdZG-CFgkKyPZkgxg1ZWpGbRBW6FQ7C0kBcyLwbZsZkRb2qMCesYmhIgsnjHgo3I6rY-EB8AYw"

	// Act
	claims, err := sut.ValidateToken(token)

	// Assert
	require.Nil(t, err)
	assert.NotEmpty(t, claims)
}

func Test_HmacOidcProvider_GenerateToken_generates_a_valid_token_with_username_as_subject(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("some-secret", "some-issuer")

	// Act
	token, err := sut.GenerateToken("some-user")
	require.Nil(t, err)
	require.NotEmpty(t, token)

	claims, err := sut.ValidateToken(token)

	// Assert
	assert.Nil(t, err)
	require.NotEmpty(t, claims)
	require.Contains(t, claims, "sub")
	assert.Equal(t, claims["sub"], "some-user")
}

func Test_HmacOidcProvider_GenerateToken_generates_a_valid_token_with_specified_issuer(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("some-secret", "some-issuer")

	// Act
	token, err := sut.GenerateToken("some-user")
	require.Nil(t, err)
	require.NotEmpty(t, token)

	claims, err := sut.ValidateToken(token)

	// Assert
	assert.Nil(t, err)
	require.NotEmpty(t, claims)
	require.Contains(t, claims, "iss")
	assert.Equal(t, claims["iss"], "some-issuer")
}

func Test_HmacOidcProvider_GenerateToken_generates_a_token_with_1h_expiration(t *testing.T) {
	// Arrange
	sut := NewHmacOidcProvider("some-secret", "some-issuer")
	now, _ := time.Parse(time.RFC3339, "2000-01-02T03:04:05.00Z")
	now_1h, _ := time.Parse(time.RFC3339, "2000-01-02T04:04:05.00Z")
	sut.(*HmacOidcProvider).now = func() time.Time {
		return now
	}

	// Act
	token, err := sut.GenerateToken("some-user")
	require.Nil(t, err)
	require.NotEmpty(t, token)

	claims, err := sut.ValidateToken(token)

	// Assert
	assert.Nil(t, err)
	require.NotEmpty(t, claims)
	require.Contains(t, claims, "nbf")
	require.Contains(t, claims, "iat")
	require.Contains(t, claims, "exp")
	assert.Equal(t, float64(now.Unix()), claims["nbf"])
	assert.Equal(t, float64(now.Unix()), claims["iat"])
	assert.Equal(t, float64(now_1h.Unix()), claims["exp"])
}
