package api_handlers

import (
	"coding_exercise/internal/app_handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AuthHandler_returns_400_on_empty_body(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.AuthHandlerMock{}
	sut := NewAuthHandler(app_handler_mock)

	req := httptest.NewRequest("POST", "/", nil)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.False(t, app_handler_mock.HandleCalled)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func Test_AuthHandler_returns_400_on_invalid_json_in_body(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.AuthHandlerMock{}
	sut := NewAuthHandler(app_handler_mock)

	body := strings.NewReader("invalid-json")
	req := httptest.NewRequest("POST", "/", body)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.False(t, app_handler_mock.HandleCalled)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func Test_AuthHandler_calls_app_handler_with_specified_username_and_password(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.AuthHandlerMock{}
	sut := NewAuthHandler(app_handler_mock)

	body := strings.NewReader(`{"username":"some-username","password":"some-password"}`)
	req := httptest.NewRequest("POST", "/", body)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.True(t, app_handler_mock.HandleCalled)
	require.NotNil(t, app_handler_mock.LastRequest)
	assert.Equal(t, "some-username", app_handler_mock.LastRequest.Username)
	assert.Equal(t, "some-password", app_handler_mock.LastRequest.Password)
}

func Test_AuthHandler_returns_400_on_validation_error(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.AuthHandlerMock{
		NextError: app_handlers.ErrAuthValidationError,
	}
	sut := NewAuthHandler(app_handler_mock)

	body := strings.NewReader(`{}`)
	req := httptest.NewRequest("POST", "/", body)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.True(t, app_handler_mock.HandleCalled)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), app_handlers.ErrAuthValidationError.Error())
}

func Test_AuthHandler_returns_500_on_token_generation_failure(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.AuthHandlerMock{
		NextError: app_handlers.ErrAuthTokenGenerationError,
	}
	sut := NewAuthHandler(app_handler_mock)

	body := strings.NewReader(`{}`)
	req := httptest.NewRequest("POST", "/", body)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.True(t, app_handler_mock.HandleCalled)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func Test_AuthHandler_returns_generated_token(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.AuthHandlerMock{
		NextResponse: &app_handlers.AuthResponse{Token: "some-token"},
	}
	sut := NewAuthHandler(app_handler_mock)

	body := strings.NewReader(`{"username":"some-username","password":"some-password"}`)
	req := httptest.NewRequest("POST", "/", body)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.True(t, app_handler_mock.HandleCalled)
	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"token":"some-token"}`, recorder.Body.String())
}
