package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	hashservice "github.com/YagoSchramm/Golinkr/infrastructure/service/hash"
	"github.com/google/uuid"
)

type fakeUserRepository struct {
	users         map[string]entity.User
	registerCalls int
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		users: make(map[string]entity.User),
	}
}

func (r *fakeUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, derr.NotFoundError
	}
	return &user, nil
}

func (r *fakeUserRepository) AttemptRegister(ctx context.Context, user entity.User) (*uuid.UUID, error) {
	r.registerCalls++
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	r.users[user.Email] = user
	return &user.ID, nil
}

func (r *fakeUserRepository) AttemptLogin(ctx context.Context, credentials entity.UserCredentials) (*entity.User, error) {
	user, ok := r.users[credentials.Email]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

func TestAuthUsecaseAttemptRegisterStoresHashedPasswordAndReturnsToken(t *testing.T) {
	repository := newFakeUserRepository()
	usecase := NewAuthUsecase(repository, "test-secret")

	token, err := usecase.AttemptRegister(context.Background(), entity.User{
		Name:     "Yago",
		Email:    "yago@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("AttemptRegister returned error: %v", err)
	}

	if token == "" {
		t.Fatal("expected token")
	}
	if repository.registerCalls != 1 {
		t.Fatalf("expected AttemptRegister to be called once, got %d", repository.registerCalls)
	}
	savedUser := repository.users["yago@example.com"]
	if savedUser.Password == "secret123" {
		t.Fatal("expected saved password to be hashed")
	}
	if !hashservice.CheckPassword("secret123", savedUser.Password) {
		t.Fatal("expected saved password hash to match original password")
	}
}

func TestAuthUsecaseAttemptRegisterRejectsExistingUser(t *testing.T) {
	repository := newFakeUserRepository()
	repository.users["yago@example.com"] = entity.User{
		ID:       uuid.New(),
		Name:     "Yago",
		Email:    "yago@example.com",
		Password: "already-hashed",
	}
	usecase := NewAuthUsecase(repository, "test-secret")

	_, err := usecase.AttemptRegister(context.Background(), entity.User{
		Name:     "Yago",
		Email:    "yago@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, derr.UserAlreadyExists) {
		t.Fatalf("expected UserAlreadyExists, got %v", err)
	}
	if repository.registerCalls != 0 {
		t.Fatalf("expected AttemptRegister not to be called, got %d calls", repository.registerCalls)
	}
}

func TestAuthUsecaseAttemptLoginReturnsTokenForValidCredentials(t *testing.T) {
	repository := newFakeUserRepository()
	passwordHash, err := hashservice.HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	repository.users["yago@example.com"] = entity.User{
		ID:       uuid.New(),
		Email:    "yago@example.com",
		Password: passwordHash,
	}
	usecase := NewAuthUsecase(repository, "test-secret")

	token, err := usecase.AttemptLogin(context.Background(), entity.UserCredentials{
		Email:    "yago@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("AttemptLogin returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
}

func TestAuthUsecaseAttemptLoginRejectsInvalidPassword(t *testing.T) {
	repository := newFakeUserRepository()
	passwordHash, err := hashservice.HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	repository.users["yago@example.com"] = entity.User{
		ID:       uuid.New(),
		Email:    "yago@example.com",
		Password: passwordHash,
	}
	usecase := NewAuthUsecase(repository, "test-secret")

	_, err = usecase.AttemptLogin(context.Background(), entity.UserCredentials{
		Email:    "yago@example.com",
		Password: "wrong123",
	})
	if !errors.Is(err, derr.InvalidCredentials) {
		t.Fatalf("expected InvalidCredentials, got %v", err)
	}
}

func TestAuthUsecaseValidateSessionChecksUserIDAndEmail(t *testing.T) {
	userID := uuid.New()
	repository := newFakeUserRepository()
	repository.users["yago@example.com"] = entity.User{
		ID:    userID,
		Email: "yago@example.com",
	}
	usecase := NewAuthUsecase(repository, "test-secret")

	if err := usecase.ValidateSession(context.Background(), userID, "yago@example.com"); err != nil {
		t.Fatalf("ValidateSession returned error: %v", err)
	}

	err := usecase.ValidateSession(context.Background(), uuid.New(), "yago@example.com")
	if !errors.Is(err, derr.UnauthorizedError) {
		t.Fatalf("expected UnauthorizedError, got %v", err)
	}
}
