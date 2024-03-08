package main

import (
	"coding_exercise/internal/lib"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Integration_Main_initializeRouter_configures_auth_endpoint(t *testing.T) {
	// Arrange
	config := &config{
		secret: "some-secret",
		issuer: "some-issuer",
	}

	sut := initializeRouter(config)
	recorder := httptest.NewRecorder()
	body := strings.NewReader(`{"username":"some-user","password":"some-password"}`)
	req := httptest.NewRequest("POST", "/auth", body)
	req.Header.Add("Content-Type", "application/json")

	// Act
	sut.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, 200, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `{"token":"`)
}

func Test_Integration_Main_initializeRouter_configures_sum_endpoint_with_authentication_401(t *testing.T) {
	// Arrange
	config := &config{
		secret: "some-secret",
		issuer: "some-issuer",
	}

	sut := initializeRouter(config)
	recorder := httptest.NewRecorder()
	body := strings.NewReader(`[1,2]`)
	req := httptest.NewRequest("POST", "/sum", body)
	req.Header.Add("Content-Type", "application/json")

	// Act
	sut.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, 401, recorder.Code)
}

func Test_Integration_Main_initializeRouter_configures_sum_endpoint_with_authentication_200(t *testing.T) {
	// Arrange
	config := &config{
		secret: "some-secret",
		issuer: "some-issuer",
	}

	sut := initializeRouter(config)
	recorder := httptest.NewRecorder()
	body := strings.NewReader(`[1,2]`)
	req := httptest.NewRequest("POST", "/sum", body)
	req.Header.Add("Content-Type", "application/json")

	oidc_provider := lib.NewHmacOidcProvider(config.secret, config.issuer)
	token, err := oidc_provider.GenerateToken("some-username")
	require.Nil(t, err)

	req.Header.Add("Authorization", "Bearer "+token)

	// Act
	sut.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, 200, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `{"sha256Sum":"`)
}
