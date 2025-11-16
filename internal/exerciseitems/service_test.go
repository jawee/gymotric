package exerciseitems

import (
	"context"
	"testing"

	"weight-tracker/internal/exercises"
	"weight-tracker/internal/exercisetypes"
	"weight-tracker/internal/repository"

	"github.com/stretchr/testify/mock"
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

type exerciseRepositoryMock struct {
	mock.Mock
}

func (r *exerciseRepositoryMock) GetAll(ctx context.Context, userId string) ([]exercises.Exercise, error) {
	args := r.Called(ctx, userId)
	return args.Get(0).([]exercises.Exercise), args.Error(1)
}

func (r *exerciseRepositoryMock) GetByWorkoutId(ctx context.Context, arg repository.GetExercisesByWorkoutIdParams) ([]exercises.Exercise, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).([]exercises.Exercise), args.Error(1)
}

func (r *exerciseRepositoryMock) GetByExerciseItemId(ctx context.Context, exerciseItemId string, userId string) ([]exercises.Exercise, error) {
	args := r.Called(ctx, exerciseItemId, userId)
	return args.Get(0).([]exercises.Exercise), args.Error(1)
}

func (r *exerciseRepositoryMock) DeleteById(ctx context.Context, arg repository.DeleteExerciseByIdParams) error {
	args := r.Called(ctx, arg)
	return args.Error(0)
}

func (r *exerciseRepositoryMock) CreateAndReturnId(ctx context.Context, exercise repository.CreateExerciseAndReturnIdParams) (string, error) {
	args := r.Called(ctx, exercise)
	return args.String(0), args.Error(1)
}

func (r *exerciseRepositoryMock) GetExerciseTypeById(ctx context.Context, arg repository.GetExerciseTypeByIdParams) (*exercisetypes.ExerciseType, error) {
	args := r.Called(ctx, arg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*exercisetypes.ExerciseType), args.Error(1)
}

func TestNewService(t *testing.T) {
	mockRepo := new(repositoryMock)
	mockExerciseRepo := new(exerciseRepositoryMock)
	service := NewService(mockRepo, mockExerciseRepo)

	if service == nil {
		t.Error("Expected service to be created")
	}
}
