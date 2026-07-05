package usecase

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/google/uuid"
)

type LinkUsecase interface {
	Create(ctx context.Context, link entity.Link) (*entity.Link, error)
	FindByCode(ctx context.Context, code string) (*entity.Link, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Link, error)
}
