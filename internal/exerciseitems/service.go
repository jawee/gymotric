package exerciseitems

import (
	"context"
	"fmt"
	"time"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	GetById(ctx context.Context, arg repository.GetExerciseItemByIdParams) (ExerciseItem, error)
	GetByWorkoutId(ctx context.Context, arg repository.GetExerciseItemsByWorkoutIdParams) ([]ExerciseItem, error)
	CreateAndReturnId(ctx context.Context, itemType string, workoutId string, userId string) (string, error)
	UpdateType(ctx context.Context, arg repository.UpdateExerciseItemTypeParams) (int64, error)
	DeleteById(ctx context.Context, arg repository.DeleteExerciseItemByIdParams) (int64, error)
}

type exerciseItemService struct {
	repo ExerciseItemRepository
}

func NewService(repo ExerciseItemRepository) Service {
	return &exerciseItemService{repo}
}

func (s *exerciseItemService) CreateAndReturnId(ctx context.Context, itemType string, workoutId string, userId string) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	params := repository.CreateExerciseItemAndReturnIdParams{
		ID:        uuid.String(),
		Type:      itemType,
		UserID:    userId,
		WorkoutID: workoutId,
		CreatedOn: now,
		UpdatedOn: now,
	}

	id, err := s.repo.CreateAndReturnId(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to create exercise item: %w", err)
	}
	return id, nil
}

func (s *exerciseItemService) GetById(ctx context.Context, arg repository.GetExerciseItemByIdParams) (ExerciseItem, error) {
	return s.repo.GetById(ctx, arg)
}

func (s *exerciseItemService) GetByWorkoutId(ctx context.Context, arg repository.GetExerciseItemsByWorkoutIdParams) ([]ExerciseItem, error) {
	return s.repo.GetByWorkoutId(ctx, arg)
}

func (s *exerciseItemService) UpdateType(ctx context.Context, arg repository.UpdateExerciseItemTypeParams) (int64, error) {
	return s.repo.UpdateType(ctx, arg)
}

func (s *exerciseItemService) DeleteById(ctx context.Context, arg repository.DeleteExerciseItemByIdParams) (int64, error) {
	return s.repo.DeleteById(ctx, arg)
}
