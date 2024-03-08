package api_handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HttpError_return_correct_message_and_status_code(t *testing.T) {
	// Arrange
	recorder := httptest.NewRecorder()

	// Act
	HttpError(recorder, "some error", 500)

	// Assert
	assert.Equal(t, 500, recorder.Code)
	assert.Equal(t, `{"error":"some error"}`, recorder.Body.String())
}

func Test_HttpSuccess_return_correct_message_and_status_code(t *testing.T) {
	// Arrange
	recorder := httptest.NewRecorder()

	// Act
	HttpSuccess(recorder, map[string]interface{}{"test": "value"})

	// Assert
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, `{"test":"value"}`, recorder.Body.String())
}
