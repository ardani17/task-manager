package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ardani17/taskmanager/internal/middleware"
	"github.com/ardani17/taskmanager/internal/models"
	"github.com/ardani17/taskmanager/internal/repository"
	"github.com/ardani17/taskmanager/internal/services"
	"github.com/ardani17/taskmanager/pkg/utils"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	jwtService *services.JWTService
	userRepo   *repository.DeveloperRepository
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtService *services.JWTService, userRepo *repository.DeveloperRepository) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Validate email format
	if !utils.ValidateEmail(req.Email) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Validate password strength
	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if email already exists
	existingUser, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to check email")
		return
	}
	if existingUser != nil {
		utils.ErrorResponse(w, http.StatusConflict, "Email already registered")
		return
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create developer model
	developer := req.ToDeveloper(passwordHash)

	// Save to database
	if err := h.userRepo.Create(developer); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT tokens (convert int ID to string)
	tokenPair, err := h.jwtService.GenerateToken(strconv.Itoa(developer.ID), developer.Email, developer.Role)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Update status to online
	h.userRepo.UpdateStatus(developer.ID, "online")

	// Return response
	utils.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Registration successful",
		"data": map[string]interface{}{
			"developer": developer,
			"token": map[string]interface{}{
				"access_token": tokenPair.AccessToken,
				"token_type":   tokenPair.TokenType,
				"expires_in":   tokenPair.ExpiresIn,
			},
		},
	})
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Fetch developer from database
	developer, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	if developer == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Verify password (check if password_hash column exists and has value)
	// For existing users without password, allow any password for now
	if developer.PasswordHash == "" {
		// Legacy user - set password hash
		passwordHash, _ := utils.HashPassword(req.Password)
		developer.PasswordHash = passwordHash
	}

	if !utils.CheckPassword(req.Password, developer.PasswordHash) {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT tokens
	tokenPair, err := h.jwtService.GenerateToken(strconv.Itoa(developer.ID), developer.Email, developer.Role)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Update developer status
	h.userRepo.UpdateStatus(developer.ID, "online")
	developer.Status = "online"

	// Return response
	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data": map[string]interface{}{
			"developer": developer,
			"token": map[string]interface{}{
				"access_token": tokenPair.AccessToken,
				"token_type":   tokenPair.TokenType,
				"expires_in":   tokenPair.ExpiresIn,
			},
		},
	})
}

// Me handles GET /api/v1/auth/me
// Returns the currently authenticated user
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Fetch full developer data from database
	developer, err := h.userRepo.GetByID(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	if developer == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    developer,
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := middleware.GetUserID(r)
	if userID == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Update developer status to offline
	h.userRepo.UpdateStatus(userID, "inactive")

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logout successful",
	})
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// Refresh the token
	tokenPair, err := h.jwtService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Token refreshed successfully",
		"data": map[string]interface{}{
			"token": map[string]interface{}{
				"access_token": tokenPair.AccessToken,
				"token_type":   tokenPair.TokenType,
				"expires_in":   tokenPair.ExpiresIn,
			},
		},
	})
}
