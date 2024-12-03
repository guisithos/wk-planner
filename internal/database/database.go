package database

import (
	"database/sql"
	"time"

	"github.com/guisithos/workout-planner/internal/models"
)

// WorkoutStore handles all database operations for workouts
type WorkoutStore struct {
	db *sql.DB
}

// NewWorkoutStore creates a new WorkoutStore
func NewWorkoutStore(db *sql.DB) *WorkoutStore {
	return &WorkoutStore{db: db}
}

// CreateWorkout inserts a new workout into the database
func (s *WorkoutStore) CreateWorkout(workout *models.Workout) error {
	query := `
		INSERT INTO workouts (sport_type, duration, distance, date, is_completed, notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query,
		workout.SportType,
		workout.Duration,
		workout.Distance,
		workout.Date,
		workout.IsCompleted,
		workout.Notes,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	workout.ID = id
	return nil
}

// GetWorkouts retrieves all workouts for a given time range
func (s *WorkoutStore) GetWorkouts(startDate, endDate time.Time) ([]models.Workout, error) {
	query := `
		SELECT id, sport_type, duration, distance, date, is_completed, notes
		FROM workouts
		WHERE date BETWEEN ? AND ?
		ORDER BY date DESC
	`
	rows, err := s.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []models.Workout
	for rows.Next() {
		var w models.Workout
		err := rows.Scan(
			&w.ID,
			&w.SportType,
			&w.Duration,
			&w.Distance,
			&w.Date,
			&w.IsCompleted,
			&w.Notes,
		)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	return workouts, nil
}

// ToggleWorkoutCompletion updates the completion status of a workout
func (s *WorkoutStore) ToggleWorkoutCompletion(id int64) error {
	query := `
		UPDATE workouts
		SET is_completed = NOT is_completed
		WHERE id = ?
	`
	_, err := s.db.Exec(query, id)
	return err
}
