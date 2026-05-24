package usecase

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
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
	return u.repository.UpdateAnalytics(ctx, updatedAnalytics)
}

func (u analyticsUseCase) GetByLinkId(ctx context.Context, link_id string) (*entity.Analytics, error) {
	return u.repository.GetAnalyticsByLinkById(ctx, link_id)
}
