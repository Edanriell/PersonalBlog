package database

import (
	"Server/internal/features/articles"
	"Server/internal/features/newsletter"
	"Server/internal/features/users"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Run migrations for all models
	return db.AutoMigrate(
		&users.User{},
		&articles.Article{},
		&newsletter.Subscriber{},
	)
}
