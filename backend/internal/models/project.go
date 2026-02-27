package models

import (
	"time"
)

// Project represents a project in the system
type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	StartDate   *string   `json:"start_date,omitempty"`
	EndDate     *string   `json:"end_date,omitempty"`
	TeamID      *int      `json:"team_id,omitempty"`
	TaskCount   int       `json:"task_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProjectRequest represents a project creation request
type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	TeamID      *int   `json:"team_id,omitempty"`
}

// UpdateProjectRequest represents a project update request
type UpdateProjectRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	TeamID      *int   `json:"team_id,omitempty"`
}

// ProjectListResponse represents a list of projects
type ProjectListResponse struct {
	Success bool       `json:"success"`
	Data    []*Project `json:"data"`
	Total   int        `json:"total"`
}

// Validate validates the create project request
func (r *CreateProjectRequest) Validate() []string {
	var errors []string

	if r.Name == "" {
		errors = append(errors, "Name is required")
	}
	if len(r.Name) < 3 {
		errors = append(errors, "Name must be at least 3 characters")
	}

	// Set default status
	if r.Status == "" {
		r.Status = "active"
	}

	// Validate status
	validStatuses := map[string]bool{"active": true, "archived": true, "completed": true}
	if !validStatuses[r.Status] {
		errors = append(errors, "Invalid status. Must be one of: active, archived, completed")
	}

	return errors
}
