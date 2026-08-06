package usecase

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/domain/rules"
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
	"github.com/google/uuid"
)

func NewAnalyticsUseCase(repository repository.AnalyticsRepository) usecase.AnalyticsUseCase {
	return &analyticsUseCase{
		repository: repository,
	}
}

type analyticsUseCase struct {
	repository repository.AnalyticsRepository
}

func (u analyticsUseCase) UpdateAnalytics(ctx context.Context, updatedAnalytics entity.Analytics) (*entity.Analytics, error) {
	if err := rules.ValidateAnalytics(updatedAnalytics); err != nil {
		return nil, err
	}

	return u.repository.UpdateAnalytics(ctx, updatedAnalytics)
}

func (u analyticsUseCase) GetByLinkId(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) (*entity.Analytics, error) {
	if linkID == uuid.Nil {
		return nil, derr.InvalidLinkId
	}
	if userID == uuid.Nil {
		return nil, derr.UnauthorizedError
	}

	return u.repository.GetAnalyticsByLinkById(ctx, linkID, userID)
}

func (u analyticsUseCase) ListHourlyClickAverages(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) ([]entity.HourlyClickAverage, error) {
	if err := rules.ValidateHourlyClickAverage(linkID, userID); err != nil {
		return nil, err
	}

	return u.repository.ListHourlyClickAverages(ctx, linkID, userID)
}

func (u analyticsUseCase) ListMonthlyWeekClickAverages(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) ([]entity.MonthlyWeekClickAverage, error) {
	if err := rules.ValidateMonthlyWeekClickAverage(linkID, userID); err != nil {
		return nil, err
	}

	return u.repository.ListMonthlyWeekClickAverages(ctx, linkID, userID)
}

func (u analyticsUseCase) ListWeekdayClickAverages(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) ([]entity.WeekdayClickAverage, error) {
	if err := rules.ValidateWeekdayClickAverage(linkID, userID); err != nil {
		return nil, err
	}

	return u.repository.ListWeekdayClickAverages(ctx, linkID, userID)
}
