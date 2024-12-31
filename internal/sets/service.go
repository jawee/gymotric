package sets

import (
	"context"
	"weight-tracker/internal/repository"
)

type Service interface {
	GetByExerciseId(context context.Context, exerciseId string) ([]Set, error)
	DeleteById(context context.Context, setId string) error
	CreateAndReturnId(context context.Context, t createSetRequest, exerciseId string) (string, error)
}

func (s *setsService) GetByExerciseId(context context.Context, exerciseId string) ([]Set, error) {
	return s.repo.GetsByExerciseId(context, exerciseId)
}

func (s *setsService) DeleteById(context context.Context, setId string) error {
	_, err := s.repo.DeleteById(context, setId)
	return err
}

func (s *setsService) CreateAndReturnId(context context.Context, t createSetRequest, exerciseId string) (string, error) {
	set := repository.CreateSetAndReturnIdParams{
		ID:          generateUuid(),
		Repetitions: int64(t.Repetitions),
		Weight:      t.Weight,
		ExerciseID:  exerciseId,
	}
	id, err := s.repo.CreateAndReturnId(context, set)
	return id, err
}

type setsService struct {
	repo SetsRepository
}

func NewService(repo SetsRepository) Service {
	return &setsService{repo}
}
