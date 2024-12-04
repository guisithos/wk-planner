package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/guisithos/wk-planner/internal/database"
	"github.com/guisithos/wk-planner/internal/models"
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

// CalendarDay represents a single day in the calendar view
type CalendarDay struct {
	Date     time.Time
	Workouts []models.Workout
	IsToday  bool
}

// HandleHome displays the main calendar view
func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	now := time.Now()
	year := now.Year()
	month := now.Month()

	// Get the first day of the month and the total days
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	// Get all workouts for the month
	workouts, err := h.store.GetWorkouts(firstDay, lastDay)
	if err != nil {
		http.Error(w, "Failed to fetch workouts", http.StatusInternalServerError)
		return
	}

	// Create calendar grid
	var calendar [][]CalendarDay

	// Get the day of week for the first day (0 = Sunday, 1 = Monday, etc.)
	firstDayWeekday := int(firstDay.Weekday())

	// Create calendar weeks
	currentDay := firstDay.AddDate(0, 0, -firstDayWeekday) // Start from the first Sunday
	for week := 0; week < 6; week++ {                      // Maximum 6 weeks in a month
		var weekDays []CalendarDay
		for day := 0; day < 7; day++ {
			// Create calendar day
			calDay := CalendarDay{
				Date:    currentDay,
				IsToday: currentDay.Year() == now.Year() && currentDay.Month() == now.Month() && currentDay.Day() == now.Day(),
			}

			// Add workouts for this day
			for _, w := range workouts {
				if w.Date.Year() == currentDay.Year() &&
					w.Date.Month() == currentDay.Month() &&
					w.Date.Day() == currentDay.Day() {
					calDay.Workouts = append(calDay.Workouts, w)
				}
			}

			weekDays = append(weekDays, calDay)
			currentDay = currentDay.AddDate(0, 0, 1)
		}
		calendar = append(calendar, weekDays)

		// Break if we've gone past the end of the month
		if currentDay.Month() != month && week >= 4 {
			break
		}
	}

	data := struct {
		Calendar [][]CalendarDay
		Month    time.Time
	}{
		Calendar: calendar,
		Month:    firstDay,
	}

	if err := h.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
