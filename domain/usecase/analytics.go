package usecase

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
)

type AnalyticsUseCase interface {
	GetByLinkId(ctx context.Context, link_id string) (*entity.LinkAnalyticsResponse, error)
}
