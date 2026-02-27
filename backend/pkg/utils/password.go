package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default bcrypt cost (10 = ~100ms)
	DefaultCost = 12
)

// HashPassword creates a bcrypt hash from the given password
func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", ErrEmptyPassword
	}
	if len(password) < 6 {
		return "", ErrPasswordTooShort
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a password with its hash
// Returns true if the password matches the hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Password validation errors
var (
	ErrEmptyPassword     = &PasswordError{Message: "password cannot be empty"}
	ErrPasswordTooShort  = &PasswordError{Message: "password must be at least 6 characters"}
	ErrPasswordNoUpper   = &PasswordError{Message: "password must contain at least one uppercase letter"}
	ErrPasswordNoLower   = &PasswordError{Message: "password must contain at least one lowercase letter"}
	ErrPasswordNoNumber  = &PasswordError{Message: "password must contain at least one number"}
	ErrPasswordNoSpecial = &PasswordError{Message: "password must contain at least one special character"}
)

// PasswordError represents a password validation error
type PasswordError struct {
	Message string
}

func (e *PasswordError) Error() string {
	return e.Message
}

// ValidatePassword checks if password meets security requirements
// Returns nil if valid, or an error describing the issue
func ValidatePassword(password string) error {
	if len(password) == 0 {
		return ErrEmptyPassword
	}
	if len(password) < 6 {
		return ErrPasswordTooShort
	}

	// Basic validation - you can add more rules if needed
	hasUpper := false
	hasLower := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		}
	}

	// For development/testing, we only require basic length
	// In production, you might want to enforce all rules
	if !hasUpper && !hasLower {
		return ErrPasswordNoLower
	}

	return nil
}

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	// Basic email validation
	atCount := 0
	for _, char := range email {
		if char == '@' {
			atCount++
		}
	}
	return atCount == 1 && len(email) > 3 && email[0] != '@'
}
