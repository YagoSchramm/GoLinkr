package impl

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository"
)

func NewAnalyticsRepository(db *sql.DB) repository.AnalyticsRepository {
	return &analyticsRepository{
		db: db,
	}
}

type analyticsRepository struct {
	db *sql.DB
}

//go:embed _query/analytics/getAnalyticsByLinkId.sql
var getAnalyticsByLinkIdQuery string

//go:embed _query/analytics/updateAnalytics.sql
var updateAnalyticsQuery string

func (a analyticsRepository) GetAnalyticsByLinkById(ctx context.Context, linkID string) (*entity.LinkAnalyticsResponse, error) {
	var result entity.LinkAnalyticsResponse
	row := a.db.QueryRowContext(
		ctx,
		getLinkByCodeQuery,
		linkID,
	)
	if err := row.Scan(
		&result.LinkID,
		&result.Code,
		&result.OriginalURL,
		&result.CreatedAt,
		&result.Clicks,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to scan the link", err)
	}
	return &result, nil
}

func (a analyticsRepository) UpdateAnalytics(ctx context.Context, updatedAnalytics entity.Analytics) error {
	//TODO implement me
	panic("implement me")
}
