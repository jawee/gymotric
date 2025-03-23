package statistics

import "context"

func NewService(repo StatisticsRepository) Service {
	return &statisticsService{repo}
}

type Service interface {
	GetStatistics(context context.Context, userId string) (Statistics, error)
} 


type statisticsService struct {
	repo StatisticsRepository
}

func (s *statisticsService) GetStatistics(context context.Context, userId string) (Statistics, error) {
	return s.repo.GetStatistics(context, userId)
}
