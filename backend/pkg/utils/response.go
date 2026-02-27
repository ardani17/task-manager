package utils

import (
	"encoding/json"
	"net/http"
)

// JSON response helper
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error response helper
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{
		"error":   http.StatusText(status),
		"message": message,
	})
}

// Success response helper
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// Created response helper
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// ErrorResponse sends a structured error response
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    status,
			"message": message,
		},
	})
}

// ValidationErrorResponse sends validation error details
func ValidationErrorResponse(w http.ResponseWriter, errors []string) {
	JSON(w, http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Validation failed",
			"details": errors,
		},
	})
}
