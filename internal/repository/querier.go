// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"context"
)

type Querier interface {
	CompleteWorkoutById(ctx context.Context, arg CompleteWorkoutByIdParams) (int64, error)
	CreateExerciseAndReturnId(ctx context.Context, arg CreateExerciseAndReturnIdParams) (string, error)
	CreateExerciseTypeAndReturnId(ctx context.Context, arg CreateExerciseTypeAndReturnIdParams) (string, error)
	CreateSetAndReturnId(ctx context.Context, arg CreateSetAndReturnIdParams) (string, error)
	CreateUserAndReturnId(ctx context.Context, arg CreateUserAndReturnIdParams) (string, error)
	CreateWorkoutAndReturnId(ctx context.Context, arg CreateWorkoutAndReturnIdParams) (string, error)
	DeleteExerciseById(ctx context.Context, arg DeleteExerciseByIdParams) (int64, error)
	DeleteExerciseTypeById(ctx context.Context, arg DeleteExerciseTypeByIdParams) (int64, error)
	DeleteSetById(ctx context.Context, arg DeleteSetByIdParams) (int64, error)
	DeleteWorkoutById(ctx context.Context, arg DeleteWorkoutByIdParams) (int64, error)
	GetAllExerciseTypes(ctx context.Context, userID string) ([]ExerciseType, error)
	GetAllExercises(ctx context.Context, userID string) ([]Exercise, error)
	GetAllSets(ctx context.Context, userID string) ([]Set, error)
	GetAllWorkouts(ctx context.Context, userID string) ([]Workout, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetExerciseById(ctx context.Context, arg GetExerciseByIdParams) (Exercise, error)
	GetExerciseTypeById(ctx context.Context, arg GetExerciseTypeByIdParams) (ExerciseType, error)
	GetExercisesByWorkoutId(ctx context.Context, arg GetExercisesByWorkoutIdParams) ([]Exercise, error)
	GetLastWeightRepsByExerciseTypeId(ctx context.Context, arg GetLastWeightRepsByExerciseTypeIdParams) (GetLastWeightRepsByExerciseTypeIdRow, error)
	GetMaxWeightRepsByExerciseTypeId(ctx context.Context, arg GetMaxWeightRepsByExerciseTypeIdParams) (GetMaxWeightRepsByExerciseTypeIdRow, error)
	GetSetById(ctx context.Context, arg GetSetByIdParams) (Set, error)
	GetSetsByExerciseId(ctx context.Context, arg GetSetsByExerciseIdParams) ([]Set, error)
	GetStatisticsSinceDate(ctx context.Context, arg GetStatisticsSinceDateParams) (int64, error)
	GetWorkoutById(ctx context.Context, arg GetWorkoutByIdParams) (Workout, error)
	UpdateWorkoutById(ctx context.Context, arg UpdateWorkoutByIdParams) (int64, error)
}

var _ Querier = (*Queries)(nil)
