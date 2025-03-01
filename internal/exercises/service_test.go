package exercises

import (
	"context"
	"testing"
	"weight-tracker/internal/exercisetypes"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	mock.Mock
}

func (r *repoMock) CreateAndReturnId(context context.Context, exercise repository.CreateExerciseAndReturnIdParams) (string, error) {
	args := r.Called(context, exercise)
	return args.String(0), args.Error(1)
}

func (r *repoMock) DeleteById(context context.Context, arg repository.DeleteExerciseByIdParams) error {
	args := r.Called(context, arg)
	return args.Error(0)
}

func (r *repoMock) GetAll(context context.Context, userId string) ([]Exercise, error) {
	args := r.Called(context, userId)
	return args.Get(0).([]Exercise), args.Error(1)
}

func (r *repoMock) GetByWorkoutId(context context.Context, arg repository.GetExercisesByWorkoutIdParams) ([]Exercise, error) {
	args := r.Called(context, arg)
	return args.Get(0).([]Exercise), args.Error(1)
}

func (r *repoMock) GetExerciseTypeById(context context.Context, arg repository.GetExerciseTypeByIdParams) (*exercisetypes.ExerciseType, error) {
	args := r.Called(context, arg)
	return args.Get(0).(*exercisetypes.ExerciseType), args.Error(1)
}

func TestGetAll(t *testing.T) {
	userId, _ := uuid.NewV7()
	expected := []Exercise{
		{ID: "a", Name: "a", WorkoutID: "", ExerciseTypeID: ""},
		{ID: "b", Name: "b", WorkoutID: "", ExerciseTypeID: ""},
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAll", ctx, userId.String()).Return([]Exercise{
		{ID: "a", Name: "a", WorkoutID: "", ExerciseTypeID: ""},
		{ID: "b", Name: "b", WorkoutID: "", ExerciseTypeID: ""},
	}, nil).Once()

	service := NewService(&repoMock)

	result, err := service.GetAll(ctx, userId.String())

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}

func TestCreateAndReturnId(t *testing.T) {
	userId, _ := uuid.NewV7()
	ctx := context.Background()

	exerciseTypeId, _ := uuid.NewV7()
	exerciseTypeName := "example-exercise"
	workoutId, _ := uuid.NewV7()
	exerciseId, _ := uuid.NewV7()

	repoMock := repoMock{}
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateExerciseAndReturnIdParams) bool {
		return input.UserID == userId.String() && input.Name == exerciseTypeName && input.WorkoutID == workoutId.String() && input.CreatedOn != "" && input.UpdatedOn != ""
	})).Return(exerciseId.String(), nil).Once()

	repoMock.On("GetExerciseTypeById", ctx, mock.MatchedBy(func(input repository.GetExerciseTypeByIdParams) bool {
		return input.ID == exerciseTypeId.String() && input.UserID == userId.String()
	})).Return(&exercisetypes.ExerciseType{ID: exerciseTypeId.String(), Name: exerciseTypeName}, nil).Once()

	service := NewService(&repoMock)
	id, err := service.CreateAndReturnId(context.Background(), createExerciseRequest{
		ExerciseTypeID: exerciseTypeId.String(),
	}, workoutId.String(), userId.String())

	assert.Nil(t, err)
	assert.Equal(t, exerciseId.String(), id)
	repoMock.AssertExpectations(t)
}

func TestDeleteById(t *testing.T) {
	exerciseId, _ := uuid.NewV7()
	userId, _ := uuid.NewV7()
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("DeleteById", ctx, mock.MatchedBy(func(input repository.DeleteExerciseByIdParams) bool {
		return input.UserID == userId.String() && input.ID == exerciseId.String()
	})).Return(nil).Once()

	service := NewService(&repoMock)
	err := service.DeleteById(ctx, exerciseId.String(), userId.String())

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestGetByWorkoutId(t *testing.T) {
userId, _ := uuid.NewV7()
	workoutId, _ := uuid.NewV7()
	expected := []Exercise{
		{ID: "a", Name: "a", WorkoutID: workoutId.String(), ExerciseTypeID: ""},
		{ID: "b", Name: "b", WorkoutID: workoutId.String(), ExerciseTypeID: ""},
	}
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetByWorkoutId", ctx, mock.MatchedBy(func(input repository.GetExercisesByWorkoutIdParams) bool {
		return input.WorkoutID == workoutId.String() && input.UserID == userId.String()
	})).Return([]Exercise{
		{ID: "a", Name: "a", WorkoutID: workoutId.String(), ExerciseTypeID: ""},
		{ID: "b", Name: "b", WorkoutID: workoutId.String(), ExerciseTypeID: ""},
	}, nil).Once()

	service := NewService(&repoMock)
	result, err := service.GetByWorkoutId(ctx, workoutId.String(), userId.String())

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}
