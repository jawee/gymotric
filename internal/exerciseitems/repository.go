package exerciseitems

import (
	"context"
	"fmt"
	"weight-tracker/internal/repository"
)

type ExerciseItemRepository interface {
	GetById(ctx context.Context, arg repository.GetExerciseItemByIdParams) (ExerciseItem, error)
	GetByWorkoutId(ctx context.Context, arg repository.GetExerciseItemsByWorkoutIdParams) ([]ExerciseItem, error)
	CreateAndReturnId(ctx context.Context, arg repository.CreateExerciseItemAndReturnIdParams) (string, error)
	UpdateType(ctx context.Context, arg repository.UpdateExerciseItemTypeParams) (int64, error)
	DeleteById(ctx context.Context, arg repository.DeleteExerciseItemByIdParams) (int64, error)
}

func NewExerciseItemRepository(repo repository.Querier) ExerciseItemRepository {
	return exerciseItemRepository{repo: repo}
}

type exerciseItemRepository struct {
	repo repository.Querier
}

func (e exerciseItemRepository) CreateAndReturnId(ctx context.Context, arg repository.CreateExerciseItemAndReturnIdParams) (string, error) {
	id, err := e.repo.CreateExerciseItemAndReturnId(ctx, arg)
	if err != nil {
		return "", fmt.Errorf("failed to create exercise item: %w", err)
	}
	return id, nil
}

func (e exerciseItemRepository) GetById(ctx context.Context, arg repository.GetExerciseItemByIdParams) (ExerciseItem, error) {
	item, err := e.repo.GetExerciseItemById(ctx, arg)
	if err != nil {
		return ExerciseItem{}, fmt.Errorf("failed to get exercise item: %w", err)
	}
	return newExerciseItem(item), nil
}

func (e exerciseItemRepository) GetByWorkoutId(ctx context.Context, arg repository.GetExerciseItemsByWorkoutIdParams) ([]ExerciseItem, error) {
	items, err := e.repo.GetExerciseItemsByWorkoutId(ctx, arg)
	if err != nil {
		return []ExerciseItem{}, fmt.Errorf("failed to get exercise items: %w", err)
	}

	result := []ExerciseItem{}
	for _, v := range items {
		result = append(result, newExerciseItem(v))
	}
	return result, nil
}

func (e exerciseItemRepository) UpdateType(ctx context.Context, arg repository.UpdateExerciseItemTypeParams) (int64, error) {
	return e.repo.UpdateExerciseItemType(ctx, arg)
}

func (e exerciseItemRepository) DeleteById(ctx context.Context, arg repository.DeleteExerciseItemByIdParams) (int64, error) {
	return e.repo.DeleteExerciseItemById(ctx, arg)
}

func newExerciseItem(v repository.ExerciseItem) ExerciseItem {
	return ExerciseItem{
		ID:        v.ID,
		Type:      v.Type,
		UserID:    v.UserID,
		WorkoutID: v.WorkoutID,
		CreatedOn: v.CreatedOn,
		UpdatedOn: v.UpdatedOn,
	}
}
