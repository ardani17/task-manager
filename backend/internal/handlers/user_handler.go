package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ardani17/taskmanager/internal/middleware"
	"github.com/ardani17/taskmanager/internal/models"
	"github.com/ardani17/taskmanager/internal/repository"
	"github.com/ardani17/taskmanager/pkg/utils"
	"github.com/go-chi/chi/v5"
)

// UserHandler handles user endpoints
type UserHandler struct {
	repo *repository.DeveloperRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo *repository.DeveloperRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// List handles GET /api/v1/users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse pagination params
	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	users, total, err := h.repo.List(limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    users,
		"total":   total,
	})
}

// Get handles GET /api/v1/users/{id}
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.repo.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	if user == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// Update handles PUT /api/v1/users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get current user ID from context
	currentUserID := middleware.GetUserID(r)
	currentUserRole := middleware.GetUserRole(r)

	// Only allow users to update their own profile, or admins to update any
	if id != currentUserID && currentUserRole != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, "You can only update your own profile")
		return
	}

	var req models.UpdateDeveloperRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.repo.Update(id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	if user == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

// Delete handles DELETE /api/v1/users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Only admins can delete users
	currentUserRole := middleware.GetUserRole(r)
	if currentUserRole != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, "Only admins can delete users")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}

// UpdateStatus handles PATCH /api/v1/users/{id}/status
func (h *UserHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Status == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Status is required")
		return
	}

	if err := h.repo.UpdateStatus(id, req.Status); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update status")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Status updated successfully",
	})
}

// Helper to get current time
func now() time.Time {
	return time.Now()
}
