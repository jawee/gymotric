package workouts

import (
	"context"
	"errors"
	"testing"
	"time"
	"weight-tracker/internal/exercises"
	"weight-tracker/internal/exercisetypes"
	"weight-tracker/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testError = errors.New("Testerror")

type exerciseRepoMock struct {
	mock.Mock
}

func (r *exerciseRepoMock) GetAll(context context.Context, userId string) ([]exercises.Exercise, error) {
	args := r.Called(context, userId)
	return args.Get(0).([]exercises.Exercise), args.Error(1)
}
func (r *exerciseRepoMock) GetByWorkoutId(context context.Context, arg repository.GetExercisesByWorkoutIdParams) ([]exercises.Exercise, error) {
	args := r.Called(context, arg)
	return args.Get(0).([]exercises.Exercise), args.Error(1)
}
func (r *exerciseRepoMock) DeleteById(context context.Context, arg repository.DeleteExerciseByIdParams) error {
	args := r.Called(context, arg)
	return args.Error(0)
}
func (r *exerciseRepoMock) CreateAndReturnId(context context.Context, exercise repository.CreateExerciseAndReturnIdParams) (string, error) {
	args := r.Called(context, exercise)
	return args.String(0), args.Error(1)
}
func (r *exerciseRepoMock) GetExerciseTypeById(context context.Context, arg repository.GetExerciseTypeByIdParams) (*exercisetypes.ExerciseType, error) {
	args := r.Called(context, arg)
	return args.Get(0).(*exercisetypes.ExerciseType), args.Error(1)
}

type repoMock struct {
	mock.Mock
}

func (r *repoMock) UpdateById(ctx context.Context, arg repository.UpdateWorkoutByIdParams) error {
	args := r.Called(ctx, arg)
	return args.Error(0)
}

func (r *repoMock) CompleteById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) (int64, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repoMock) CreateAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error) {
	args := r.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

func (r *repoMock) GetAll(ctx context.Context, arg repository.GetAllWorkoutsParams) ([]Workout, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).([]Workout), args.Error(1)
}

func (r *repoMock) GetAllCount(ctx context.Context, userId string) (int64, error) {
	args := r.Called(ctx, userId)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repoMock) GetById(ctx context.Context, arg repository.GetWorkoutByIdParams) (Workout, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).(Workout), args.Error(1)
}

func (r *repoMock) DeleteById(ctx context.Context, arg repository.DeleteWorkoutByIdParams) error {
	args := r.Called(ctx, arg)
	return args.Error(0)
}

func TestGetAll(t *testing.T) {
	userId := "userid"
	expected := []Workout{
		{ID: "b", Name: "B", CreatedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339), CompletedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339), UpdatedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339)},
		{ID: "a", Name: "A", CreatedOn: time.Now().UTC().Format(time.RFC3339), CompletedOn: time.Now().UTC().Format(time.RFC3339), UpdatedOn: time.Now().UTC().Format(time.RFC3339)},
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAll", ctx, mock.MatchedBy(func(input repository.GetAllWorkoutsParams) bool {
		return input.UserID == userId && input.Offset == 0 && input.Limit == 10
	})).Return([]Workout{
		{ID: "a", Name: "A", CreatedOn: time.Now().UTC().Format(time.RFC3339), CompletedOn: time.Now().UTC().Format(time.RFC3339), UpdatedOn: time.Now().UTC().Format(time.RFC3339)},
		{ID: "b", Name: "B", CreatedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339), CompletedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339), UpdatedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339)},
	}, nil).Once()
	service := NewService(&repoMock, nil)

	result, err := service.GetAll(ctx, userId, 1, 10)

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expected[0].ID, result[0].ID)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}

func TestGetAllDefaultPage(t *testing.T) {
	userId := "userid"
	expected := []Workout{
		{ID: "b", Name: "B", CreatedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339), CompletedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339), UpdatedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339)},
		{ID: "a", Name: "A", CreatedOn: time.Now().UTC().Format(time.RFC3339), CompletedOn: time.Now().UTC().Format(time.RFC3339), UpdatedOn: time.Now().UTC().Format(time.RFC3339)},
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAll", ctx, mock.MatchedBy(func(input repository.GetAllWorkoutsParams) bool {
		return input.UserID == userId && input.Offset == 0 && input.Limit == 10
	})).Return([]Workout{
		{ID: "a", Name: "A", CreatedOn: time.Now().UTC().Format(time.RFC3339), CompletedOn: time.Now().UTC().Format(time.RFC3339), UpdatedOn: time.Now().UTC().Format(time.RFC3339)},
		{ID: "b", Name: "B", CreatedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339), CompletedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339), UpdatedOn: time.Now().Add(time.Minute * 2).UTC().Format(time.RFC3339)},
	}, nil).Once()
	service := NewService(&repoMock, nil)

	result, err := service.GetAll(ctx, userId, 0, 0)

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expected[0].ID, result[0].ID)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)

}

func TestGetAllRepoErr(t *testing.T) {
	userId := "userid"
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAll", ctx, mock.MatchedBy(func(input repository.GetAllWorkoutsParams) bool {
		return input.UserID == userId && input.Offset == 0 && input.Limit == 10
	})).Return([]Workout{}, testError).Once()

	service := NewService(&repoMock, nil)

	result, err := service.GetAll(ctx, userId, 1, 10)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	assert.Len(t, result, 0)
	repoMock.AssertExpectations(t)
}

func TestGetById(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	expected := Workout{
		ID:          workoutId,
		Name:        "A",
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn:   time.Now().UTC().Format(time.RFC3339),
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetById", ctx, mock.MatchedBy(func(input repository.GetWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId
	})).Return(Workout{
		ID:          workoutId,
		Name:        "A",
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn:   time.Now().UTC().Format(time.RFC3339),
	}, nil).Once()

	service := NewService(&repoMock, nil)

	result, err := service.GetById(ctx, workoutId, userId)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}

func TestCreateAndReturnId(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	ctx := context.Background()

	request := createWorkoutRequest{
		Name: "A",
	}

	repoMock := repoMock{}
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateWorkoutAndReturnIdParams) bool {
		return input.Name == "A" && input.ID != "" && input.CreatedOn != "" && input.UpdatedOn != "" && input.UserID == userId
	})).Return(workoutId, nil).Once()

	service := NewService(&repoMock, nil)

	result, err := service.CreateAndReturnId(ctx, request, userId)

	assert.Nil(t, err)
	assert.Equal(t, workoutId, result)
	repoMock.AssertExpectations(t)
}

func TestCompleteById(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("CompleteById", ctx, mock.MatchedBy(func(input repository.CompleteWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId && input.CompletedOn != ""
	})).Return(int64(1), nil).Once()

	service := NewService(&repoMock, nil)

	err := service.CompleteById(ctx, workoutId, userId)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestDeleteById(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("DeleteById", ctx, mock.MatchedBy(func(input repository.DeleteWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId
	})).Return(nil).Once()

	service := NewService(&repoMock, nil)
	err := service.DeleteById(ctx, workoutId, userId)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestUpdateById(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	ctx := context.Background()

	request := updateWorkoutRequest{
		Note: "The note",
	}

	repoMock := repoMock{}
	repoMock.On("GetById", ctx, mock.MatchedBy(func(input repository.GetWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId
	})).Return(Workout{
		ID:          workoutId,
		Name:        "A",
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn:   time.Now().UTC().Format(time.RFC3339),
	}, nil).Once()
	repoMock.On("UpdateById", ctx, mock.MatchedBy(func(input repository.UpdateWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId && input.Note == request.Note
	})).Return(nil).Once()

	service := NewService(&repoMock, nil)

	err := service.UpdateWorkoutById(ctx, workoutId, request, userId)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestUpdateByIdNotFound(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	ctx := context.Background()

	request := updateWorkoutRequest{
		Note: "The note",
	}

	repoMock := repoMock{}
	repoMock.On("GetById", ctx, mock.MatchedBy(func(input repository.GetWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId
	})).Return(Workout{}, testError).Once()

	service := NewService(&repoMock, nil)

	err := service.UpdateWorkoutById(ctx, workoutId, request, userId)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestUpdateByIdUpdateErr(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	ctx := context.Background()

	request := updateWorkoutRequest{
		Note: "The note",
	}

	repoMock := repoMock{}
	repoMock.On("GetById", ctx, mock.MatchedBy(func(input repository.GetWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId
	})).Return(Workout{
		ID:          workoutId,
		Name:        "A",
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn:   time.Now().UTC().Format(time.RFC3339),
	}, nil).Once()
	repoMock.On("UpdateById", ctx, mock.Anything).Return(testError).Once()

	service := NewService(&repoMock, nil)

	err := service.UpdateWorkoutById(ctx, workoutId, request, userId)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestCloneByIdAndReturnId(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	newWorkoutId := "newWorkoutId"
	ctx := context.Background()

	request := createWorkoutRequest{
		Name: "A",
	}

	repoMock := repoMock{}
	repoMock.On("GetById", ctx, mock.MatchedBy(func(input repository.GetWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId
	})).Return(Workout{
		ID:          workoutId,
		Name:        "A",
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn:   time.Now().UTC().Format(time.RFC3339),
	}, nil).Once()
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateWorkoutAndReturnIdParams) bool {
		return input.Name == request.Name && input.UserID == userId
	})).Return(newWorkoutId, nil).Once()

	exerciseRepoMock := exerciseRepoMock{}
	exerciseRepoMock.On("GetByWorkoutId", ctx, mock.MatchedBy(func(input repository.GetExercisesByWorkoutIdParams) bool {
		return input.UserID == userId && input.WorkoutID == workoutId
	})).Return([]exercises.Exercise{
		{
			ID:             "exerciseId",
			Name:           "Exercise A",
			ExerciseTypeID: "exerciseTypeId",
			WorkoutID:      workoutId,
		},
	}, nil).Once()

	exerciseRepoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateExerciseAndReturnIdParams) bool {
		return input.WorkoutID == newWorkoutId && input.UserID == userId && input.Name == "Exercise A" && input.ExerciseTypeID == "exerciseTypeId"
	})).Return("newExerciseId", nil).Once()

	service := NewService(&repoMock, &exerciseRepoMock)

	result, err := service.CloneByIdAndReturnId(ctx, workoutId, userId)

	assert.Nil(t, err)
	assert.Equal(t, newWorkoutId, result)
	repoMock.AssertExpectations(t)
	exerciseRepoMock.AssertExpectations(t)
}

func TestCloneByIdAndReturnIdSourceNotFound(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetById", ctx, mock.MatchedBy(func(input repository.GetWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId
	})).Return(Workout{}, testError).Once()

	service := NewService(&repoMock, nil)

	result, err := service.CloneByIdAndReturnId(ctx, workoutId, userId)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	assert.Equal(t, "", result)
	repoMock.AssertExpectations(t)
}

func TestCloneByIdAndReturnIdCreateAndReturnIdErr(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetById", ctx, mock.MatchedBy(func(input repository.GetWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId
	})).Return(Workout{
		ID:          workoutId,
		Name:        "A",
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn:   time.Now().UTC().Format(time.RFC3339),
	}, nil).Once()

	repoMock.On("CreateAndReturnId", ctx, mock.Anything).Return("", testError).Once()

	service := NewService(&repoMock, nil)

	result, err := service.CloneByIdAndReturnId(ctx, workoutId, userId)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	assert.Equal(t, "", result)
	repoMock.AssertExpectations(t)
}

func TestCloneByIdAndReturnIdCreateExerciseErr(t *testing.T) {
	userId := "userid"
	workoutId := "workoutId"
	newWorkoutId := "newWorkoutId"
	ctx := context.Background()

	request := createWorkoutRequest{
		Name: "A",
	}

	repoMock := repoMock{}
	repoMock.On("GetById", ctx, mock.MatchedBy(func(input repository.GetWorkoutByIdParams) bool {
		return input.ID == workoutId && input.UserID == userId
	})).Return(Workout{
		ID:          workoutId,
		Name:        "A",
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn:   time.Now().UTC().Format(time.RFC3339),
	}, nil).Once()
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateWorkoutAndReturnIdParams) bool {
		return input.Name == request.Name && input.UserID == userId
	})).Return(newWorkoutId, nil).Once()

	exerciseRepoMock := exerciseRepoMock{}
	exerciseRepoMock.On("GetByWorkoutId", ctx, mock.MatchedBy(func(input repository.GetExercisesByWorkoutIdParams) bool {
		return input.UserID == userId && input.WorkoutID == workoutId
	})).Return([]exercises.Exercise{
		{
			ID:             "exerciseId",
			Name:           "Exercise A",
			ExerciseTypeID: "exerciseTypeId",
			WorkoutID:      workoutId,
		},
	}, nil).Once()

	exerciseRepoMock.On("CreateAndReturnId", ctx, mock.Anything).Return("", testError).Once()

	service := NewService(&repoMock, &exerciseRepoMock)

	result, err := service.CloneByIdAndReturnId(ctx, workoutId, userId)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	assert.Equal(t, "", result)
	repoMock.AssertExpectations(t)
	exerciseRepoMock.AssertExpectations(t)
}

func TestGetAllCount(t *testing.T) {
	userId := "userid"
	ctx := context.Background()
	expectedCount := 5

	repoMock := repoMock{}
	repoMock.On("GetAllCount", ctx, userId).Return(int64(expectedCount), nil).Once()

	service := NewService(&repoMock, nil)

	count, err := service.GetAllCount(ctx, userId)

	assert.Nil(t, err)
	assert.Equal(t, expectedCount, count)
	repoMock.AssertExpectations(t)
}

func TestGetAllCountError(t *testing.T) {
	userId := "userid"
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAllCount", ctx, userId).Return(int64(0), testError).Once()

	service := NewService(&repoMock, nil)

	count, err := service.GetAllCount(ctx, userId)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	assert.Equal(t, 0, count)
	repoMock.AssertExpectations(t)
}

