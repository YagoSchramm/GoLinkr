package repository

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/google/uuid"
)

type AnalyticsRepository interface {
	GetAnalyticsByLinkById(ctx context.Context, linkID string, userID uuid.UUID) (*entity.Analytics, error)
	UpdateAnalytics(ctx context.Context, updatedAnalytics entity.Analytics) (*entity.Analytics, error)
}
