package impl

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
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
