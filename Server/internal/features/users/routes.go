package users

import (
	"Server/internal/middleware"

	"gorm.io/gorm"
)

func RegisterRoutes(g *echo.Group, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	// Public routes
	g.POST("/auth/register", handler.Register)
	g.POST("/auth/login", handler.Login)

	// Protected routes
	protected := g.Group("/users")
	protected.Use(middleware.JWTAuth())
	protected.GET("/:id", handler.GetByID)
	protected.PUT("/profile", handler.Update)
}
