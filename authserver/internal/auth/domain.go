package auth

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/simple-auth-roles/internal/auth/handlers"
	"github.com/simple-auth-roles/internal/auth/repository"
	"github.com/simple-auth-roles/internal/auth/service"
	"github.com/simple-auth-roles/pkg/cache"
	"github.com/simple-auth-roles/pkg/email"
	"gorm.io/gorm"
)

// Domain represents the authentication domain
type Domain struct {
	service *service.AuthService
	handler *handlers.AuthHandler
	logger  *slog.Logger
}

// NewDomain creates a new authentication domain
func NewDomain(db *gorm.DB, cacheService cache.CacheService, emailService email.EmailService, logger *slog.Logger) *Domain {
	// Create repository
	userRepo := repository.NewUserRepository(db)

	// Create service
	authService := service.NewAuthService(userRepo, cacheService, emailService, logger)

	// Create handler
	authHandler := handlers.NewAuthHandler(authService, logger)

	return &Domain{
		service: authService,
		handler: authHandler,
		logger:  logger.With("domain", "auth"),
	}
}

// Service returns the auth service
func (d *Domain) Service() *service.AuthService {
	return d.service
}

// RegisterRoutes registers authentication routes
func (d *Domain) RegisterRoutes(router *gin.RouterGroup) {
	d.handler.RegisterRoutes(router)
}
