package users
package newsletter

import (
"Server/internal/middleware"

"github.com/labstack/echo/v4"
"gorm.io/gorm"
)

func RegisterRoutes(g *echo.Group, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	// Public routes
	g.POST("/newsletter/subscribe", handler.Subscribe)
	g.PUT("/newsletter/unsubscribe/:id", handler.Unsubscribe)

	// Protected routes (admin only)
	protected := g.Group("/newsletter")
	protected.Use(middleware.JWTAuth())
	protected.GET("/subscribers", handler.GetAll)
	protected.GET("/subscribers/:id", handler.GetByID)
	protected.DELETE("/subscribers/:id", handler.Delete)
}