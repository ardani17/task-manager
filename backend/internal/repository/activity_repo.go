package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ardani17/taskmanager/internal/models"
)

// ActivityRepository handles database operations for activity logs
type ActivityRepository struct {
	db *DB
}

// NewActivityRepository creates a new activity repository
func NewActivityRepository(db *DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

// Create creates a new activity log entry
func (r *ActivityRepository) Create(activity *models.Activity) error {
	now := time.Now()

	query := `
		INSERT INTO activities (developer_id, task_id, action, description, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		activity.DeveloperID,
		activity.TaskID,
		activity.Action,
		activity.Description,
		activity.Metadata,
		now,
	).Scan(&activity.ID, &activity.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	return nil
}

// List retrieves activity logs with pagination and filters
func (r *ActivityRepository) List(limit, offset int, developerID, taskID int) ([]*models.Activity, int, error) {
	// Build query with filters
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if developerID > 0 {
		whereClause += fmt.Sprintf(" AND a.developer_id = $%d", argIndex)
		args = append(args, developerID)
		argIndex++
	}
	if taskID > 0 {
		whereClause += fmt.Sprintf(" AND a.task_id = $%d", argIndex)
		args = append(args, taskID)
		argIndex++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM activities a %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	// Get activities
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT a.id, a.developer_id, a.task_id, a.action, a.description, a.metadata, a.created_at,
		       d.id, d.name, d.email, d.status
		FROM activities a
		LEFT JOIN developers d ON a.developer_id = d.id
		%s
		ORDER BY a.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list activities: %w", err)
	}
	defer rows.Close()

	var activities []*models.Activity
	for rows.Next() {
		a := &models.Activity{}
		var developerID, taskID sql.NullInt64
		var developer models.Developer

		err := rows.Scan(
			&a.ID,
			&developerID,
			&taskID,
			&a.Action,
			&a.Description,
			&a.Metadata,
			&a.CreatedAt,
			&developer.ID,
			&developer.Name,
			&developer.Email,
			&developer.Status,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan activity: %w", err)
		}

		if developerID.Valid {
			did := int(developerID.Int64)
			a.DeveloperID = &did
			a.Developer = &developer
		}
		if taskID.Valid {
			tid := int(taskID.Int64)
			a.TaskID = &tid
		}

		activities = append(activities, a)
	}

	return activities, total, nil
}
