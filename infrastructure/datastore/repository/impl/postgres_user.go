package impl

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
	"github.com/google/uuid"
)

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *sql.DB
}

//go:embed _query/user/GetUserByEmail.sql
var getUserByEmailQuery string

//go:embed _query/user/attemptRegister.sql
var attemptRegisterQuery string

//go:embed _query/user/attemptLogin.sql
var attemptLoginQuery string

func (r *userRepository) AttemptRegister(ctx context.Context, user entity.User) (*uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(
		ctx,
		attemptRegisterQuery,
		user.Name,
		user.Email,
		user.Password,
	).Scan(&id)
	if err != nil {
		return nil, derr.JoinError("failed to execute the query", err)
	}

	return &id, nil
}

func (r *userRepository) AttemptLogin(ctx context.Context, credentials entity.UserCredentials) (*entity.User, error) {
	var user entity.User

	row := r.db.QueryRowContext(ctx, attemptLoginQuery, credentials.Email)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, derr.JoinError("failed to execute the query", err)
	}

	return &user, err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx, getUserByEmailQuery, email)

	var user entity.User
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to get user by email", err)
	}

	return &user, nil
}
