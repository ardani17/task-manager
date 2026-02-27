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

// ProjectHandler handles project endpoints
type ProjectHandler struct {
	repo         *repository.ProjectRepository
	activityRepo *repository.ActivityRepository
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(repo *repository.ProjectRepository, activityRepo *repository.ActivityRepository) *ProjectHandler {
	return &ProjectHandler{
		repo:         repo,
		activityRepo: activityRepo,
	}
}

// List handles GET /api/v1/projects
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
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

	projects, total, err := h.repo.List(limit, offset, status)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    projects,
		"total":   total,
	})
}

// Get handles GET /api/v1/projects/{id}
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.repo.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch project")
		return
	}

	if project == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Project not found")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    project,
	})
}

// Create handles POST /api/v1/projects
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := req.Validate(); len(errors) > 0 {
		utils.ValidationErrorResponse(w, errors)
		return
	}

	// Create project model
	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		TeamID:      req.TeamID,
	}

	// Save to database
	if err := h.repo.Create(project); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	// Log activity
	userID := middleware.GetUserID(r)
	activity := &models.Activity{
		DeveloperID: &userID,
		Action:      models.ActionProjectCreated,
		Description: "Project created: " + project.Name,
		Metadata: models.JSONB{
			"name":   project.Name,
			"status": project.Status,
		},
		CreatedAt: now(),
	}
	h.activityRepo.Create(activity)

	utils.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Project created successfully",
		"data":    project,
	})
}

// Update handles PUT /api/v1/projects/{id}
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req models.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project, err := h.repo.Update(id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	if project == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Project not found")
		return
	}

	// Log activity
	userID := middleware.GetUserID(r)
	activity := &models.Activity{
		DeveloperID: &userID,
		Action:      models.ActionProjectUpdated,
		Description: "Project updated: " + project.Name,
		Metadata: models.JSONB{
			"name": project.Name,
		},
		CreatedAt: now(),
	}
	h.activityRepo.Create(activity)

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Project updated successfully",
		"data":    project,
	})
}

// Delete handles DELETE /api/v1/projects/{id}
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Get project before deleting (for activity log)
	project, _ := h.repo.GetByID(id)

	if err := h.repo.Delete(id); err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	// Log activity
	if project != nil {
		userID := middleware.GetUserID(r)
		activity := &models.Activity{
			DeveloperID: &userID,
			Action:      models.ActionProjectDeleted,
			Description: "Project deleted: " + project.Name,
			Metadata: models.JSONB{
				"name": project.Name,
			},
			CreatedAt: now(),
		}
		h.activityRepo.Create(activity)
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Project deleted successfully",
	})
}
