package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/google/uuid"
)

type fakeAnalyticsRepository struct {
	analytics                  map[uuid.UUID]entity.Analytics
	hourlyClickAverages        []entity.HourlyClickAverage
	monthlyWeekClickAverages  []entity.MonthlyWeekClickAverage
	weekdayClickAverages      []entity.WeekdayClickAverage
	updateCalls                int
	lastLinkID                 uuid.UUID
	lastUserID                 uuid.UUID
}

func newFakeAnalyticsRepository() *fakeAnalyticsRepository {
	return &fakeAnalyticsRepository{
		analytics: make(map[uuid.UUID]entity.Analytics),
	}
}

func (r *fakeAnalyticsRepository) GetAnalyticsByLinkById(
	ctx context.Context,
	linkID uuid.UUID,
	userID uuid.UUID,
) (*entity.Analytics, error) {
	r.lastLinkID = linkID
	r.lastUserID = userID
	analytics, ok := r.analytics[linkID]
	if !ok {
		return nil, derr.NotFoundError
	}
	return &analytics, nil
}

func (r *fakeAnalyticsRepository) ListHourlyClickAverages(
	ctx context.Context,
	linkID uuid.UUID,
	userID uuid.UUID,
) ([]entity.HourlyClickAverage, error) {
	r.lastLinkID = linkID
	r.lastUserID = userID
	return r.hourlyClickAverages, nil
}

func (r *fakeAnalyticsRepository) ListMonthlyWeekClickAverages(
	ctx context.Context,
	linkID uuid.UUID,
	userID uuid.UUID,
) ([]entity.MonthlyWeekClickAverage, error) {
	r.lastLinkID = linkID
	r.lastUserID = userID
	return r.monthlyWeekClickAverages, nil
}

func (r *fakeAnalyticsRepository) ListWeekdayClickAverages(
	ctx context.Context,
	linkID uuid.UUID,
	userID uuid.UUID,
) ([]entity.WeekdayClickAverage, error) {
	r.lastLinkID = linkID
	r.lastUserID = userID
	return r.weekdayClickAverages, nil
}

func (r *fakeAnalyticsRepository) UpdateAnalytics(
	ctx context.Context,
	updatedAnalytics entity.Analytics,
) (*entity.Analytics, error) {
	r.updateCalls++
	if updatedAnalytics.ID == uuid.Nil {
		updatedAnalytics.ID = uuid.New()
	}
	r.analytics[updatedAnalytics.LinkID] = updatedAnalytics
	return &updatedAnalytics, nil
}

func TestAnalyticsUsecaseUpdateAnalyticsPersistsValidAnalytics(t *testing.T) {
	repository := newFakeAnalyticsRepository()
	usecase := NewAnalyticsUseCase(repository)
	linkID := uuid.New()

	analytics, err := usecase.UpdateAnalytics(context.Background(), entity.Analytics{
		LinkID: linkID,
		Clicks: 3,
	})
	if err != nil {
		t.Fatalf("UpdateAnalytics returned error: %v", err)
	}

	if repository.updateCalls != 1 {
		t.Fatalf("expected UpdateAnalytics to be called once, got %d", repository.updateCalls)
	}
	if analytics.LinkID != linkID {
		t.Fatalf("expected link id %s, got %s", linkID, analytics.LinkID)
	}
	if repository.analytics[linkID].Clicks != 3 {
		t.Fatalf("expected 3 clicks, got %d", repository.analytics[linkID].Clicks)
	}
}

func TestAnalyticsUsecaseUpdateAnalyticsRejectsInvalidLinkID(t *testing.T) {
	repository := newFakeAnalyticsRepository()
	usecase := NewAnalyticsUseCase(repository)

	_, err := usecase.UpdateAnalytics(context.Background(), entity.Analytics{})
	if !errors.Is(err, derr.InvalidLinkId) {
		t.Fatalf("expected InvalidLinkId, got %v", err)
	}
	if repository.updateCalls != 0 {
		t.Fatalf("expected UpdateAnalytics not to be called, got %d calls", repository.updateCalls)
	}
}

func TestAnalyticsUsecaseGetByLinkIdValidatesIDs(t *testing.T) {
	repository := newFakeAnalyticsRepository()
	usecase := NewAnalyticsUseCase(repository)

	_, err := usecase.GetByLinkId(context.Background(), uuid.Nil, uuid.New())
	if !errors.Is(err, derr.InvalidLinkId) {
		t.Fatalf("expected InvalidLinkId, got %v", err)
	}

	_, err = usecase.GetByLinkId(context.Background(), uuid.New(), uuid.Nil)
	if !errors.Is(err, derr.UnauthorizedError) {
		t.Fatalf("expected UnauthorizedError, got %v", err)
	}
}

func TestAnalyticsUsecaseGetByLinkIdReturnsRepositoryAnalytics(t *testing.T) {
	repository := newFakeAnalyticsRepository()
	usecase := NewAnalyticsUseCase(repository)
	linkID := uuid.New()
	userID := uuid.New()
	expected := entity.Analytics{
		ID:     uuid.New(),
		LinkID: linkID,
		Clicks: 10,
	}
	repository.analytics[linkID] = expected

	analytics, err := usecase.GetByLinkId(context.Background(), linkID, userID)
	if err != nil {
		t.Fatalf("GetByLinkId returned error: %v", err)
	}
	if *analytics != expected {
		t.Fatalf("expected analytics %+v, got %+v", expected, *analytics)
	}
	if repository.lastLinkID != linkID || repository.lastUserID != userID {
		t.Fatalf("expected repository call with link %s and user %s", linkID, userID)
	}
}

func TestAnalyticsUsecaseListsClickAverages(t *testing.T) {
	repository := newFakeAnalyticsRepository()
	usecase := NewAnalyticsUseCase(repository)
	linkID := uuid.New()
	userID := uuid.New()
	repository.hourlyClickAverages = []entity.HourlyClickAverage{
		{Hour: 9, AverageClicks: 2.5},
	}
	repository.monthlyWeekClickAverages = []entity.MonthlyWeekClickAverage{
		{WeekOfMonth: 2, AverageClicks: 4},
	}
	repository.weekdayClickAverages = []entity.WeekdayClickAverage{
		{DayOfWeek: 1, AverageClicks: 7},
	}

	hourly, err := usecase.ListHourlyClickAverages(context.Background(), linkID, userID)
	if err != nil {
		t.Fatalf("ListHourlyClickAverages returned error: %v", err)
	}
	if !reflect.DeepEqual(hourly, repository.hourlyClickAverages) {
		t.Fatalf("expected hourly averages %+v, got %+v", repository.hourlyClickAverages, hourly)
	}

	monthly, err := usecase.ListMonthlyWeekClickAverages(context.Background(), linkID, userID)
	if err != nil {
		t.Fatalf("ListMonthlyWeekClickAverages returned error: %v", err)
	}
	if !reflect.DeepEqual(monthly, repository.monthlyWeekClickAverages) {
		t.Fatalf(
			"expected monthly week averages %+v, got %+v",
			repository.monthlyWeekClickAverages,
			monthly,
		)
	}

	weekday, err := usecase.ListWeekdayClickAverages(context.Background(), linkID, userID)
	if err != nil {
		t.Fatalf("ListWeekdayClickAverages returned error: %v", err)
	}
	if !reflect.DeepEqual(weekday, repository.weekdayClickAverages) {
		t.Fatalf("expected weekday averages %+v, got %+v", repository.weekdayClickAverages, weekday)
	}
}
