package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ardani17/taskmanager/internal/models"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project
func (r *ProjectRepository) Create(project *models.Project) error {
	now := time.Now()

	query := `
		INSERT INTO projects (name, description, status, team_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		project.Name,
		project.Description,
		project.Status,
		project.TeamID,
		now,
		now,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
	query := `
		SELECT id, name, description, status, start_date, end_date, team_id, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	project := &models.Project{}
	var startDate, endDate sql.NullString
	var teamID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&startDate,
		&endDate,
		&teamID,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	if startDate.Valid {
		project.StartDate = &startDate.String
	}
	if endDate.Valid {
		project.EndDate = &endDate.String
	}
	if teamID.Valid {
		tid := int(teamID.Int64)
		project.TeamID = &tid
	}

	// Get task count
	err = r.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE project_id = $1", id).Scan(&project.TaskCount)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get task count: %w", err)
	}

	return project, nil
}

// List retrieves all projects with pagination
func (r *ProjectRepository) List(limit, offset int, status string) ([]*models.Project, int, error) {
	// Build query with filters
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM projects %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	// Get projects
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT id, name, description, status, start_date, end_date, team_id, created_at, updated_at
		FROM projects
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		p := &models.Project{}
		var startDate, endDate sql.NullString
		var teamID sql.NullInt64
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Status,
			&startDate,
			&endDate,
			&teamID,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan project: %w", err)
		}
		if startDate.Valid {
			p.StartDate = &startDate.String
		}
		if endDate.Valid {
			p.EndDate = &endDate.String
		}
		if teamID.Valid {
			tid := int(teamID.Int64)
			p.TeamID = &tid
		}
		projects = append(projects, p)
	}

	return projects, total, nil
}

// Update updates a project
func (r *ProjectRepository) Update(id int, req *models.UpdateProjectRequest) (*models.Project, error) {
	query := `
		UPDATE projects
		SET name = COALESCE($2, name),
		    description = COALESCE($3, description),
		    status = COALESCE($4, status),
		    team_id = COALESCE($5, team_id),
		    updated_at = $6
		WHERE id = $1
		RETURNING id, name, description, status, start_date, end_date, team_id, created_at, updated_at
	`

	project := &models.Project{}
	var startDate, endDate sql.NullString
	var teamID sql.NullInt64

	err := r.db.QueryRow(
		query,
		id,
		req.Name,
		req.Description,
		req.Status,
		req.TeamID,
		time.Now(),
	).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&startDate,
		&endDate,
		&teamID,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	if startDate.Valid {
		project.StartDate = &startDate.String
	}
	if endDate.Valid {
		project.EndDate = &endDate.String
	}
	if teamID.Valid {
		tid := int(teamID.Int64)
		project.TeamID = &tid
	}

	return project, nil
}

// Delete deletes a project
func (r *ProjectRepository) Delete(id int) error {
	query := "DELETE FROM projects WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}
