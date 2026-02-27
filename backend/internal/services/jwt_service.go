package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token operations
type JWTService struct {
	secretKey     string
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
}

// JWTClaims represents the custom claims in the JWT token
type JWTClaims struct {
	DeveloperID string `json:"developer_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair contains access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"`
}

// JWT errors
var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("token has expired")
	ErrInvalidClaims     = errors.New("invalid token claims")
	ErrMissingSecretKey  = errors.New("JWT secret key is required")
)

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string) *JWTService {
	if secretKey == "" {
		panic(ErrMissingSecretKey)
	}

	return &JWTService{
		secretKey:     secretKey,
		tokenExpiry:   24 * time.Hour,       // 24 hours
		refreshExpiry: 7 * 24 * time.Hour,   // 7 days
	}
}

// NewJWTServiceWithExpiry creates a new JWT service with custom expiry times
func NewJWTServiceWithExpiry(secretKey string, tokenExpiry, refreshExpiry time.Duration) *JWTService {
	if secretKey == "" {
		panic(ErrMissingSecretKey)
	}

	return &JWTService{
		secretKey:     secretKey,
		tokenExpiry:   tokenExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateToken generates a new access token for a developer
func (s *JWTService) GenerateToken(developerID, email, role string) (*TokenPair, error) {
	now := time.Now()

	// Access token
	accessClaims := JWTClaims{
		DeveloperID: developerID,
		Email:       email,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   developerID,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token
	refreshClaims := JWTClaims{
		DeveloperID: developerID,
		Email:       email,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   developerID,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.tokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken validates an access token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *JWTService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Generate new token pair
	return s.GenerateToken(claims.DeveloperID, claims.Email, claims.Role)
}

// ExtractToken extracts the JWT token from the Authorization header
// Expected format: "Bearer <token>"
func ExtractToken(authHeader string) (string, error) {
	if len(authHeader) < 7 {
		return "", ErrInvalidToken
	}

	if authHeader[:7] != "Bearer " {
		return "", ErrInvalidToken
	}

	return authHeader[7:], nil
}
