package statistics

import (
	"context"
	"time"
	"weight-tracker/internal/repository"
)

type Statistics struct {
	Week  int
	Month int
	Year  int
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

	month, err := s.repo.GetStatisticsSinceDate(context, repository.GetStatisticsSinceDateParams{
		UserID:    userId,
		StartDate: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
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

	return Statistics{
		Week:  int(week),
		Month: int(month),
		Year:  int(year),
	}, nil
}

func weekStartDate(date time.Time) time.Time {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	result := date.Add(time.Duration(offset*24) * time.Hour)
	return result
}
