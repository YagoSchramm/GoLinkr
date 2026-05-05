package impl

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
)

func NewLinkRepository(db *sql.DB) repository.LinkRepository {
	return &linkRepository{
		db: db,
	}
}

type linkRepository struct {
	db *sql.DB
}

//go:embed _query/link/CreateLink.sql
var createLinkQuery string

//go:embed _query/link/GetLinkByCode.sql
var getLinkByCodeQuery string

func (r *linkRepository) FindByCode(ctx context.Context, code string) (*entity.Link, error) {
	row := r.db.QueryRowContext(
		ctx,
		getLinkByCodeQuery,
		code,
	)

	var link entity.Link
	if err := row.Scan(
		&link.ID,
		&link.Code,
		&link.OriginalURL,
		&link.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to scan the link", err)
	}

	return &link, nil
}

func (r *linkRepository) Save(ctx context.Context, link entity.Link) (*entity.Link, error) {
	var result entity.Link
	err := r.db.QueryRowContext(
		ctx,
		createLinkQuery,
		link.Code,
		link.OriginalURL,
	).Scan(&result.ID, &result.CreatedAt)
	if err != nil {
		return nil, derr.JoinError("failed to execute the query", err)
	}

	result.Code = link.Code
	result.OriginalURL = link.OriginalURL

	return &result, nil
}
