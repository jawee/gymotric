package migrations

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)
func init() {
	goose.AddMigrationNoTxContext(Up0008, Down0008)
}

type Exercise struct {
	ID             string         `db:"id"`
	UserID         string         `db:"user_id"`
	CreatedOn      string         `db:"created_on"`
	UpdatedOn      string         `db:"updated_on"`
	WorkoutID      string         `db:"workout_id"`
}

func fetchExercises(db *sql.DB) ([]Exercise, error) {
	rows, err := db.Query("SELECT id, user_id, created_on, updated_on, workout_id FROM exercises ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []Exercise
	for rows.Next() {
		var ex Exercise
		if err := rows.Scan(&ex.ID, &ex.UserID, &ex.CreatedOn, &ex.UpdatedOn, &ex.WorkoutID); err != nil {
			return nil, err
		}
		exercises = append(exercises, ex)
	}
	return exercises, nil
}
func Up0008(ctx context.Context, db *sql.DB) error {
	// Get all exercises, order by Id (for ordering)
	exercises, err := fetchExercises(db)

	if err != nil {
		return err
	}

	// loop through all
	t := "straight"
	for _, e := range(exercises) {

		id, err := uuid.NewV7()
		if err != nil {
			return err
		}

		// insert into exercise_items
		_, err = db.ExecContext(ctx, "INSERT INTO exercise_items (id, user_id, created_on, updated_on, workout_id, type) VALUES (?, ?, ?, ?, ?, ?)", id.String(), e.UserID, e.CreatedOn, e.UpdatedOn, e.WorkoutID, t)
		if err != nil {
			return err
		}

		// update exercise to set exercise_item_id
		_, err = db.ExecContext(ctx, "UPDATE exercises SET exercise_item_id = ? WHERE id = ?", id.String(), e.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func Down0008(ctx context.Context, db *sql.DB) error {
	return nil
}
