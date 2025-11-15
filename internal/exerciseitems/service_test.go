package exerciseitems

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"weight-tracker/internal/repository"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) GetById(ctx context.Context, arg repository.GetExerciseItemByIdParams) (ExerciseItem, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).(ExerciseItem), args.Error(1)
}

func (r *repositoryMock) GetByWorkoutId(ctx context.Context, arg repository.GetExerciseItemsByWorkoutIdParams) ([]ExerciseItem, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).([]ExerciseItem), args.Error(1)
}

func (r *repositoryMock) CreateAndReturnId(ctx context.Context, arg repository.CreateExerciseItemAndReturnIdParams) (string, error) {
	args := r.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

func (r *repositoryMock) UpdateType(ctx context.Context, arg repository.UpdateExerciseItemTypeParams) (int64, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repositoryMock) DeleteById(ctx context.Context, arg repository.DeleteExerciseItemByIdParams) (int64, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func TestNewService(t *testing.T) {
	mockRepo := new(repositoryMock)
	service := NewService(mockRepo)

	if service == nil {
		t.Error("Expected service to be created")
	}
}
