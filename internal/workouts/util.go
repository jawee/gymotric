package workouts

import "github.com/google/uuid"

func generateUuid() string {
	id, _ := uuid.NewV7()
	return id.String()
}
