package workouts

import (
	"context"
	"testing"
	"time"
	"weight-tracker/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func (r *repoMock) GetAll(ctx context.Context, userId string) ([]Workout, error) {
	args := r.Called(ctx, userId)
	return args.Get(0).([]Workout), args.Error(1)
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
		{ID: "a", Name: "A", CreatedOn: time.Now().UTC().Format(time.RFC3339), CompletedOn: time.Now().UTC().Format(time.RFC3339), UpdatedOn: time.Now().UTC().Format(time.RFC3339)},
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAll", ctx, userId).Return([]Workout{
		{ID: "a", Name: "A", CreatedOn: time.Now().UTC().Format(time.RFC3339), CompletedOn: time.Now().UTC().Format(time.RFC3339), UpdatedOn: time.Now().UTC().Format(time.RFC3339)},
	}, nil).Once()

	service := NewService(&repoMock, nil) 

	result, err := service.GetAll(ctx, userId)

	assert.Nil(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expected, result)
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
