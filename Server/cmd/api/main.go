package main

import (
	"log"
	"os"

	"Server/internal/config"
	"Server/internal/database"
	"Server/internal/features/articles"
	"Server/internal/features/newsletter"
	"Server/internal/features/users"
	"Server/internal/middleware"

	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(middleware.CORS())

	// API routes group
	api := e.Group("/api/v1")

	// Initialize features
	users.RegisterRoutes(api, db)
	articles.RegisterRoutes(api, db)
	newsletter.RegisterRoutes(api, db)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
