package main

import (
	"context"
	"fmt"
	"time"
	"weight-tracker/internal/database"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Printf("Seeding\n")
	db := database.New()

	repo := db.GetRepository()

	ctx := context.Background()

	password, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	user := repository.CreateUserAndReturnIdParams {
		ID: getUuidString(),
		Username: "test",
		Password: string(password),
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
	}

	id, err := repo.CreateUserAndReturnId(ctx, user)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", id)

	exerciseType := repository.CreateExerciseTypeAndReturnIdParams{
		ID:   getUuidString(),
		Name: "Deadlift",
	}
	repo.CreateExerciseTypeAndReturnId(ctx, exerciseType)
	repo.CreateExerciseTypeAndReturnId(ctx, repository.CreateExerciseTypeAndReturnIdParams{ ID: getUuidString(), Name: "Squats" })

	createWorkout(false, 0, exerciseType, repo, ctx)
	createWorkout(true, 2, exerciseType, repo, ctx)

	fmt.Printf("Done\n")
}

func createWorkout(completed bool, daysAgo int, exerciseType repository.CreateExerciseTypeAndReturnIdParams, repo repository.Querier, ctx context.Context) {
	workout := repository.CreateWorkoutAndReturnIdParams{
		ID:        getUuidString(),
		Name:      "back",
		CreatedOn: time.Now().UTC().AddDate(0, 0, -daysAgo).Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
	}

	exercise := repository.CreateExerciseAndReturnIdParams{
		ID:             getUuidString(),
		Name:           exerciseType.Name,
		WorkoutID:      workout.ID,
		ExerciseTypeID: exerciseType.ID,
	}

	set := repository.CreateSetAndReturnIdParams{
		ID:          getUuidString(),
		Repetitions: 8,
		Weight:      110,
		ExerciseID:  exercise.ID,
	}
	set2 := repository.CreateSetAndReturnIdParams{
		ID:          getUuidString(),
		Repetitions: 8,
		Weight:      110,
		ExerciseID:  exercise.ID,
	}

	repo.CreateWorkoutAndReturnId(ctx, workout)
	repo.CreateExerciseAndReturnId(ctx, exercise)
	repo.CreateSetAndReturnId(ctx, set)
	repo.CreateSetAndReturnId(ctx, set2)

	if completed {
		setCompleted := repository.CompleteWorkoutByIdParams {
			CompletedOn: time.Now().UTC().Format(time.RFC3339),
			ID: workout.ID,
		}

		repo.CompleteWorkoutById(ctx, setCompleted)
	}

}

func getUuidString() string {
	id, _ := uuid.NewV7()
	return id.String()
}
