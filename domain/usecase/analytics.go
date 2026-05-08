package usecase

import "context"

type AnalyticsUseCase interface {
	GetByLinkId(ctx context.Context, link_id string)
}
