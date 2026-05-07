package repository

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	AttemptRegister(ctx context.Context, user entity.User) (uuid.UUID, error)
}
