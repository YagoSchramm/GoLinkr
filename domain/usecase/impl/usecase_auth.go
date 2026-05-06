package usecase

import (
	"context"
	"errors"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
	"github.com/google/uuid"
)

func NewAuthUsecase(repository repository.UserRepository) usecase.AuthUseCase {
	return &authUsecase{
		repository: repository,
	}
}

type authUsecase struct {
	repository repository.UserRepository
}

func (a *authUsecase) Authenticate(ctx context.Context, credentials entity.UserCredentials) (*entity.User, error) {
	return nil, derr.InvalidCredentials
}

func (a *authUsecase) ValidateSession(ctx context.Context, userID uuid.UUID, email string) error {
	if userID == uuid.Nil || email == "" {
		return derr.UnauthorizedError
	}

	user, err := a.repository.GetUserByEmail(ctx, email)
	if err != nil {
		var clientErr derr.ClientError
		if errors.As(err, &clientErr) && clientErr.Code == derr.NotFoundError.Code {
			return derr.UnauthorizedError
		}
		return err
	}
	if user.ID != userID {
		return derr.UnauthorizedError
	}

	return nil
}
