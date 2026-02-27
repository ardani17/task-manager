package models

import (
	"time"
)

// Developer represents a developer in the system
type Developer struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`              // Never expose password hash
	Role         string    `json:"role"`
	TeamID       *int      `json:"team_id,omitempty"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"` // Optional, defaults to "developer"
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *LoginResponseData `json:"data"`
}

// LoginResponseData contains the login response data
type LoginResponseData struct {
	Developer *Developer `json:"developer"`
	Token     *TokenData `json:"token"`
}

// TokenData contains token information
type TokenData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// MeResponse represents the current user response
type MeResponse struct {
	Success bool       `json:"success"`
	Data    *Developer `json:"data"`
}

// UpdateDeveloperRequest represents an update request
type UpdateDeveloperRequest struct {
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Status    string `json:"status,omitempty"`
}

// Validate validates the registration request
func (r *RegisterRequest) Validate() []string {
	var errors []string

	if r.Name == "" {
		errors = append(errors, "Name is required")
	}
	if len(r.Name) < 2 {
		errors = append(errors, "Name must be at least 2 characters")
	}
	if r.Email == "" {
		errors = append(errors, "Email is required")
	}
	if r.Password == "" {
		errors = append(errors, "Password is required")
	}
	if len(r.Password) < 6 {
		errors = append(errors, "Password must be at least 6 characters")
	}

	return errors
}

// Validate validates the login request
func (r *LoginRequest) Validate() []string {
	var errors []string

	if r.Email == "" {
		errors = append(errors, "Email is required")
	}
	if r.Password == "" {
		errors = append(errors, "Password is required")
	}

	return errors
}

// ToDeveloper creates a Developer from RegisterRequest
func (r *RegisterRequest) ToDeveloper(passwordHash string) *Developer {
	role := r.Role
	if role == "" {
		role = "developer"
	}

	return &Developer{
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: passwordHash,
		Role:         role,
		Status:       "active",
	}
}
