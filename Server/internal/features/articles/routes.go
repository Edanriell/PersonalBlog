package articles

import (
	"Server/internal/middleware"

	"gorm.io/gorm"
)

func RegisterRoutes(g *echo.Group, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	// Public routes
	g.GET("/articles", handler.GetAll)
	g.GET("/articles/:id", handler.GetByID)
	g.GET("/articles/slug/:slug", handler.GetBySlug)

	// Protected routes
	protected := g.Group("/articles")
	protected.Use(middleware.JWTAuth())
	protected.POST("", handler.Create)
	protected.PUT("/:id", handler.Update)
	protected.DELETE("/:id", handler.Delete)
}
