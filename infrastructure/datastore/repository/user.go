package repository

import (
	"context"

	"github.com/YagoSchramm/Golinkr/domain/entity"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
