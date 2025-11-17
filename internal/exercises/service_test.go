package exercises

import (
	"context"
	"testing"
	"weight-tracker/internal/exercisetypes"
	"weight-tracker/internal/repository"

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

func (r *repoMock) GetByExerciseItemId(context context.Context, exerciseItemId string, userId string) ([]Exercise, error) {
	args := r.Called(context, exerciseItemId, userId)
	return args.Get(0).([]Exercise), args.Error(1)
}

func (r *repoMock) GetExerciseTypeById(context context.Context, arg repository.GetExerciseTypeByIdParams) (*exercisetypes.ExerciseType, error) {
	args := r.Called(context, arg)
	return args.Get(0).(*exercisetypes.ExerciseType), args.Error(1)
}

func TestGetAll(t *testing.T) {
	userId := "userid"
	expected := []Exercise{
		{ID: "a", Name: "a", WorkoutID: "", ExerciseTypeID: ""},
		{ID: "b", Name: "b", WorkoutID: "", ExerciseTypeID: ""},
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAll", ctx, userId).Return([]Exercise{
		{ID: "a", Name: "a", WorkoutID: "", ExerciseTypeID: ""},
		{ID: "b", Name: "b", WorkoutID: "", ExerciseTypeID: ""},
	}, nil).Once()

	service := NewService(&repoMock)

	result, err := service.GetAll(ctx, userId)

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}

func TestCreateAndReturnId(t *testing.T) {
	userId := "userid"
	ctx := context.Background()

	exerciseTypeId := "exerciseTypeId"
	exerciseTypeName := "example-exercise"
	workoutId := "workoutId"
	exerciseId := "exerciseId"

	repoMock := repoMock{}
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateExerciseAndReturnIdParams) bool {
		return input.UserID == userId && input.Name == exerciseTypeName && input.WorkoutID == workoutId && input.CreatedOn != "" && input.UpdatedOn != ""
	})).Return(exerciseId, nil).Once()

	repoMock.On("GetExerciseTypeById", ctx, mock.MatchedBy(func(input repository.GetExerciseTypeByIdParams) bool {
		return input.ID == exerciseTypeId && input.UserID == userId
	})).Return(&exercisetypes.ExerciseType{ID: exerciseTypeId, Name: exerciseTypeName}, nil).Once()

	service := NewService(&repoMock)
	id, err := service.CreateAndReturnId(context.Background(), createExerciseRequest{
		ExerciseTypeID: exerciseTypeId,
	}, workoutId, userId)

	assert.Nil(t, err)
	assert.Equal(t, exerciseId, id)
	repoMock.AssertExpectations(t)
}

func TestDeleteById(t *testing.T) {
	userId := "userid"
	exerciseId := "exerciseId"
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("DeleteById", ctx, mock.MatchedBy(func(input repository.DeleteExerciseByIdParams) bool {
		return input.UserID == userId && input.ID == exerciseId
	})).Return(nil).Once()

	service := NewService(&repoMock)
	err := service.DeleteById(ctx, exerciseId, userId)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestGetByWorkoutId(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	expected := []Exercise{
		{ID: "a", Name: "a", WorkoutID: workoutId, ExerciseTypeID: ""},
		{ID: "b", Name: "b", WorkoutID: workoutId, ExerciseTypeID: ""},
	}
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetByWorkoutId", ctx, mock.MatchedBy(func(input repository.GetExercisesByWorkoutIdParams) bool {
		return input.WorkoutID == workoutId && input.UserID == userId
	})).Return([]Exercise{
		{ID: "a", Name: "a", WorkoutID: workoutId, ExerciseTypeID: ""},
		{ID: "b", Name: "b", WorkoutID: workoutId, ExerciseTypeID: ""},
	}, nil).Once()

	service := NewService(&repoMock)
	result, err := service.GetByWorkoutId(ctx, workoutId, userId)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}
