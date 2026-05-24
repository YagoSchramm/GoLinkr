package usecase

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
)

func NewAnalytics(repository repository.AnalyticsRepository) usecase.AnalyticsUseCase {
	return &analyticsUsecase{
		repository: repository,
	}
}

type analyticsUsecase struct {
	repository repository.AnalyticsRepository
}

func (u analyticsUsecase) GetByLinkId(ctx context.Context, link_id string) (*entity.LinkAnalyticsResponse, error) {
	return u.repository.GetAnalyticsByLinkById(ctx, link_id)
}
