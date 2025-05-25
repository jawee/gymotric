package statistics

import (
	"context"
	"time"
	"weight-tracker/internal/repository"
)

type Statistics struct {
	Week          int
	PreviousWeek  int
	Month         int
	PreviousMonth int
	Year          int
	PreviousYear  int
}

type StatisticsRepository interface {
	GetStatistics(context context.Context, userId string) (Statistics, error)
}

type statisticsRepository struct {
	repo repository.Querier
}

func (s *statisticsRepository) GetStatistics(context context.Context, userId string) (Statistics, error) {
	week, err := s.repo.GetStatisticsSinceDate(context, repository.GetStatisticsSinceDateParams{
		UserID:    userId,
		StartDate: weekStartDate(time.Now()),
	})

	if err != nil {
		return Statistics{}, err
	}

	previousWeek, err := s.repo.GetStatisticsBetweenDates(context, repository.GetStatisticsBetweenDatesParams{
		UserID:    userId,
		StartDate: weekStartDate(time.Now().AddDate(0, 0, -7)),
		EndDate:   weekStartDate(time.Now()),
	})
	if err != nil {
		return Statistics{}, err
	}

	month, err := s.repo.GetStatisticsSinceDate(context, repository.GetStatisticsSinceDateParams{
		UserID:    userId,
		StartDate: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return Statistics{}, err
	}

	previousMonth, err := s.repo.GetStatisticsBetweenDates(context, repository.GetStatisticsBetweenDatesParams{
		UserID:    userId,
		StartDate: time.Date(time.Now().Year(), time.Now().Month()-1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return Statistics{}, err
	}

	year, err := s.repo.GetStatisticsSinceDate(context, repository.GetStatisticsSinceDateParams{
		UserID:    userId,
		StartDate: time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return Statistics{}, err
	}

	previousYear, err := s.repo.GetStatisticsBetweenDates(context, repository.GetStatisticsBetweenDatesParams{
		UserID:    userId,
		StartDate: time.Date(time.Now().Year()-1, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return Statistics{}, err
	}

	return Statistics{
		Week:          int(week),
		PreviousWeek:  int(previousWeek),
		Month:         int(month),
		PreviousMonth: int(previousMonth),
		Year:          int(year),
		PreviousYear:  int(previousYear),
	}, nil
}

func weekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}
