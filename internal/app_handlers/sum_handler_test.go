package app_handlers

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type calculateSumTest struct {
	json     string
	expected float64
	actual   float64
}

func Test_SumHandler_float64ToBytes_converts_float_to_little_endian_bytes(t *testing.T) {
	// Arrange
	val := 0.123
	sut := NewSumHandler().(*SumHandler)

	// Act
	res, err := sut.float64ToBytes(val)

	// Assert
	require.Nil(t, err)
	assert.ElementsMatch(t, []byte{176, 114, 104, 145, 237, 124, 191, 63}, res)
}

func Test_SumHandler_CalculateSum_calculates_sum_correctly(t *testing.T) {
	sut := NewSumHandler().(*SumHandler)

	tests := []calculateSumTest{
		{
			json:     `1`,
			expected: 1,
		},
		{
			json:     `1.2`,
			expected: 1.2,
		},
		{
			json:     `"valid json string"`,
			expected: 0,
		},
		{
			json:     `[1,2,3,4]`,
			expected: 10,
		},
		{
			json:     `{"a":6,"b":4}`,
			expected: 10,
		},
		{
			json:     `[[[2]]]`,
			expected: 2,
		},
		{
			json:     `{"a":{"b":4},"c":-2}`,
			expected: 2,
		},
		{
			json:     `{"a":[-1,1,"dark"]}`,
			expected: 0,
		},
		{
			json:     `[-1,{"a":1, "b":"light"}]`,
			expected: 0,
		},
		{
			json:     `[]`,
			expected: 0,
		},
		{
			json:     `{}`,
			expected: 0,
		},
		{
			json:     `{"a":{"b":0.2, "c":[0.3,-0.5]},"d":1.2,"e":-1.2}`,
			expected: 0,
		},
	}

	for _, test := range tests {
		// Act
		var req SumRequest
		if err := json.Unmarshal([]byte(test.json), &req); err != nil {
			t.Errorf("error during calculation: %s", err)
			t.Fail()
		}

		test.actual = sut.calculateSum(req)

		// Assert
		assert.Equal(t, test.expected, test.actual, fmt.Sprintf("Sum did not match, expected %f, actual %f, json %s", test.expected, test.actual, test.json))
	}
}

func Test_SumHandler_Handle_calculates_and_formats_sha256(t *testing.T) {
	// Arrange
	sut := NewSumHandler()
	var req SumRequest

	err := json.Unmarshal([]byte("[10]"), &req)
	require.Nil(t, err)

	// Act
	res, err := sut.Handle(req)

	// Assert
	assert.Nil(t, err)
	require.NotNil(t, res)

	// Sha256Sum for little endian bytes of 10.0
	assert.Equal(t, "24b1f4ef66b650ff816e519b01742ff1753733d36e1b4c3e3b52743168915b1f", res.Sha256Sum)
}
