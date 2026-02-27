package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ardani17/taskmanager/internal/models"
)

// DeveloperRepository handles database operations for developers
type DeveloperRepository struct {
	db *DB
}

// NewDeveloperRepository creates a new developer repository
func NewDeveloperRepository(db *DB) *DeveloperRepository {
	return &DeveloperRepository{db: db}
}

// Create creates a new developer
func (r *DeveloperRepository) Create(developer *models.Developer) error {
	now := time.Now()

	query := `
		INSERT INTO developers (name, email, role, avatar_url, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		developer.Name,
		developer.Email,
		developer.Role,
		developer.AvatarURL,
		developer.Status,
		now,
		now,
	).Scan(&developer.ID, &developer.CreatedAt, &developer.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create developer: %w", err)
	}

	return nil
}

// GetByID retrieves a developer by ID
func (r *DeveloperRepository) GetByID(id int) (*models.Developer, error) {
	query := `
		SELECT id, name, email, COALESCE(role, 'developer') as role, team_id, avatar_url, status, created_at, updated_at
		FROM developers
		WHERE id = $1
	`

	developer := &models.Developer{}
	var role, avatarURL sql.NullString
	var teamID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&developer.ID,
		&developer.Name,
		&developer.Email,
		&role,
		&teamID,
		&avatarURL,
		&developer.Status,
		&developer.CreatedAt,
		&developer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get developer: %w", err)
	}

	if role.Valid {
		developer.Role = role.String
	}
	if avatarURL.Valid {
		developer.AvatarURL = avatarURL.String
	}
	if teamID.Valid {
		tid := int(teamID.Int64)
		developer.TeamID = &tid
	}

	return developer, nil
}

// GetByEmail retrieves a developer by email
func (r *DeveloperRepository) GetByEmail(email string) (*models.Developer, error) {
	query := `
		SELECT id, name, email, COALESCE(role, 'developer') as role, team_id, avatar_url, status, created_at, updated_at
		FROM developers
		WHERE email = $1
	`

	developer := &models.Developer{}
	var role, avatarURL sql.NullString
	var teamID sql.NullInt64

	err := r.db.QueryRow(query, email).Scan(
		&developer.ID,
		&developer.Name,
		&developer.Email,
		&role,
		&teamID,
		&avatarURL,
		&developer.Status,
		&developer.CreatedAt,
		&developer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get developer by email: %w", err)
	}

	if role.Valid {
		developer.Role = role.String
	}
	if avatarURL.Valid {
		developer.AvatarURL = avatarURL.String
	}
	if teamID.Valid {
		tid := int(teamID.Int64)
		developer.TeamID = &tid
	}

	return developer, nil
}

// List retrieves all developers with pagination
func (r *DeveloperRepository) List(limit, offset int) ([]*models.Developer, int, error) {
	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM developers").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count developers: %w", err)
	}

	// Get developers
	query := `
		SELECT id, name, email, COALESCE(role, 'developer') as role, team_id, avatar_url, status, created_at, updated_at
		FROM developers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list developers: %w", err)
	}
	defer rows.Close()

	var developers []*models.Developer
	for rows.Next() {
		d := &models.Developer{}
		var role, avatarURL sql.NullString
		var teamID sql.NullInt64
		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Email,
			&role,
			&teamID,
			&avatarURL,
			&d.Status,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan developer: %w", err)
		}
		if role.Valid {
			d.Role = role.String
		}
		if avatarURL.Valid {
			d.AvatarURL = avatarURL.String
		}
		if teamID.Valid {
			tid := int(teamID.Int64)
			d.TeamID = &tid
		}
		developers = append(developers, d)
	}

	return developers, total, nil
}

// Update updates a developer
func (r *DeveloperRepository) Update(id int, req *models.UpdateDeveloperRequest) (*models.Developer, error) {
	query := `
		UPDATE developers
		SET name = COALESCE($2, name),
		    avatar_url = COALESCE($3, avatar_url),
		    status = COALESCE($4, status),
		    updated_at = $5
		WHERE id = $1
		RETURNING id, name, email, COALESCE(role, 'developer') as role, team_id, avatar_url, status, created_at, updated_at
	`

	developer := &models.Developer{}
	var role, avatarURL sql.NullString
	var teamID sql.NullInt64

	var avatarURLParam interface{}
	if req.AvatarURL != "" {
		avatarURLParam = req.AvatarURL
	}

	err := r.db.QueryRow(
		query,
		id,
		req.Name,
		avatarURLParam,
		req.Status,
		time.Now(),
	).Scan(
		&developer.ID,
		&developer.Name,
		&developer.Email,
		&role,
		&teamID,
		&avatarURL,
		&developer.Status,
		&developer.CreatedAt,
		&developer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update developer: %w", err)
	}

	if role.Valid {
		developer.Role = role.String
	}
	if avatarURL.Valid {
		developer.AvatarURL = avatarURL.String
	}
	if teamID.Valid {
		tid := int(teamID.Int64)
		developer.TeamID = &tid
	}

	return developer, nil
}

// UpdateStatus updates developer status
func (r *DeveloperRepository) UpdateStatus(id int, status string) error {
	query := `
		UPDATE developers
		SET status = $2, updated_at = $3
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update developer status: %w", err)
	}

	return nil
}

// Delete deletes a developer
func (r *DeveloperRepository) Delete(id int) error {
	query := "DELETE FROM developers WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete developer: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("developer not found")
	}

	return nil
}
