package app_handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

type SumRequest interface{}

type SumResponse struct {
	Sha256Sum string `json:"sha256Sum"`
}

type SumHandler struct {
	// precision string
}

func NewSumHandler() AppHandler[SumRequest, SumResponse] {
	return &SumHandler{}
}

func (h *SumHandler) Handle(request SumRequest) (*SumResponse, error) {
	sum := h.calculateSum(request)
	sum_bytes, err := h.float64ToBytes(sum)
	if err != nil {
		return nil, err
	}

	sum_sha256 := sha256.Sum256([]byte(sum_bytes))

	return &SumResponse{
		Sha256Sum: fmt.Sprintf("%x", sum_sha256),
	}, nil
}

func (h *SumHandler) float64ToBytes(val float64) ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, val); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *SumHandler) calculateSum(request SumRequest) float64 {
	return h.processValue(request)
}

func (h *SumHandler) processValue(jsonPart interface{}) float64 {
	sum := 0.0
	switch actualValue := jsonPart.(type) {
	case map[string]interface{}:
		sum += h.iterateMap(actualValue)

	case []interface{}:
		sum += h.iterateSlice(actualValue)

	case float64:
		sum += actualValue
	}

	return sum
}

func (h *SumHandler) iterateSlice(jsonSlice []interface{}) float64 {
	sum := 0.0
	for _, val := range jsonSlice {
		sum += h.processValue(val)
	}

	return sum
}

func (h *SumHandler) iterateMap(jsonMap map[string]interface{}) float64 {
	sum := 0.0
	for _, val := range jsonMap {
		sum += h.processValue(val)
	}

	return sum
}
