package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardani17/taskmanager/internal/config"
	"github.com/ardani17/taskmanager/internal/handlers"
	"github.com/ardani17/taskmanager/internal/middleware"
	"github.com/ardani17/taskmanager/internal/repository"
	"github.com/ardani17/taskmanager/internal/services"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	setupLogger(cfg.IsDevelopment())

	log.Info().
		Str("port", cfg.AppPort).
		Str("env", cfg.AppEnv).
		Msg("Starting TaskManager API")

	// Connect to database
	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewDeveloperRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	activityRepo := repository.NewActivityRepository(db)

	// Initialize services
	jwtService := services.NewJWTService(cfg.JWTSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(jwtService, userRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	taskHandler := handlers.NewTaskHandler(taskRepo, activityRepo)
	projectHandler := handlers.NewProjectHandler(projectRepo, activityRepo)
	activityHandler := handlers.NewActivityHandler(activityRepo)

	// Create router
	r := chi.NewRouter()

	// Setup middleware
	setupMiddleware(r, cfg, jwtService)

	// Setup routes
	setupRoutes(r, authHandler, userHandler, taskHandler, projectHandler, activityHandler, jwtService)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Str("addr", server.Addr).Msg("HTTP server listening")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	}

	log.Info().Msg("Server stopped")
}

func setupLogger(isDevelopment bool) {
	if isDevelopment {
		// Pretty console output for development
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if isDevelopment {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func setupMiddleware(r chi.Router, cfg *config.Config, jwtService *services.JWTService) {
	// Request ID
	r.Use(chiMiddleware.RequestID)

	// Real IP
	r.Use(chiMiddleware.RealIP)

	// Logger
	r.Use(middleware.Logger)

	// Recover from panics
	r.Use(chiMiddleware.Recoverer)

	// Timeout
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: Restrict in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Compress responses
	r.Use(chiMiddleware.Compress(5))

	log.Debug().Msg("Middleware configured")
}

func setupRoutes(
	r chi.Router,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
	projectHandler *handlers.ProjectHandler,
	activityHandler *handlers.ActivityHandler,
	jwtService *services.JWTService,
) {
	// Health check (public)
	r.Get("/health", handlers.Health)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// API info (public)
		r.Get("/", handlers.APIInfo)
		r.Get("/system", handlers.SystemInfo)

		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.RefreshToken)
			
			// Protected auth routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthMiddleware(jwtService))
				r.Get("/me", authHandler.Me)
				r.Post("/logout", authHandler.Logout)
			})
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtService))

			// Users
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.List)
				r.Get("/{id}", userHandler.Get)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
				r.Patch("/{id}/status", userHandler.UpdateStatus)
			})

			// Projects
			r.Route("/projects", func(r chi.Router) {
				r.Get("/", projectHandler.List)
				r.Post("/", projectHandler.Create)
				r.Get("/{id}", projectHandler.Get)
				r.Put("/{id}", projectHandler.Update)
				r.Delete("/{id}", projectHandler.Delete)
			})

			// Tasks
			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", taskHandler.List)
				r.Post("/", taskHandler.Create)
				r.Get("/{id}", taskHandler.Get)
				r.Put("/{id}", taskHandler.Update)
				r.Delete("/{id}", taskHandler.Delete)
				r.Patch("/{id}/status", taskHandler.UpdateStatus)
			})

			// Activity
			r.Get("/activity", activityHandler.List)
		})
	})

	log.Debug().Msg("Routes configured")
}
