package api_handlers

import (
	"coding_exercise/internal/app_handlers"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SumHandler_returns_400_on_empty_body(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.SumHandlerMock{}
	sut := NewSumHandler(app_handler_mock)

	req := httptest.NewRequest("POST", "/", nil)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.False(t, app_handler_mock.HandleCalled)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func Test_SumHandler_returns_400_on_invalid_json_in_body(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.SumHandlerMock{}
	sut := NewSumHandler(app_handler_mock)

	body := strings.NewReader("invalid-json")
	req := httptest.NewRequest("POST", "/", body)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.False(t, app_handler_mock.HandleCalled)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func Test_SumHandler_calls_handle_with_the_specified_json(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.SumHandlerMock{}
	sut := NewSumHandler(app_handler_mock)

	body := strings.NewReader("[1]")
	req := httptest.NewRequest("POST", "/", body)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.True(t, app_handler_mock.HandleCalled)
	last_req, err := json.Marshal(app_handler_mock.LastRequest)
	require.Nil(t, err)
	assert.Equal(t, "[1]", string(last_req))
}

func Test_SumHandler_returns_the_received_result(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.SumHandlerMock{
		NextResponse: &app_handlers.SumResponse{Sha256Sum: "handler-response"},
	}
	sut := NewSumHandler(app_handler_mock)

	body := strings.NewReader("[1]")
	req := httptest.NewRequest("POST", "/", body)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.True(t, app_handler_mock.HandleCalled)
	assert.Equal(t, http.StatusOK, recorder.Code)

	require.NotNil(t, recorder.Body)
	res, err := io.ReadAll(recorder.Body)
	require.Nil(t, err)
	assert.Equal(t, `{"sha256Sum":"handler-response"}`, string(res))
}

func Test_SumHandler_returns_internalservererror_on_sumhandler_error(t *testing.T) {
	// Arrange
	app_handler_mock := &app_handlers.SumHandlerMock{
		NextError: errors.New("some-error"),
	}
	sut := NewSumHandler(app_handler_mock)

	body := strings.NewReader("[1]")
	req := httptest.NewRequest("POST", "/", body)
	recorder := httptest.NewRecorder()

	// Act
	sut.Handle(recorder, req)

	// Assert
	assert.True(t, app_handler_mock.HandleCalled)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	require.NotNil(t, recorder.Body)
	res, err := io.ReadAll(recorder.Body)
	require.Nil(t, err)
	assert.Contains(t, string(res), `{"error":`)
}
