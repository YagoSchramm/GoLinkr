package entity

import (
	"github.com/google/uuid"
)

type Analytics struct {
	ID     uuid.UUID `json:"id,omitempty"`
	LinkID uuid.UUID `json:"link_id,omitempty"`
	Clicks int64     `json:"clicks"`
}

type GetAnalyticsDTO struct {
	UserID uuid.UUID `json:"user_id"`
	LinkID string    `json:"link_id,omitempty"`
}
