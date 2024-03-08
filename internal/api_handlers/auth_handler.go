package api_handlers

import (
	"coding_exercise/internal/app_handlers"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthHandler struct {
	app_handler app_handlers.AppHandler[app_handlers.AuthRequest, app_handlers.AuthResponse]
}

func NewAuthHandler(app_handler app_handlers.AppHandler[app_handlers.AuthRequest, app_handlers.AuthResponse]) *AuthHandler {
	return &AuthHandler{
		app_handler: app_handler,
	}
}

func (h *AuthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req app_handlers.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HttpError(w, "unable to read body", http.StatusBadRequest)
		return
	}

	res, err := h.app_handler.Handle(req)

	if err != nil {
		if errors.Is(err, app_handlers.ErrAuthValidationError) {
			HttpError(w, err.Error(), http.StatusBadRequest)
		} else {
			HttpError(w, "error while generating token", http.StatusInternalServerError)
		}

		return
	}

	HttpSuccess(w, res)
}
