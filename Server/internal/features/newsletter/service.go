package newsletter

import (
	"errors"
	"math"
	"time"

	"Server/pkg/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetAll(req SubscriberListRequest) ([]SubscriberResponse, response.PaginationMeta, error)
	GetByID(id uint) (*SubscriberResponse, error)
	Subscribe(req SubscribeRequest) (*SubscriberResponse, error)
	Unsubscribe(id uint) error
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(req SubscriberListRequest) ([]SubscriberResponse, response.PaginationMeta, error) {
	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	subscribers, total, err := s.repo.GetAll(req)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	subscriberResponses := make([]SubscriberResponse, len(subscribers))
	for i, subscriber := range subscribers {
		subscriberResponses[i] = *s.toSubscriberResponse(&subscriber)
	}

	pages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	pagination := response.PaginationMeta{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		Pages:    pages,
	}

	return subscriberResponses, pagination, nil
}

func (s *service) GetByID(id uint) (*SubscriberResponse, error) {
	subscriber, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toSubscriberResponse(subscriber), nil
}

func (s *service) Subscribe(req SubscribeRequest) (*SubscriberResponse, error) {
	// Check if subscriber already exists
	existingSubscriber, err := s.repo.GetByEmail(req.Email)
	if err == nil {
		// If exists and not unsubscribed, return existing
		if existingSubscriber.UnsubscribedAt == nil {
			return s.toSubscriberResponse(existingSubscriber), nil
		}
		// If previously unsubscribed, resubscribe
		existingSubscriber.UnsubscribedAt = nil
		existingSubscriber.IsConfirmed = false
		token := uuid.New()
		existingSubscriber.ConfirmationToken = &token

		if err := s.repo.Update(existingSubscriber); err != nil {
			return nil, err
		}

		return s.toSubscriberResponse(existingSubscriber), nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new subscriber
	token := uuid.New()
	subscriber := &Subscriber{
		Email:             req.Email,
		Name:              req.Name,
		Source:            req.Source,
		GDPRConsent:       req.GDPRConsent,
		IsConfirmed:       false,
		ConfirmationToken: &token,
		SubscribedAt:      time.Now(),
	}

	if subscriber.Source == "" {
		subscriber.Source = "website"
	}

	if err := s.repo.Create(subscriber); err != nil {
		return nil, err
	}

	return s.toSubscriberResponse(subscriber), nil
}

func (s *service) Unsubscribe(id uint) error {
	subscriber, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	subscriber.UnsubscribedAt = &now

	return s.repo.Update(subscriber)
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) toSubscriberResponse(subscriber *Subscriber) *SubscriberResponse {
	return &SubscriberResponse{
		ID:           subscriber.ID,
		Email:        subscriber.Email,
		Name:         subscriber.Name,
		SubscribedAt: subscriber.SubscribedAt,
		IsConfirmed:  subscriber.IsConfirmed,
		Source:       subscriber.Source,
		CreatedAt:    subscriber.CreatedAt,
		UpdatedAt:    subscriber.UpdatedAt,
	}
}
