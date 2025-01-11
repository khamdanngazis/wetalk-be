package middleware

import (
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewHTTPResponse creates a new HTTPResponse instance with the provided code, message, and optional data.
func NewHTTPResponse(message string, data interface{}) *HTTPResponse {
	return &HTTPResponse{
		Message: message,
		Data:    data,
	}
}

func WriteResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := NewHTTPResponse(message, data)
	WriteJSONResponse(w, response, statusCode)
}

func WriteJSONResponse(w http.ResponseWriter, response *HTTPResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
