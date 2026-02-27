package handlers

import (
	"net/http"
	"strconv"

	"github.com/ardani17/taskmanager/internal/repository"
	"github.com/ardani17/taskmanager/pkg/utils"
)

// ActivityHandler handles activity endpoints
type ActivityHandler struct {
	repo *repository.ActivityRepository
}

// NewActivityHandler creates a new activity handler
func NewActivityHandler(repo *repository.ActivityRepository) *ActivityHandler {
	return &ActivityHandler{repo: repo}
}

// List handles GET /api/v1/activity
func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) {
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

	// Parse filters
	developerID := 0
	taskID := 0
	if d := r.URL.Query().Get("developer_id"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val > 0 {
			developerID = val
		}
	}
	if t := r.URL.Query().Get("task_id"); t != "" {
		if val, err := strconv.Atoi(t); err == nil && val > 0 {
			taskID = val
		}
	}

	activities, total, err := h.repo.List(limit, offset, developerID, taskID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch activities")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    activities,
		"total":   total,
	})
}
