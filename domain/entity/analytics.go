package entity

import (
	"github.com/google/uuid"
)

type Analytics struct {
	ID     uuid.UUID `json:"id,omitempty"`
	LinkID uuid.UUID `json:"link_id,omitempty"`
	Clicks int64     `json:"clicks"`
}

type HourlyClickAverage struct {
	Hour          int     `json:"hour"`
	AverageClicks float64 `json:"average_clicks"`
}
