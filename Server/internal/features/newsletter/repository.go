package newsletter

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(req SubscriberListRequest) ([]Subscriber, int64, error)
	GetByID(id uint) (*Subscriber, error)
	GetByEmail(email string) (*Subscriber, error)
	Create(subscriber *Subscriber) error
	Update(subscriber *Subscriber) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(req SubscriberListRequest) ([]Subscriber, int64, error) {
	var subscribers []Subscriber
	var total int64

	query := r.db.Model(&Subscriber{})

	// Apply filters
	if req.Active {
		query = query.Where("unsubscribed_at IS NULL")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Limit(req.PageSize).Offset(offset).Find(&subscribers).Error; err != nil {
		return nil, 0, err
	}

	return subscribers, total, nil
}

func (r *repository) GetByID(id uint) (*Subscriber, error) {
	var subscriber Subscriber
	if err := r.db.First(&subscriber, id).Error; err != nil {
		return nil, err
	}
	return &subscriber, nil
}

func (r *repository) GetByEmail(email string) (*Subscriber, error) {
	var subscriber Subscriber
	if err := r.db.Where("email = ?", email).First(&subscriber).Error; err != nil {
		return nil, err
	}
	return &subscriber, nil
}

func (r *repository) Create(subscriber *Subscriber) error {
	return r.db.Create(subscriber).Error
}

func (r *repository) Update(subscriber *Subscriber) error {
	return r.db.Save(subscriber).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Subscriber{}, id).Error
}
