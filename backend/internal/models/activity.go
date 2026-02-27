package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Activity represents an activity log entry
type Activity struct {
	ID          int        `json:"id"`
	DeveloperID *int       `json:"developer_id,omitempty"`
	TaskID      *int       `json:"task_id,omitempty"`
	Developer   *Developer `json:"developer,omitempty"`
	Action      string     `json:"action"`
	Description string     `json:"description,omitempty"`
	Metadata    JSONB      `json:"metadata,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// JSONB is a custom type for handling JSONB data
type JSONB map[string]interface{}

// Value implements driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// ActivityListResponse represents a list of activities
type ActivityListResponse struct {
	Success bool         `json:"success"`
	Data    []*Activity  `json:"data"`
	Total   int          `json:"total"`
}

// CreateActivityRequest represents an activity creation request
type CreateActivityRequest struct {
	DeveloperID int    `json:"developer_id"`
	TaskID      *int   `json:"task_id,omitempty"`
	Action      string `json:"action"`
	Description string `json:"description,omitempty"`
	Metadata    JSONB  `json:"metadata,omitempty"`
}

// Common activity actions
const (
	ActionTaskCreated   = "task_created"
	ActionTaskUpdated   = "task_updated"
	ActionTaskDeleted   = "task_deleted"
	ActionTaskCompleted = "task_completed"

	ActionProjectCreated = "project_created"
	ActionProjectUpdated = "project_updated"
	ActionProjectDeleted = "project_deleted"

	ActionUserLoggedIn  = "user_logged_in"
	ActionUserLoggedOut = "user_logged_out"
)
