package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/guisithos/workout-planner/internal/database"
	"github.com/guisithos/workout-planner/internal/models"
)

type Handler struct {
	store *database.WorkoutStore
	tmpl  *template.Template
}

func NewHandler(store *database.WorkoutStore, templates *template.Template) *Handler {
	return &Handler{
		store: store,
		tmpl:  templates,
	}
}

// HandleHome displays the main calendar view
func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get the current month's start and end dates
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	workouts, err := h.store.GetWorkouts(startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to fetch workouts", http.StatusInternalServerError)
		return
	}

	data := struct {
		Workouts []models.Workout
		Month    time.Time
	}{
		Workouts: workouts,
		Month:    startDate,
	}

	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleNewWorkout displays the workout creation form
func (h *Handler) HandleNewWorkout(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "workout_form.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleSaveWorkout processes the workout form submission
func (h *Handler) HandleSaveWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	duration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		http.Error(w, "Invalid duration", http.StatusBadRequest)
		return
	}

	distance, err := strconv.ParseFloat(r.FormValue("distance"), 64)
	if err != nil {
		http.Error(w, "Invalid distance", http.StatusBadRequest)
		return
	}

	workout := &models.Workout{
		SportType: r.FormValue("sport_type"),
		Duration:  duration,
		Distance:  distance,
		Date:      date,
		Notes:     r.FormValue("notes"),
	}

	if err := h.store.CreateWorkout(workout); err != nil {
		http.Error(w, "Failed to save workout", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleToggleWorkout toggles the completion status of a workout
func (h *Handler) HandleToggleWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	if err := h.store.ToggleWorkoutCompletion(id); err != nil {
		http.Error(w, "Failed to toggle workout completion", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
