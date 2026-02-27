package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Uptime    string `json:"uptime,omitempty"`
}

type APIInfoResponse struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Endpoints []struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Desc   string `json:"desc"`
	} `json:"endpoints"`
}

var startTime = time.Now()

// Health check endpoint
func Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime).String(),
	}

	JSON(w, http.StatusOK, response)
}

// API info endpoint
func APIInfo(w http.ResponseWriter, r *http.Request) {
	response := APIInfoResponse{
		Name:    "TaskManager API",
		Version: "1.0.0",
		Endpoints: []struct {
			Method string `json:"method"`
			Path   string `json:"path"`
			Desc   string `json:"desc"`
		}{
			{"GET", "/health", "Health check"},
			{"GET", "/api/v1", "API info"},
			{"POST", "/api/v1/auth/register", "Register new user"},
			{"POST", "/api/v1/auth/login", "Login user"},
			{"GET", "/api/v1/auth/me", "Get current user"},
			{"GET", "/api/v1/users", "List users"},
			{"GET", "/api/v1/projects", "List projects"},
			{"POST", "/api/v1/projects", "Create project"},
			{"GET", "/api/v1/tasks", "List tasks"},
			{"GET", "/api/v1/activity", "Activity log"},
		},
	}

	JSON(w, http.StatusOK, response)
}

// System info endpoint (for debugging)
func SystemInfo(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"go_version": runtime.Version(),
		"goroutines": runtime.NumGoroutine(),
		"memory": map[string]interface{}{
			"alloc_mb":       memStats.Alloc / 1024 / 1024,
			"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
			"sys_mb":         memStats.Sys / 1024 / 1024,
		},
	}

	JSON(w, http.StatusOK, response)
}

// JSON helper
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error helper
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{
		"error":   http.StatusText(status),
		"message": message,
	})
}
