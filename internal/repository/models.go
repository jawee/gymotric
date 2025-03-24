// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

type Exercise struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	CreatedOn      string      `json:"created_on"`
	UpdatedOn      string      `json:"updated_on"`
	UserID         string      `json:"user_id"`
	WorkoutID      string      `json:"workout_id"`
	ExerciseTypeID string      `json:"exercise_type_id"`
	Foreign        interface{} `json:"foreign"`
}

type ExerciseType struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	UserID    string `json:"user_id"`
}

type Set struct {
	ID          string      `json:"id"`
	Repetitions int64       `json:"repetitions"`
	Weight      float64     `json:"weight"`
	CreatedOn   string      `json:"created_on"`
	UpdatedOn   string      `json:"updated_on"`
	UserID      string      `json:"user_id"`
	ExerciseID  string      `json:"exercise_id"`
	Foreign     interface{} `json:"foreign"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

type Workout struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	CompletedOn interface{} `json:"completed_on"`
	CreatedOn   string      `json:"created_on"`
	UpdatedOn   string      `json:"updated_on"`
	UserID      string      `json:"user_id"`
	Note        interface{} `json:"note"`
}
