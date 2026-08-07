package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/google/uuid"
)

type fakeLinkRepository struct {
	links     map[string]entity.Link
	saveCalls int
	listCalls int
}

func newFakeLinkRepository() *fakeLinkRepository {
	return &fakeLinkRepository{
		links: make(map[string]entity.Link),
	}
}

func (r *fakeLinkRepository) Save(ctx context.Context, link entity.Link) (*entity.Link, error) {
	r.saveCalls++
	if link.ID == uuid.Nil {
		link.ID = uuid.New()
	}
	r.links[link.Code] = link
	return &link, nil
}

func (r *fakeLinkRepository) FindByCode(ctx context.Context, code string) (*entity.Link, error) {
	link, ok := r.links[code]
	if !ok {
		return nil, derr.NotFoundError
	}
	return &link, nil
}

func (r *fakeLinkRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Link, error) {
	r.listCalls++
	links := make([]entity.Link, 0)
	for _, link := range r.links {
		if link.UserId == userID {
			links = append(links, link)
		}
	}
	return links, nil
}

func TestLinkUsecaseCreateGeneratesCodeAndSaves(t *testing.T) {
	repository := newFakeLinkRepository()
	usecase := NewLinkUsecase(repository)
	userID := uuid.New()

	link, err := usecase.Create(context.Background(), entity.Link{
		UserId:      userID,
		OriginalURL: "https://example.com/docs",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if repository.saveCalls != 1 {
		t.Fatalf("expected Save to be called once, got %d", repository.saveCalls)
	}
	if link.Code == "" {
		t.Fatal("expected generated code")
	}
	if len(link.Code) != 6 {
		t.Fatalf("expected generated code with 6 chars, got %q", link.Code)
	}
	if link.UserId != userID {
		t.Fatalf("expected user id %s, got %s", userID, link.UserId)
	}
}

func TestLinkUsecaseCreateKeepsCustomCode(t *testing.T) {
	repository := newFakeLinkRepository()
	usecase := NewLinkUsecase(repository)

	link, err := usecase.Create(context.Background(), entity.Link{
		UserId:      uuid.New(),
		Code:        "custom",
		OriginalURL: "example.com",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if link.Code != "custom" {
		t.Fatalf("expected custom code to be preserved, got %q", link.Code)
	}
}

func TestLinkUsecaseCreateRejectsUnauthorizedUser(t *testing.T) {
	repository := newFakeLinkRepository()
	usecase := NewLinkUsecase(repository)

	_, err := usecase.Create(context.Background(), entity.Link{
		OriginalURL: "https://example.com",
	})
	if !errors.Is(err, derr.UnauthorizedError) {
		t.Fatalf("expected UnauthorizedError, got %v", err)
	}
	if repository.saveCalls != 0 {
		t.Fatalf("expected Save not to be called, got %d calls", repository.saveCalls)
	}
}

func TestLinkUsecaseListByUserIDRejectsUnauthorizedUser(t *testing.T) {
	repository := newFakeLinkRepository()
	usecase := NewLinkUsecase(repository)

	_, err := usecase.ListByUserID(context.Background(), uuid.Nil)
	if !errors.Is(err, derr.UnauthorizedError) {
		t.Fatalf("expected UnauthorizedError, got %v", err)
	}
	if repository.listCalls != 0 {
		t.Fatalf("expected ListByUserID not to be called, got %d calls", repository.listCalls)
	}
}

func TestLinkUsecaseFindByCodeReturnsRepositoryLink(t *testing.T) {
	repository := newFakeLinkRepository()
	usecase := NewLinkUsecase(repository)
	expected := entity.Link{
		ID:          uuid.New(),
		UserId:      uuid.New(),
		Code:        "abc123",
		OriginalURL: "https://example.com",
	}
	repository.links[expected.Code] = expected

	link, err := usecase.FindByCode(context.Background(), expected.Code)
	if err != nil {
		t.Fatalf("FindByCode returned error: %v", err)
	}
	if *link != expected {
		t.Fatalf("expected link %+v, got %+v", expected, *link)
	}
}
