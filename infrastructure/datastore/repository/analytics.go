package repository

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
)

type AnalyticsRepository interface {
	GetAnalyticsByLinkById(ctx context.Context, linkID string) (*entity.LinkAnalyticsResponse, error)
	UpdateAnalytics(ctx context.Context, updatedAnalytics entity.Analytics) error
}
