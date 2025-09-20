package articles

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(req ArticleListRequest) ([]Article, int64, error)
	GetByID(id uint) (*Article, error)
	GetBySlug(slug string) (*Article, error)
	Create(article *Article) error
	Update(article *Article) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(req ArticleListRequest) ([]Article, int64, error) {
	var articles []Article
	var total int64

	query := r.db.Model(&Article{})

	// Apply filters
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.Tag != "" {
		query = query.Where("? = ANY(tags)", req.Tag)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Limit(req.PageSize).Offset(offset).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *repository) GetByID(id uint) (*Article, error) {
	var article Article
	if err := r.db.First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *repository) GetBySlug(slug string) (*Article, error) {
	var article Article
	if err := r.db.Where("slug = ?", slug).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *repository) Create(article *Article) error {
	return r.db.Create(article).Error
}

func (r *repository) Update(article *Article) error {
	return r.db.Save(article).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Article{}, id).Error
}
