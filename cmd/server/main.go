package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/guisithos/wk-planner/internal/database"
	"github.com/guisithos/wk-planner/internal/handlers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize database
	db, err := sql.Open("sqlite3", "./workouts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables if they don't exist
	initDB(db)

	// Load templates with debug logging
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	log.Printf("Loaded templates: %v", tmpl.DefinedTemplates())

	// Initialize handlers
	store := database.NewWorkoutStore(db)
	handler := handlers.NewHandler(store, tmpl)

	// Initialize handlers
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register routes
	mux.HandleFunc("/", handler.HandleHome)
	mux.HandleFunc("/workout/new", handler.HandleNewWorkout)
	mux.HandleFunc("/workout/save", handler.HandleSaveWorkout)
	mux.HandleFunc("/workout/toggle", handler.HandleToggleWorkout)

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
