package handlers

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Message string      `json:"message,omitempty"`
	Data    any `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}