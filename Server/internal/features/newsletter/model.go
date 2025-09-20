package newsletter

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscriber struct {
	ID                uint           `json:"id" gorm:"primarykey"`
	Email             string         `json:"email" gorm:"unique;not null"`
	Name              string         `json:"name"`
	SubscribedAt      time.Time      `json:"subscribed_at" gorm:"default:now()"`
	UnsubscribedAt    *time.Time     `json:"unsubscribed_at"`
	IsConfirmed       bool           `json:"is_confirmed" gorm:"default:false"`
	ConfirmationToken *uuid.UUID     `json:"-" gorm:"type:uuid"`
	Source            string         `json:"source" gorm:"default:website"`
	Preferences       string         `json:"preferences" gorm:"type:jsonb"`
	GDPRConsent       bool           `json:"gdpr_consent" gorm:"default:true"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

type SubscribeRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	GDPRConsent bool   `json:"gdpr_consent"`
}

type SubscriberResponse struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	SubscribedAt time.Time `json:"subscribed_at"`
	IsConfirmed  bool      `json:"is_confirmed"`
	Source       string    `json:"source"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SubscriberListRequest struct {
	Page     int  `query:"page" validate:"min=1"`
	PageSize int  `query:"page_size" validate:"min=1,max=100"`
	Active   bool `query:"active"`
}
