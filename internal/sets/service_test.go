package sets

//
// import (
// 	"context"
// 	"testing"
// 	"weight-tracker/internal/exercisetypes"
// 	"weight-tracker/internal/repository"
//
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )
//
// type repoMock struct {
// 	mock.Mock
// }
//
// func (r *repoMock) CreateAndReturnId(context context.Context, exercise repository.CreateExerciseAndReturnIdParams, workoutId string) (string, error) {
// 	args := r.Called(context, exercise, workoutId)
// 	return args.String(0), args.Error(1)
// }
//
// func (r *repoMock) DeleteById(context context.Context, id string) error {
// 	args := r.Called(context, id)
// 	return args.Error(0)
// }
//
// func (r *repoMock) GetAll(context context.Context) ([]Exercise, error) {
// 	args := r.Called(context)
// 	return args.Get(0).([]Exercise), args.Error(1)
// }
//
// func (r *repoMock) GetByWorkoutId(context context.Context, workoutId string) ([]Exercise, error) {
// 	args := r.Called(context, workoutId)
// 	return args.Get(0).([]Exercise), args.Error(1)
// }
//
// func (r *repoMock) GetExerciseTypeById(context context.Context, exerciseTypeId string) (*exercisetypes.ExerciseType, error) {
// 	args := r.Called(context, exerciseTypeId)
// 	return args.Get(0).(*exercisetypes.ExerciseType), args.Error(1)
// }
//
// func TestGetAll(t *testing.T) {
// 	expected := []Exercise{
// 		{ID: "a", Name: "a", WorkoutID: "", ExerciseTypeID: ""},
// 		{ID: "b", Name: "b", WorkoutID: "", ExerciseTypeID: ""},
// 	}
//
// 	ctx := context.Background()
//
// 	repoMock := repoMock{}
// 	repoMock.On("GetAll", ctx).Return([]Exercise{
// 		{ID: "a", Name: "a", WorkoutID: "", ExerciseTypeID: ""},
// 		{ID: "b", Name: "b", WorkoutID: "", ExerciseTypeID: ""},
// 	}, nil).Once()
//
// 	service := NewService(&repoMock)
//
// 	result, err := service.GetAll(ctx)
//
// 	assert.Nil(t, err)
// 	assert.Len(t, result, 2)
// 	assert.Equal(t, expected, result)
// 	repoMock.AssertExpectations(t)
// }
//
// func TestCreateAndReturnId(t *testing.T) {
// 	ctx := context.Background()
//
// 	exerciseTypeId := generateUuid()
// 	exerciseTypeName := "example-exercise"
// 	workoutId := generateUuid()
// 	exerciseId := generateUuid()
//
//
// 	repoMock := repoMock{}
// 	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateExerciseAndReturnIdParams) bool {
// 		return input.Name == exerciseTypeName && input.WorkoutID == workoutId
// 	}), workoutId).Return(exerciseId, nil).Once()
//
// 	repoMock.On("GetExerciseTypeById", ctx, exerciseTypeId).Return(&exercisetypes.ExerciseType{ID: exerciseTypeId, Name: exerciseTypeName}, nil).Once()
//
// 	service := NewService(&repoMock)
// 	id, err := service.CreateAndReturnId(context.Background(), createExerciseRequest{
// 		ExerciseTypeID: exerciseTypeId,
// 	}, workoutId)
//
// 	assert.Nil(t, err)
// 	assert.Equal(t, exerciseId, id)
// 	repoMock.AssertExpectations(t)
// }
//
// func TestDeleteById(t *testing.T) {
// 	ctx := context.Background()
//
// 	repoMock := repoMock{}
// 	repoMock.On("DeleteById", ctx, "a").Return(nil).Once()
//
// 	service := NewService(&repoMock)
// 	err := service.DeleteById(ctx, "a")
//
// 	assert.Nil(t, err)
// 	repoMock.AssertExpectations(t)
// }
//
// func TestGetByWorkoutId(t *testing.T) {
// 	workoutId := generateUuid()
// 	expected := []Exercise{
// 		{ID: "a", Name: "a", WorkoutID: workoutId, ExerciseTypeID: ""},
// 		{ID: "b", Name: "b", WorkoutID: workoutId, ExerciseTypeID: ""},
// 	}
// 	ctx := context.Background()
//
// 	repoMock := repoMock{}
// 	repoMock.On("GetByWorkoutId", ctx, workoutId).Return([]Exercise{
// 		{ID: "a", Name: "a", WorkoutID: workoutId, ExerciseTypeID: ""},
// 		{ID: "b", Name: "b", WorkoutID: workoutId, ExerciseTypeID: ""},
// 	}, nil).Once()
//
// 	service := NewService(&repoMock)
// 	result, err := service.GetByWorkoutId(ctx, workoutId)
//
// 	assert.Nil(t, err)
// 	assert.Equal(t, expected, result)
// 	repoMock.AssertExpectations(t)
// }
