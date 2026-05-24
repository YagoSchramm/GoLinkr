package entity

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}
