package articles

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Article struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	Title           string         `json:"title" gorm:"not null"`
	Slug            string         `json:"slug" gorm:"unique;not null"`
	AuthorID        uint           `json:"author_id"`
	ContentMarkdown string         `json:"content_markdown" gorm:"type:text;not null"`
	ContentHTML     string         `json:"content_html" gorm:"type:text"`
	Excerpt         string         `json:"excerpt" gorm:"type:text"`
	Tags            pq.StringArray `json:"tags" gorm:"type:text[]"`
	Status          string         `json:"status" gorm:"default:draft"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateArticleRequest struct {
	Title           string   `json:"title" validate:"required"`
	Slug            string   `json:"slug" validate:"required"`
	ContentMarkdown string   `json:"content_markdown" validate:"required"`
	Excerpt         string   `json:"excerpt"`
	Tags            []string `json:"tags"`
	Status          string   `json:"status" validate:"required,oneof=draft published archived"`
}

type UpdateArticleRequest struct {
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	ContentMarkdown string   `json:"content_markdown"`
	Excerpt         string   `json:"excerpt"`
	Tags            []string `json:"tags"`
	Status          string   `json:"status" validate:"omitempty,oneof=draft published archived"`
}

type ArticleResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	AuthorID        uint      `json:"author_id"`
	ContentMarkdown string    `json:"content_markdown"`
	ContentHTML     string    `json:"content_html"`
	Excerpt         string    `json:"excerpt"`
	Tags            []string  `json:"tags"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ArticleListRequest struct {
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"page_size" validate:"min=1,max=100"`
	Status   string `query:"status"`
	Tag      string `query:"tag"`
}
