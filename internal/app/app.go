package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/bsnack-backend/config"
	"github.com/savanyv/bsnack-backend/internal/database"
	"github.com/savanyv/bsnack-backend/internal/middlewares"
)

type Server struct {
	app *fiber.App
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	app := fiber.New(fiber.Config{
		AppName: config.AppName,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 10 * time.Second,
	})

	return &Server{
		app: app,
		config: config,
	}
}

func (s *Server) Start() error {
	// Initialize Database
	if _, err := database.InitDatabase(s.config); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	if err := database.AutoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Middleware & Routes
	s.app.Use(middlewares.CORSMiddleware())
	s.app.Use(middlewares.MethodValidationMiddleware())

	// Start Server
	addr := fmt.Sprintf(":%s", s.config.AppPort)
	go func() {
		if err := s.app.Listen(addr); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server gracefully stopped")
	return nil
}
