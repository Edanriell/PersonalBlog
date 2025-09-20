package articles

import (
	"errors"
	"math"

	"Server/internal/utils"
	"Server/pkg/response"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Service interface {
	GetAll(req ArticleListRequest) ([]ArticleResponse, response.PaginationMeta, error)
	GetByID(id uint) (*ArticleResponse, error)
	GetBySlug(slug string) (*ArticleResponse, error)
	Create(authorID uint, req CreateArticleRequest) (*ArticleResponse, error)
	Update(id uint, authorID uint, req UpdateArticleRequest) (*ArticleResponse, error)
	Delete(id uint, authorID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(req ArticleListRequest) ([]ArticleResponse, response.PaginationMeta, error) {
	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	articles, total, err := s.repo.GetAll(req)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	articleResponses := make([]ArticleResponse, len(articles))
	for i, article := range articles {
		articleResponses[i] = *s.toArticleResponse(&article)
	}

	pages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	pagination := response.PaginationMeta{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		Pages:    pages,
	}

	return articleResponses, pagination, nil
}

func (s *service) GetByID(id uint) (*ArticleResponse, error) {
	article, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toArticleResponse(article), nil
}

func (s *service) GetBySlug(slug string) (*ArticleResponse, error) {
	article, err := s.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.toArticleResponse(article), nil
}

func (s *service) Create(authorID uint, req CreateArticleRequest) (*ArticleResponse, error) {
	// Check if slug already exists
	_, err := s.repo.GetBySlug(req.Slug)
	if err == nil {
		return nil, errors.New("article with this slug already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Convert markdown to HTML
	contentHTML := utils.MarkdownToHTML(req.ContentMarkdown)

	article := &Article{
		Title:           req.Title,
		Slug:            req.Slug,
		AuthorID:        authorID,
		ContentMarkdown: req.ContentMarkdown,
		ContentHTML:     contentHTML,
		Excerpt:         req.Excerpt,
		Tags:            pq.StringArray(req.Tags),
		Status:          req.Status,
	}

	if err := s.repo.Create(article); err != nil {
		return nil, err
	}

	return s.toArticleResponse(article), nil
}

func (s *service) Update(id uint, authorID uint, req UpdateArticleRequest) (*ArticleResponse, error) {
	article, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user is the author
	if article.AuthorID != authorID {
		return nil, errors.New("unauthorized to update this article")
	}

	// Update fields if provided
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Slug != "" && req.Slug != article.Slug {
		// Check if new slug already exists
		_, err := s.repo.GetBySlug(req.Slug)
		if err == nil {
			return nil, errors.New("article with this slug already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		article.Slug = req.Slug
	}
	if req.ContentMarkdown != "" {
		article.ContentMarkdown = req.ContentMarkdown
		article.ContentHTML = utils.MarkdownToHTML(req.ContentMarkdown)
	}
	if req.Excerpt != "" {
		article.Excerpt = req.Excerpt
	}
	if req.Tags != nil {
		article.Tags = pq.StringArray(req.Tags)
	}
	if req.Status != "" {
		article.Status = req.Status
	}

	if err := s.repo.Update(article); err != nil {
		return nil, err
	}

	return s.toArticleResponse(article), nil
}

func (s *service) Delete(id uint, authorID uint) error {
	article, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if user is the author
	if article.AuthorID != authorID {
		return errors.New("unauthorized to delete this article")
	}

	return s.repo.Delete(id)
}

func (s *service) toArticleResponse(article *Article) *ArticleResponse {
	return &ArticleResponse{
		ID:              article.ID,
		Title:           article.Title,
		Slug:            article.Slug,
		AuthorID:        article.AuthorID,
		ContentMarkdown: article.ContentMarkdown,
		ContentHTML:     article.ContentHTML,
		Excerpt:         article.Excerpt,
		Tags:            []string(article.Tags),
		Status:          article.Status,
		CreatedAt:       article.CreatedAt,
		UpdatedAt:       article.UpdatedAt,
	}
}
