package exerciseitems

type ExerciseItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	UserID    string `json:"user_id"`
	WorkoutID string `json:"workout_id"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}
