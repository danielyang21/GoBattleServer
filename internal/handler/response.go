package handler

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	ErrCodeBadRequest          = "bad_request"
	ErrCodeUnauthorized        = "unauthorized"
	ErrCodeNotFound            = "not_found"
	ErrCodeConflict            = "conflict"
	ErrCodeInternalServerError = "internal_server_error"
	ErrCodeCooldownActive      = "cooldown_active"
	ErrCodeInsufficientCoins   = "insufficient_coins"
)

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func RespondError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// RespondBadRequest sends a 400 Bad Request response
func RespondBadRequest(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusBadRequest, ErrCodeBadRequest, message)
}

// RespondNotFound sends a 404 Not Found response
func RespondNotFound(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusNotFound, ErrCodeNotFound, message)
}

// RespondConflict sends a 409 Conflict response
func RespondConflict(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusConflict, ErrCodeConflict, message)
}

// RespondInternalError sends a 500 Internal Server Error response
func RespondInternalError(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusInternalServerError, ErrCodeInternalServerError, message)
}
