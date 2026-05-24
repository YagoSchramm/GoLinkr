package repository

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
)

type AnalyticsRepository interface {
	GetAnalyticsByLinkById(ctx context.Context, linkID string) (*entity.Analytics, error)
	UpdateAnalytics(ctx context.Context, updatedAnalytics entity.Analytics) (*entity.Analytics, error)
}
