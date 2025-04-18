package workouts

import (
	"context"
	"sort"
	"time"
	"weight-tracker/internal/exercises"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	GetAll(context context.Context, userId string, page int, pageSize int) ([]Workout, error)
	GetAllCount(context context.Context, userId string) (int, error)
	GetById(context context.Context, id string, userId string) (Workout, error)
	CreateAndReturnId(context context.Context, t createWorkoutRequest, userId string) (string, error)
	CompleteById(context context.Context, workoutId string, userId string) error
	DeleteById(context context.Context, workoutId string, userId string) error
	CloneByIdAndReturnId(context context.Context, workoutId string, userId string) (string, error)
	UpdateWorkoutById(context context.Context, workoutId string, t updateWorkoutRequest, userId string) error
}

type workoutsService struct {
	repo         WorkoutsRepository
	exerciseRepo exercises.ExerciseRepository
}

func (w *workoutsService) UpdateWorkoutById(context context.Context, workoutId string, t updateWorkoutRequest, userId string) error {
	_, err := w.GetById(context, workoutId, userId)
	if err != nil {
		return err
	}

	arg := repository.UpdateWorkoutByIdParams{
		Note:      t.Note,
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
		ID:        workoutId,
		UserID:    userId,
	}

	err = w.repo.UpdateById(context, arg)
	if err != nil {
		return err
	}

	return nil
}

func (w *workoutsService) CloneByIdAndReturnId(context context.Context, workoutId string, userId string) (string, error) {
	workout, err := w.GetById(context, workoutId, userId)
	if err != nil {
		return "", err
	}

	cloneId, err := w.CreateAndReturnId(context, createWorkoutRequest{
		Name: workout.Name,
	}, userId)

	if err != nil {
		return "", err
	}

	exercises, err := w.exerciseRepo.GetByWorkoutId(context, repository.GetExercisesByWorkoutIdParams{UserID: userId, WorkoutID: workoutId})

	for _, exercise := range exercises {
		uuid, err := uuid.NewV7()
		if err != nil {
			return "", err
		}

		_, err = w.exerciseRepo.CreateAndReturnId(context, repository.CreateExerciseAndReturnIdParams{
			ID:             uuid.String(),
			WorkoutID:      cloneId,
			Name:           exercise.Name,
			ExerciseTypeID: exercise.ExerciseTypeID,
			CreatedOn:      time.Now().UTC().Format(time.RFC3339),
			UserID:         userId,
			UpdatedOn:      time.Now().UTC().Format(time.RFC3339),
		})

		if err != nil {
			return "", err
		}
	}

	return cloneId, nil
}

func (w *workoutsService) DeleteById(context context.Context, workoutId string, userId string) error {
	arg := repository.DeleteWorkoutByIdParams{
		ID:     workoutId,
		UserID: userId,
	}
	return w.repo.DeleteById(context, arg)
}

func (w *workoutsService) CompleteById(context context.Context, workoutId string, userId string) error {
	completeParams := repository.CompleteWorkoutByIdParams{
		ID:          workoutId,
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UserID:      userId,
	}

	_, err := w.repo.CompleteById(context, completeParams)

	return err
}

func (w *workoutsService) CreateAndReturnId(context context.Context, t createWorkoutRequest, userId string) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	workout := repository.CreateWorkoutAndReturnIdParams{
		ID:        uuid.String(),
		Name:      t.Name,
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
		UserID:    userId,
	}

	id, err := w.repo.CreateAndReturnId(context, workout)

	return id, err
}

func (w *workoutsService) GetAllCount(context context.Context, userId string) (int, error) {
	count, err := w.repo.GetAllCount(context, userId)
	return int(count), err
}

func (w *workoutsService) GetAll(context context.Context, userId string, page int, pageSize int) ([]Workout, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	arg := repository.GetAllWorkoutsParams{
		UserID: userId,
		Limit:  int64(pageSize),
		Offset: int64((page - 1) * pageSize),
	}

	workouts, err := w.repo.GetAll(context, arg)

	if err != nil {
		return []Workout{}, err
	}

	sort.Slice(workouts, func(i, j int) bool {
		ta, err := time.Parse(time.RFC3339, workouts[i].CreatedOn)

		if err != nil {
			return false
		}

		tb, err := time.Parse(time.RFC3339, workouts[j].CreatedOn)
		if err != nil {
			return false
		}

		return tb.Before(ta)
	})

	return workouts, nil
}

func (w *workoutsService) GetById(context context.Context, id string, userId string) (Workout, error) {
	arg := repository.GetWorkoutByIdParams{
		ID:     id,
		UserID: userId,
	}
	workout, err := w.repo.GetById(context, arg)
	return workout, err
}

func NewService(repo WorkoutsRepository, exerciseRepo exercises.ExerciseRepository) Service {
	return &workoutsService{repo, exerciseRepo}
}
