package sets

import (
	"context"
	"weight-tracker/internal/repository"
)

type Set struct {
	ID          string  `json:"id"`
	Repetitions int64   `json:"repetitions"`
	Weight      float64 `json:"weight"`
	ExerciseID  string  `json:"exercise_id"`
}

type SetsRepository interface {
	GetAll(ctx context.Context) ([]Set, error)
	GetById(ctx context.Context, id string) (Set, error)
	GetByExerciseId(ctx context.Context, exerciseID string) ([]Set, error)
	CreateAndReturnId(ctx context.Context, arg repository.CreateSetAndReturnIdParams) (string, error)
	DeleteById(ctx context.Context, id string) (int64, error)
}

type setsRepository struct {
	repo repository.Querier
}

func (s *setsRepository) CreateAndReturnId(ctx context.Context, arg repository.CreateSetAndReturnIdParams) (string, error) {
	return s.repo.CreateSetAndReturnId(ctx, arg)
}

func (s *setsRepository) DeleteById(ctx context.Context, id string) (int64, error) {
	return s.repo.DeleteSetById(ctx, id)
}

func (s *setsRepository) GetAll(ctx context.Context) ([]Set, error) {
	sets, err := s.repo.GetAllSets(ctx)
	if err != nil {
		return []Set{}, err
	}

	result := []Set{}
	for _, v := range sets {
		result = append(result, newSet(v))
	}

	return result, nil
}

func newSet(v repository.Set) Set {
	set := Set{
		ID:          v.ID,
		Repetitions: v.Repetitions,
		Weight:      v.Weight,
		ExerciseID:  v.ExerciseID,
	}

	return set
}

func (s *setsRepository) GetById(ctx context.Context, id string) (Set, error) {
	set, err := s.repo.GetSetById(ctx, id)
	if err != nil {
		return Set{}, err
	}

	result := newSet(set)

	return result, nil
}

func (s *setsRepository) GetByExerciseId(ctx context.Context, exerciseID string) ([]Set, error) {
	sets, err := s.repo.GetSetsByExerciseId(ctx, exerciseID)
	if err != nil {
		return []Set{}, err
	}

	result := []Set{}
	for _, v := range sets {
		result = append(result, newSet(v))
	}

	return result, nil
}
