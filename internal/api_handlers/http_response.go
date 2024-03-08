package api_handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpErrorResponse struct {
	Error string `json:"error"`
}

func HttpError(w http.ResponseWriter, err string, status_code int) {
	http_err := &HttpErrorResponse{
		Error: err,
	}

	response, marshal_err := json.Marshal(http_err)
	if marshal_err != nil {
		log.Printf("unable to marshal error repsonse, %s\n", err)
		response = []byte{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	w.Write(response)
}

func HttpSuccess(w http.ResponseWriter, result interface{}) {
	response, marshal_err := json.Marshal(result)
	if marshal_err != nil {
		log.Println("unable to marshal success repsonse")
		HttpError(w, "error while handling request", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
