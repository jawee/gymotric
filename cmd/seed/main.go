package main

import (
	"context"
	"fmt"
	"time"
	"weight-tracker/internal/database"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

func main() {
	fmt.Printf("Seeding\n")
	db := database.New()

	repo := db.GetRepository()

	ctx := context.Background()
	workout := repository.CreateWorkoutAndReturnIdParams{
		ID:        getUuidString(),
		Name:      "back",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	exerciseType := repository.CreateExerciseTypeAndReturnIdParams{
		ID: getUuidString(),
		Name: "Deadlift",
	}

	exercise := repository.CreateExerciseAndReturnIdParams {
		ID: getUuidString(),
		Name: exerciseType.Name,
		WorkoutID: workout.ID,
		ExerciseTypeID: exerciseType.ID,
	}

	set := repository.CreateSetAndReturnIdParams {
		ID: getUuidString(),
		Repetitions: 8,
		Weight: 110,
		ExerciseID: exercise.ID,
	}
	set2 := repository.CreateSetAndReturnIdParams {
		ID: getUuidString(),
		Repetitions: 8,
		Weight: 110,
		ExerciseID: exercise.ID,
	}

	repo.CreateWorkoutAndReturnId(ctx, workout)
	repo.CreateExerciseTypeAndReturnId(ctx, exerciseType)
	repo.CreateExerciseAndReturnId(ctx, exercise)
	repo.CreateSetAndReturnId(ctx, set)
	repo.CreateSetAndReturnId(ctx, set2)

	fmt.Printf("Done\n")
}

func getUuidString() string {
	id, _ := uuid.NewV7()
	return id.String()
}
