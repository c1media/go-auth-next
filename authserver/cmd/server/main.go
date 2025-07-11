package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/simple-auth-roles/internal/auth"
	"github.com/simple-auth-roles/internal/config"
	"github.com/simple-auth-roles/internal/middleware"
	"github.com/simple-auth-roles/pkg/cache"
	"github.com/simple-auth-roles/pkg/database"
	"github.com/simple-auth-roles/pkg/email"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	config     *config.Config
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// .env file is optional, continue without it
		slog.Debug("No .env file found, using system environment variables")
	}

	// Parse command line flags
	var (
		runMigrations = flag.Bool("migrate", false, "Run migrations before starting server")
		migrateOnly   = flag.Bool("migrate-only", false, "Run migrations only and exit")
		seedOnly      = flag.Bool("seed-only", false, "Run admin seeding only and exit")
		showHelp      = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	// Check for environment variable override
	if os.Getenv("MIGRATE") == "true" {
		*runMigrations = true
	}

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Initialize database
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Configure database connection pool
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get database instance", "error", err)
		os.Exit(1)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)

	// Run migrations if requested
	if *runMigrations || *migrateOnly {
		logger.Info("Running database migrations...")
		if err := database.RunMigrations(db); err != nil {
			logger.Error("Failed to run migrations", "error", err)
			os.Exit(1)
		}
		logger.Info("Database migrations completed successfully")

		// Seed admin user on first migration
		logger.Info("Checking for admin user...")
		if err := database.SeedAdminUser(db); err != nil {
			logger.Error("Failed to seed admin user", "error", err)
			os.Exit(1)
		}

		// Exit if migrate-only flag is set
		if *migrateOnly {
			logger.Info("Migration-only mode: exiting after migrations")
			os.Exit(0)
		}
	}

	// Run admin seeding only if requested
	if *seedOnly {
		logger.Info("Running admin user seeding...")
		if err := database.SeedAdminUser(db); err != nil {
			logger.Error("Failed to seed admin user", "error", err)
			os.Exit(1)
		}
		logger.Info("Admin seeding completed, exiting...")
		os.Exit(0)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize services
	cacheService := cache.NewCacheService(cfg, logger)
	emailService := email.NewEmailService(cfg, logger)

	// Initialize auth domain
	authDomain := auth.NewDomain(db, cacheService, emailService, logger, cfg)

	// Setup router
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.CORSMiddleware(cfg))

	// Setup routes
	setupRoutes(router, authDomain)

	// Create server
	server := &Server{
		httpServer: &http.Server{
			Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			Handler:        router,
			ReadTimeout:    cfg.Server.ReadTimeout,
			WriteTimeout:   cfg.Server.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
		logger: logger,
		config: cfg,
	}

	// Start server
	go func() {
		logger.Info("Starting server", "addr", server.httpServer.Addr)
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}

func setupRoutes(router *gin.Engine, authDomain *auth.Domain) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"time":        time.Now().UTC(),
			"api_version": "v1",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	authDomain.RegisterRoutes(api)

	// Protected demo routes
	protected := api.Group("/protected")
	protected.Use(middleware.RequireAuth(authDomain.Service()))
	{
		protected.GET("/profile", func(c *gin.Context) {
			user := middleware.GetCurrentUser(c)
			c.JSON(http.StatusOK, gin.H{
				"user":    user,
				"message": "This is a protected route",
			})
		})

		// Admin only routes
		adminOnly := protected.Group("/admin")
		adminOnly.Use(middleware.RequireRole("admin"))
		{
			adminOnly.GET("/users", func(c *gin.Context) {
				users, err := authDomain.Service().GetAllUsers(c.Request.Context())
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"users": users})
			})
		}
	}
}
