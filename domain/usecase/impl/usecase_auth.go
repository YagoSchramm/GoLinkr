package usecase

import (
	"context"
	"errors"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/domain/rules"
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
	service "github.com/YagoSchramm/Golinkr/infrastructure/service/hash"
	service2 "github.com/YagoSchramm/Golinkr/infrastructure/service/jwt"
	"github.com/google/uuid"
)

func NewAuthUsecase(repository repository.UserRepository, secret string) usecase.AuthUseCase {
	return &authUsecase{
		repository: repository,
		secret:     secret,
	}
}

type authUsecase struct {
	repository repository.UserRepository
	secret     string
}

func (a *authUsecase) AttemptRegister(ctx context.Context, user entity.User) (string, error) {
	err := rules.ValidateRegister(user)
	if err != nil {
		return "", err
	}

	existedUser, err := a.repository.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, derr.NotFoundError) {
		return "", derr.JoinError("failed to get user by email", err)
	}

	if existedUser != nil {
		return "", derr.UserAlreadyExists
	}

	hashedPassword, err := service.HashPassword(user.Password)
	if err != nil {
		return "", derr.JoinError("failed to hash the password", err)
	}

	user.Password = hashedPassword

	id, err := a.repository.AttemptRegister(ctx, user)
	if err != nil {
		return "", derr.JoinError("failed to attempt register the user", err)
	}

	token, err := service2.GenerateToken(*id, user.Email, []byte(a.secret))
	if err != nil {
		return "", derr.JoinError("failed to generate the token", err)

	}

	return token, err
}

func (a *authUsecase) AttemptLogin(ctx context.Context, credentials entity.UserCredentials) (string, error) {
	err := rules.ValidateLogin(credentials)
	if err != nil {
		return "", err
	}

	existedUser, err := a.repository.GetUserByEmail(ctx, credentials.Email)
	if err != nil && !errors.Is(err, derr.NotFoundError) {
		return "", derr.JoinError("failed to get user by email", err)
	}

	if existedUser == nil {
		return "", derr.NotFoundError
	}

	user, err := a.repository.AttemptLogin(ctx, credentials)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", derr.NewNotFoundError("user not found")
	}

	valid := service.CheckPassword(credentials.Password, user.Password)
	if !valid {
		return "", derr.InvalidCredentials
	}

	token, err := service2.GenerateToken(existedUser.ID, credentials.Email, []byte(a.secret))
	if err != nil {
		return "", err
	}

	return token, nil
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
