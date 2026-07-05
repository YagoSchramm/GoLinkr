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

func (r analyticsRepository) GetAnalyticsByLinkById(ctx context.Context, linkID string, userID uuid.UUID) (*entity.Analytics, error) {
	var result entity.Analytics
	row := r.db.QueryRowContext(
		ctx,
		getAnalyticsByLinkIdQuery,
		linkID,
		userID,
	)
	if err := row.Scan(
		&result.ID,
		&result.LinkID,
		&result.Clicks,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to scan the link", err)
	}
	return &result, nil
}

func (r analyticsRepository) UpdateAnalytics(ctx context.Context, updatedAnalytics entity.Analytics) (*entity.Analytics, error) {
	var analytics entity.Analytics
	row := r.db.QueryRowContext(
		ctx,
		updateAnalyticsQuery,
		updatedAnalytics.LinkID,
	)
	if err := row.Scan(
		&analytics.ID,
		&analytics.LinkID,
		&analytics.Clicks,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, derr.NotFoundError
		}
		return nil, derr.JoinError("failed to scan the link", err)
	}

	return &analytics, nil
}
