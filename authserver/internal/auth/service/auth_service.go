package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simple-auth-roles/internal/auth/repository"
	"github.com/simple-auth-roles/internal/config"
	"github.com/simple-auth-roles/internal/types"
	"github.com/simple-auth-roles/pkg/cache"
	"github.com/simple-auth-roles/pkg/email"
)

type AuthService struct {
	userRepo        *repository.UserRepository
	cacheService    cache.CacheService
	emailService    email.EmailService
	logger          *slog.Logger
	jwtSecret       string
	jwtExpiry       time.Duration
	webauthnService *WebAuthnService
}

func NewAuthService(userRepo *repository.UserRepository, cacheService cache.CacheService, emailService email.EmailService, logger *slog.Logger, cfg *config.Config) *AuthService {
	webAuthnService := NewWebAuthnService(userRepo, cacheService, logger, cfg)

	return &AuthService{
		userRepo:        userRepo,
		cacheService:    cacheService,
		emailService:    emailService,
		logger:          logger.With("service", "auth"),
		jwtSecret:       cfg.JWT.Secret,
		jwtExpiry:       cfg.JWT.Expiration,
		webauthnService: webAuthnService,
	}
}

func (s *AuthService) WebAuthnService() *WebAuthnService {
	return s.webauthnService
}

func (s *AuthService) UserRepository() *repository.UserRepository {
	return s.userRepo
}

// SendLoginCode creates user if needed and generates a login code
func (s *AuthService) SendLoginCode(ctx context.Context, email string, name string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Create user if not found
	if user == nil {
		user = &types.User{
			Email:    email,
			Name:     name,           // Use provided name for new users
			Role:     types.RoleUser, // Default role
			IsActive: true,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate magic link code
	code := s.generateCode()

	// Store code in cache with 10 minute expiration
	cacheKey := fmt.Sprintf("login_code:%s", strings.ToLower(email))
	if err := s.cacheService.Set(ctx, cacheKey, code, 10*time.Minute); err != nil {
		s.logger.Error("Failed to store login code", "error", err, "email", email)
		return fmt.Errorf("failed to store login code: %w", err)
	}

	// Send email with magic link code
	if err := s.emailService.SendLoginCodeEmail(ctx, email, code); err != nil {
		s.logger.Error("Failed to send login code email", "error", err, "email", email)
		return fmt.Errorf("failed to send login code email: %w", err)
	}

	s.logger.Info("Login code sent", "email", email)
	return nil
}

// VerifyLoginCode verifies the login code and returns a JWT token
func (s *AuthService) VerifyLoginCode(ctx context.Context, email, code string) (*types.AuthResponse, error) {
	// Verify code from cache
	cacheKey := fmt.Sprintf("login_code:%s", strings.ToLower(email))
	storedCode, err := s.cacheService.Get(ctx, cacheKey)
	if err != nil || storedCode == "" {
		s.logger.Warn("Login code not found or expired", "email", email)
		return nil, fmt.Errorf("code expired or not found")
	}

	if storedCode != code {
		s.logger.Warn("Invalid login code provided", "email", email)
		return nil, fmt.Errorf("invalid code")
	}

	// Delete the used code
	_ = s.cacheService.Delete(ctx, cacheKey)

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("User authenticated successfully", "email", email, "user_id", user.ID)

	return &types.AuthResponse{
		User:    user,
		Token:   token,
		Message: "Authentication successful",
	}, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *AuthService) ValidateToken(tokenString string) (*types.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		userID := uint((*claims)["user_id"].(float64))
		user, err := s.userRepo.FindByID(context.Background(), userID)
		if err != nil {
			return nil, fmt.Errorf("failed to find user: %w", err)
		}
		return user, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// GetAllUsers returns all users (admin only)
func (s *AuthService) GetAllUsers(ctx context.Context) ([]*types.User, error) {
	return s.userRepo.FindAll(ctx)
}

// UpdateUserRole updates a user's role (admin only)
func (s *AuthService) UpdateUserRole(ctx context.Context, userID uint, role string) error {
	if !types.ValidateRole(role) {
		return fmt.Errorf("invalid role: %s", role)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.Role = role
	return s.userRepo.Update(ctx, user)
}

// CreateUser creates a new user with specified role (admin only)
func (s *AuthService) CreateUser(ctx context.Context, req *types.CreateUserRequest) (*types.User, error) {
	// Check if user already exists
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = types.RoleUser
	}
	if !types.ValidateRole(role) {
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	user := &types.User{
		Email:    req.Email,
		Name:     req.Name,
		Company:  req.Company,
		Role:     role,
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GenerateJWT creates a JWT token for the user (public method)
func (s *AuthService) GenerateJWT(user *types.User) (string, error) {
	return s.generateJWT(user)
}

// generateJWT creates a JWT token for the user
func (s *AuthService) generateJWT(user *types.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateCode creates a random 6-character code
func (s *AuthService) generateCode() string {
	b := make([]byte, 5)
	_, _ = rand.Read(b)
	return strings.ToUpper(base32.StdEncoding.EncodeToString(b))[:6]
}
