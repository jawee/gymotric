package sets

import (
	"context"
	"time"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	GetByExerciseId(context context.Context, exerciseId string, userId string) ([]Set, error)
	DeleteById(context context.Context, setId string, userId string) error
	CreateAndReturnId(context context.Context, t createSetRequest, exerciseId string, userId string) (string, error)
}

func (s *setsService) GetByExerciseId(context context.Context, exerciseId string, userId string) ([]Set, error) {
	arg := repository.GetSetsByExerciseIdParams{
		ExerciseID: exerciseId,
		UserID: userId,
	}
	return s.repo.GetByExerciseId(context, arg)
}

func (s *setsService) DeleteById(context context.Context, setId string, userId string) error {
	arg := repository.DeleteSetByIdParams{
		ID:     setId,
		UserID: userId,
	}
	_, err := s.repo.DeleteById(context, arg)
	return err
}

func (s *setsService) CreateAndReturnId(context context.Context, t createSetRequest, exerciseId string, userId string) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	set := repository.CreateSetAndReturnIdParams{
		ID:          uuid.String(),
		Repetitions: int64(t.Repetitions),
		Weight:      t.Weight,
		ExerciseID:  exerciseId,
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
		UserID: userId,
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
