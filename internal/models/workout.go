package models

import "time"

// Workout represents a workout plan entry
type Workout struct {
	ID          int64
	SportType   string
	Duration    int     // in minutes
	Distance    float64 // in kilometers
	Date        time.Time
	IsCompleted bool
	Notes       string
}
