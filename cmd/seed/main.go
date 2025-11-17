package main

import (
	"context"
	"fmt"
	"time"
	"weight-tracker/internal/database"
	"weight-tracker/internal/repository"

	"math/rand"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Printf("Seeding\n")
	db := database.New()

	repo := db.GetRepository()

	ctx := context.Background()

	password, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	user := repository.CreateUserAndReturnIdParams{
		ID:        getUuidString(),
		Username:  "test",
		Password:  string(password),
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
	}

	userId, err := repo.CreateUserAndReturnId(ctx, user)

	if err != nil {
		panic(err)
	}

	// set user verified
	updateUser := repository.UpdateUserParams{
		ID:         userId,
		Password:   user.Password,
		UpdatedOn:  user.UpdatedOn,
		IsVerified: true,
	}

	_, err = repo.UpdateUser(ctx, updateUser)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", userId)

	exerciseTypeNames := []string{
		"Bench Press",
		"Deadlift",
		"Squats",
		"Overhead Press",
		"Pull Ups",
		"Barbell Rows",
		"Leg Press",
		"Lat Pulldown",
		"Chest Fly",
		"Tricep Pushdown",
		"Bicep Curl",
		"Leg Curl",
		"Leg Extension",
		"Calf Raise",
		"Shoulder Shrug",
		"Face Pull",
		"Plank",
		"Russian Twist",
		"Mountain Climber",
		"Burpee",
		"Jumping Jack",
		"High Knees",
		"Box Jump",
		"Battle Rope",
		"Medicine Ball Slam",
		"Wall Ball",
		"Box Squat",
		"Front Squat",
	}

	exerciseTypes := map[string]string{}
	for _, name := range exerciseTypeNames {
		id := getUuidString()
		exerciseTypes[name] = id
		exerciseType := repository.CreateExerciseTypeAndReturnIdParams{
			ID:     id,
			Name:   name,
			UserID: userId,
		}
		repo.CreateExerciseTypeAndReturnId(ctx, exerciseType)
	}

	for i := range 100 {
		createWorkout(true, 100-i, exerciseTypes, repo, ctx, userId)
	}

	for i := range 70 {
		createWorkout(true, 365-(100-i), exerciseTypes, repo, ctx, userId)
	}

	fmt.Printf("Done\n")
}

func createWorkout(completed bool, daysAgo int, exerciseTypes map[string]string, repo repository.Querier, ctx context.Context, userId string) {
	workoutNames := []string{
		"Legs",
		"Back",
		"Chest",
	}
	r := rand.Intn(len(workoutNames))
	workout := repository.CreateWorkoutAndReturnIdParams{
		ID:        getUuidString(),
		Name:      workoutNames[r],
		CreatedOn: time.Now().UTC().AddDate(0, 0, -daysAgo).Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().AddDate(0, 0, -daysAgo).Add(time.Hour).Format(time.RFC3339),
		UserID:    userId,
	}

	repo.CreateWorkoutAndReturnId(ctx, workout)

	numberOfExercises := rand.Intn(5) + 1

	for range numberOfExercises {
		exerciseTypeName := getRandomExerciseType(exerciseTypes)
		exerciseTypeId := exerciseTypes[exerciseTypeName]
		createExerciseWithSets(workout.ID, exerciseTypeName, exerciseTypeId, repo, ctx, userId)
	}

	if completed {
		setCompleted := repository.CompleteWorkoutByIdParams{
			CompletedOn: time.Now().UTC().AddDate(0, 0, -daysAgo).Add(time.Duration(time.Hour)).Format(time.RFC3339),
			ID:          workout.ID,
			UserID:      userId,
		}

		repo.CompleteWorkoutById(ctx, setCompleted)
	}
}
func getRandomExerciseType(m map[string]string) string {
	r := rand.Intn(len(m))
	for k := range m {
		if r == 0 {
			return k
		}
		r--
	}
	panic("should not happen")
}

func createExerciseWithSets(workoutId string, exerciseTypeName string, exerciseTypeId string, repo repository.Querier, ctx context.Context, userId string) {
	// Create exercise_item first
	exerciseItemId := getUuidString()
	now := time.Now().UTC().Format(time.RFC3339)
	exerciseItem := repository.CreateExerciseItemAndReturnIdParams{
		ID:        exerciseItemId,
		Type:      "exercise",
		UserID:    userId,
		WorkoutID: workoutId,
		CreatedOn: now,
		UpdatedOn: now,
	}
	repo.CreateExerciseItemAndReturnId(ctx, exerciseItem)

	// Create exercise with the exercise_item_id
	exercise := repository.CreateExerciseAndReturnIdParams{
		ID:             getUuidString(),
		Name:           exerciseTypeName,
		WorkoutID:      workoutId,
		ExerciseTypeID: exerciseTypeId,
		ExerciseItemID: exerciseItemId,
		CreatedOn:      now,
		UpdatedOn:      now,
		UserID:         userId,
	}

	repo.CreateExerciseAndReturnId(ctx, exercise)

	numberOfSets := rand.Intn(2) + 2
	for range numberOfSets {
		set := repository.CreateSetAndReturnIdParams{
			ID:          getUuidString(),
			Repetitions: 9,
			Weight:      110,
			ExerciseID:  exercise.ID,
			UserID:      userId,
		}
		repo.CreateSetAndReturnId(ctx, set)
	}
}

func getUuidString() string {
	id, _ := uuid.NewV7()
	return id.String()
}
