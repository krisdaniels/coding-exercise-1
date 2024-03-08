package api_handlers

import (
	"coding_exercise/internal/app_handlers"
	"encoding/json"
	"log"
	"net/http"
)

type SumHandler struct {
	app_handler app_handlers.AppHandler[app_handlers.SumRequest, app_handlers.SumResponse]
}

func NewSumHandler(app_handler app_handlers.AppHandler[app_handlers.SumRequest, app_handlers.SumResponse]) *SumHandler {
	return &SumHandler{
		app_handler: app_handler,
	}
}

func (h *SumHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req app_handlers.SumRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HttpError(w, "unable to read body", http.StatusBadRequest)
		return
	}

	res, err := h.app_handler.Handle(req)

	if err != nil {
		log.Printf("unable to handle sum request: %s\n", err)
		HttpError(w, "error while handling sum request", http.StatusInternalServerError)
		return
	}

	HttpSuccess(w, res)
}
