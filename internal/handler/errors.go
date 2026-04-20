package handler

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, message string, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Error:     message,
		RequestID: requestID,
	})
}
