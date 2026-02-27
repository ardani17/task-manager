package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ardani17/taskmanager/internal/models"
)

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	db *DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create creates a new task
func (r *TaskRepository) Create(task *models.Task) error {
	now := time.Now()

	query := `
		INSERT INTO tasks (title, description, status, priority, project_id, assignee_id, due_date, estimated_hours, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.ProjectID,
		task.AssigneeID,
		task.DueDate,
		task.EstimatedHours,
		now,
		now,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(id int) (*models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, project_id, assignee_id, 
		       due_date, COALESCE(estimated_hours, 0)::float, COALESCE(actual_hours, 0)::float, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	task := &models.Task{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.ProjectID,
		&task.AssigneeID,
		&task.DueDate,
		&task.EstimatedHours,
		&task.ActualHours,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// List retrieves all tasks with pagination and filters
func (r *TaskRepository) List(limit, offset int, status, priority string) ([]*models.Task, int, error) {
	// Build query with filters
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if priority != "" {
		whereClause += fmt.Sprintf(" AND priority = $%d", argIndex)
		args = append(args, priority)
		argIndex++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	// Get tasks
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT id, title, description, status, priority, project_id, assignee_id, 
		       due_date, COALESCE(estimated_hours, 0)::float, COALESCE(actual_hours, 0)::float, created_at, updated_at
		FROM tasks
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		t := &models.Task{}
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.ProjectID,
			&t.AssigneeID,
			&t.DueDate,
			&t.EstimatedHours,
			&t.ActualHours,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, total, nil
}

// Update updates a task
func (r *TaskRepository) Update(id int, req *models.UpdateTaskRequest) (*models.Task, error) {
	query := `
		UPDATE tasks
		SET title = COALESCE($2, title),
		    description = COALESCE($3, description),
		    status = COALESCE($4, status),
		    priority = COALESCE($5, priority),
		    project_id = COALESCE($6, project_id),
		    assignee_id = COALESCE($7, assignee_id),
		    due_date = COALESCE($8, due_date),
		    estimated_hours = COALESCE($9, estimated_hours),
		    actual_hours = COALESCE($10, actual_hours),
		    updated_at = $11
		WHERE id = $1
		RETURNING id, title, description, status, priority, project_id, assignee_id, 
		          due_date, COALESCE(estimated_hours, 0)::float, COALESCE(actual_hours, 0)::float, created_at, updated_at
	`

	task := &models.Task{}
	err := r.db.QueryRow(
		query,
		id,
		req.Title,
		req.Description,
		req.Status,
		req.Priority,
		req.ProjectID,
		req.AssigneeID,
		req.DueDate,
		req.EstimatedHours,
		req.ActualHours,
		time.Now(),
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.ProjectID,
		&task.AssigneeID,
		&task.DueDate,
		&task.EstimatedHours,
		&task.ActualHours,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// UpdateStatus updates task status
func (r *TaskRepository) UpdateStatus(id int, status string) error {
	query := "UPDATE tasks SET status = $2, updated_at = $3 WHERE id = $1"
	_, err := r.db.Exec(query, id, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}

// Delete deletes a task
func (r *TaskRepository) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
