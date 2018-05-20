package api

import (
	"encoding/json"
	"net/http"
)

var serverError = newError("An error occurred, please try again later")

func newError(description string) apiError {
	return apiError{
		Description: description,
	}
}

func httpError(w http.ResponseWriter, error apiError, code int) {
	writeJSON(w, error, code)
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data == nil {
		return
	}
	switch v := data.(type) {
	case []byte:
		w.Write(v)
	case string:
		w.Write([]byte(v))
	default:
		json.NewEncoder(w).Encode(data)
	}
}

type apiError struct {
	Error       bool   `json:"error"`
	Description string `json:"error_description"`
}
