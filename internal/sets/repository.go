package sets

import (
	"context"
	"fmt"
	"weight-tracker/internal/repository"
)

type Set struct {
	ID          string  `json:"id"`
	Repetitions int64   `json:"repetitions"`
	Weight      float64 `json:"weight"`
	ExerciseID  string  `json:"exercise_id"`
}

type SetsRepository interface {
	GetAll(ctx context.Context, userId string) ([]Set, error)
	GetById(ctx context.Context, arg repository.GetSetByIdParams) (Set, error)
	GetByExerciseId(ctx context.Context, arg repository.GetSetsByExerciseIdParams) ([]Set, error)
	CreateAndReturnId(ctx context.Context, arg repository.CreateSetAndReturnIdParams) (string, error)
	DeleteById(ctx context.Context, arg repository.DeleteSetByIdParams) (int64, error)
}

type setsRepository struct {
	repo repository.Querier
}

func (s *setsRepository) CreateAndReturnId(ctx context.Context, arg repository.CreateSetAndReturnIdParams) (string, error) {
	return s.repo.CreateSetAndReturnId(ctx, arg)
}

func (s *setsRepository) DeleteById(ctx context.Context, arg repository.DeleteSetByIdParams) (int64, error) {
	return s.repo.DeleteSetById(ctx, arg)
}

func (s *setsRepository) GetAll(ctx context.Context, userId string) ([]Set, error) {
	sets, err := s.repo.GetAllSets(ctx, userId)
	if err != nil {
		return []Set{}, fmt.Errorf("failed to get all sets: %w", err)
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

func (s *setsRepository) GetById(ctx context.Context, arg repository.GetSetByIdParams) (Set, error) {
	set, err := s.repo.GetSetById(ctx, arg)
	if err != nil {
		return Set{}, fmt.Errorf("failed to get set by id: %w", err)
	}

	result := newSet(set)

	return result, nil
}

func (s *setsRepository) GetByExerciseId(ctx context.Context, arg repository.GetSetsByExerciseIdParams) ([]Set, error) {
	sets, err := s.repo.GetSetsByExerciseId(ctx, arg)
	if err != nil {
		return []Set{}, fmt.Errorf("failed to get sets by exercise id: %w", err)
	}

	result := []Set{}
	for _, v := range sets {
		result = append(result, newSet(v))
	}

	return result, nil
}
