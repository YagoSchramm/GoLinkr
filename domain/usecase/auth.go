package usecase

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/google/uuid"
)

type AuthUseCase interface {
	Authenticate(ctx context.Context, credentials entity.UserCredentials) (*entity.User, error)
	ValidateSession(ctx context.Context, userID uuid.UUID, email string) error
}
