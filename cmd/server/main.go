package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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

func main() {
	// Initialize database
	db, err := sql.Open("sqlite3", "./workouts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables if they don't exist
	initDB(db)

	// Initialize handlers
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/workout/new", handleNewWorkout)
	mux.HandleFunc("/workout/save", handleSaveWorkout)
	mux.HandleFunc("/workout/toggle", handleToggleWorkout)

	// Start server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func initDB(db *sql.DB) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS workouts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sport_type TEXT NOT NULL,
		duration INTEGER NOT NULL,
		distance REAL NOT NULL,
		date DATE NOT NULL,
		is_completed BOOLEAN DEFAULT FALSE,
		notes TEXT
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
