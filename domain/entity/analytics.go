package entity

import (
	"github.com/google/uuid"
)

type Analytics struct {
	ID     uuid.UUID `json:"id,omitempty"`
	LinkID uuid.UUID `json:"link_id,omitempty"`
	Clicks int64     `json:"clicks"`
}
