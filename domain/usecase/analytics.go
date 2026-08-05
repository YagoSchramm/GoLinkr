package usecase

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/google/uuid"
)

type AnalyticsUseCase interface {
	GetByLinkId(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) (*entity.Analytics, error)
	ListHourlyClickAverages(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) ([]entity.HourlyClickAverage, error)
	ListWeekdayClickAverages(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) ([]entity.WeekdayClickAverage, error)
	UpdateAnalytics(ctx context.Context, updatedAnalytics entity.Analytics) (*entity.Analytics, error)
}
