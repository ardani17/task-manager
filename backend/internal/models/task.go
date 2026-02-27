package models

import (
	"time"
)

// Task represents a task in the system
type Task struct {
	ID             int        `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description,omitempty"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	ProjectID      *int       `json:"project_id,omitempty"`
	AssigneeID     *int       `json:"assignee_id,omitempty"`
	Assignee       *Developer `json:"assignee,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedHours float64    `json:"estimated_hours,omitempty"`
	ActualHours    float64    `json:"actual_hours,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CreateTaskRequest represents a task creation request
type CreateTaskRequest struct {
	Title          string     `json:"title"`
	Description    string     `json:"description,omitempty"`
	Status         string     `json:"status,omitempty"`
	Priority       string     `json:"priority,omitempty"`
	ProjectID      *int       `json:"project_id,omitempty"`
	AssigneeID     *int       `json:"assignee_id,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedHours float64    `json:"estimated_hours,omitempty"`
}

// UpdateTaskRequest represents a task update request
type UpdateTaskRequest struct {
	Title          string     `json:"title,omitempty"`
	Description    string     `json:"description,omitempty"`
	Status         string     `json:"status,omitempty"`
	Priority       string     `json:"priority,omitempty"`
	ProjectID      *int       `json:"project_id,omitempty"`
	AssigneeID     *int       `json:"assignee_id,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedHours float64    `json:"estimated_hours,omitempty"`
	ActualHours    float64    `json:"actual_hours,omitempty"`
}

// TaskListResponse represents a list of tasks with pagination
type TaskListResponse struct {
	Success bool    `json:"success"`
	Data    []*Task `json:"data"`
	Total   int     `json:"total"`
}

// Validate validates the create task request
func (r *CreateTaskRequest) Validate() []string {
	var errors []string

	if r.Title == "" {
		errors = append(errors, "Title is required")
	}
	if len(r.Title) < 3 {
		errors = append(errors, "Title must be at least 3 characters")
	}

	// Set defaults
	if r.Status == "" {
		r.Status = "todo"
	}
	if r.Priority == "" {
		r.Priority = "medium"
	}

	// Validate status
	validStatuses := map[string]bool{"todo": true, "in_progress": true, "review": true, "done": true}
	if !validStatuses[r.Status] {
		errors = append(errors, "Invalid status. Must be one of: todo, in_progress, review, done")
	}

	// Validate priority
	validPriorities := map[string]bool{"low": true, "medium": true, "high": true}
	if !validPriorities[r.Priority] {
		errors = append(errors, "Invalid priority. Must be one of: low, medium, high")
	}

	return errors
}
