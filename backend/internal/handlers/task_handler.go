package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ardani17/taskmanager/internal/middleware"
	"github.com/ardani17/taskmanager/internal/models"
	"github.com/ardani17/taskmanager/internal/repository"
	"github.com/ardani17/taskmanager/pkg/utils"
	"github.com/go-chi/chi/v5"
)

// TaskHandler handles task endpoints
type TaskHandler struct {
	repo         *repository.TaskRepository
	activityRepo *repository.ActivityRepository
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(repo *repository.TaskRepository, activityRepo *repository.ActivityRepository) *TaskHandler {
	return &TaskHandler{
		repo:         repo,
		activityRepo: activityRepo,
	}
}

// List handles GET /api/v1/tasks
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
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
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	tasks, total, err := h.repo.List(limit, offset, status, priority)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tasks,
		"total":   total,
	})
}

// Get handles GET /api/v1/tasks/{id}
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.repo.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch task")
		return
	}

	if task == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Task not found")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    task,
	})
}

// Create handles POST /api/v1/tasks
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := req.Validate(); len(errors) > 0 {
		utils.ValidationErrorResponse(w, errors)
		return
	}

	// Create task model
	task := &models.Task{
		Title:          req.Title,
		Description:    req.Description,
		AssigneeID:     req.AssigneeID,
		ProjectID:      req.ProjectID,
		Status:         req.Status,
		Priority:       req.Priority,
		DueDate:        req.DueDate,
		EstimatedHours: req.EstimatedHours,
	}

	// Save to database
	if err := h.repo.Create(task); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	// Log activity
	userID := middleware.GetUserID(r)
	activity := &models.Activity{
		DeveloperID: &userID,
		TaskID:      &task.ID,
		Action:      models.ActionTaskCreated,
		Description: "Task created: " + task.Title,
		Metadata: models.JSONB{
			"title":    task.Title,
			"status":   task.Status,
			"priority": task.Priority,
		},
		CreatedAt: now(),
	}
	h.activityRepo.Create(activity)

	utils.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Task created successfully",
		"data":    task,
	})
}

// Update handles PUT /api/v1/tasks/{id}
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.repo.Update(id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	if task == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Task not found")
		return
	}

	// Log activity
	userID := middleware.GetUserID(r)
	activity := &models.Activity{
		DeveloperID: &userID,
		TaskID:      &task.ID,
		Action:      models.ActionTaskUpdated,
		Description: "Task updated: " + task.Title,
		Metadata: models.JSONB{
			"title": task.Title,
		},
		CreatedAt: now(),
	}
	h.activityRepo.Create(activity)

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Task updated successfully",
		"data":    task,
	})
}

// Delete handles DELETE /api/v1/tasks/{id}
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get task before deleting (for activity log)
	task, _ := h.repo.GetByID(id)

	if err := h.repo.Delete(id); err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	// Log activity
	if task != nil {
		userID := middleware.GetUserID(r)
		activity := &models.Activity{
			DeveloperID: &userID,
			TaskID:      &task.ID,
			Action:      models.ActionTaskDeleted,
			Description: "Task deleted: " + task.Title,
			Metadata: models.JSONB{
				"title": task.Title,
			},
			CreatedAt: now(),
		}
		h.activityRepo.Create(activity)
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Task deleted successfully",
	})
}

// UpdateStatus handles PATCH /api/v1/tasks/{id}/status
func (h *TaskHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
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

	// Validate status
	validStatuses := map[string]bool{"todo": true, "in_progress": true, "review": true, "done": true}
	if !validStatuses[req.Status] {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid status")
		return
	}

	if err := h.repo.UpdateStatus(id, req.Status); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update status")
		return
	}

	// Log activity
	userID := middleware.GetUserID(r)
	action := models.ActionTaskUpdated
	if req.Status == "done" {
		action = models.ActionTaskCompleted
	}
	activity := &models.Activity{
		DeveloperID: &userID,
		TaskID:      &id,
		Action:      action,
		Description: "Task status changed to: " + req.Status,
		Metadata: models.JSONB{
			"new_status": req.Status,
		},
		CreatedAt: now(),
	}
	h.activityRepo.Create(activity)

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Status updated successfully",
	})
}
