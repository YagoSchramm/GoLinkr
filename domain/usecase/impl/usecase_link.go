package usecase

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/domain/rules"
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
	"github.com/google/uuid"
)

func NewLinkUsecase(repository repository.LinkRepository) usecase.LinkUsecase {
	return &linkUsecase{
		repository: repository,
	}
}

type linkUsecase struct {
	repository repository.LinkRepository
}

func (l *linkUsecase) Create(ctx context.Context, link entity.Link) (*entity.Link, error) {
	if link.UserId == uuid.Nil {
		return nil, derr.UnauthorizedError
	}

	if err := rules.ValidateURL(link.OriginalURL); err != nil {
		return nil, err
	}

	if link.Code == "" {
		code, err := generateCode(6)
		if err != nil {
			return nil, err
		}
		link.Code = code
	}

	return l.repository.Save(ctx, link)
}

func (l *linkUsecase) FindByCode(ctx context.Context, code string) (*entity.Link, error) {
	return l.repository.FindByCode(ctx, code)
}

func (l *linkUsecase) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Link, error) {
	if userID == uuid.Nil {
		return nil, derr.UnauthorizedError
	}

	return l.repository.ListByUserID(ctx, userID)
}

func generateCode(size int) (string, error) {
	if size <= 0 {
		size = 6
	}

	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	max := big.NewInt(int64(len(alphabet)))
	out := make([]byte, size)

	for i := 0; i < size; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", derr.JoinError("failed to generate link code", err)
		}
		out[i] = alphabet[n.Int64()]
	}

	return string(out), nil
}
