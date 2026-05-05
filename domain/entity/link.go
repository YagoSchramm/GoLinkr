package entity

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID
	Code        string
	OriginalURL string
	CreatedAt   time.Time
}
