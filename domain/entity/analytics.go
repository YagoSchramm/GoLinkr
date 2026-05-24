package entity

import (
	"time"

	"github.com/google/uuid"
)

type Analytics struct {
	ID     uuid.UUID `json:"id,omitempty"`
	LinkID uuid.UUID `json:"link_id,omitempty"`
	Clicks int64     `json:"clicks"`
}

type GetAnalyticsDTO struct {
	LinkID string `json:"link_id,omitempty"`
}

type LinkAnalyticsResponse struct {
	LinkID      uuid.UUID `json:"link_id" db:"link_id"`
	Code        string    `json:"code" db:"code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Clicks      int64     `json:"clicks" db:"clicks"`
}
