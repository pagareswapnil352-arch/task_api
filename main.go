package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ─────────────────────────────────────────
// 1. TASK — defines what a task looks like
// ─────────────────────────────────────────

type Task struct {
	ID    int    `json:"id"`    // Task ID
	Title string `json:"title"` // Task title
}

// ─────────────────────────────────────────
// 2. TASK STORE — stores all tasks
// ─────────────────────────────────────────

type TaskStore struct {
	tasks  map[int]Task // Stores tasks using ID as key
	nextID int          // Keeps track of the next ID
}

// ─────────────────────────────────────────
// 3. CONSTRUCTOR — creates an empty store
// ─────────────────────────────────────────

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task), // Create an empty map
		nextID: 1,                  // First task starts with ID 1
	}
}

// ─────────────────────────────────────────
// 4. CREATE — creates a new task
// ─────────────────────────────────────────

func (s *TaskStore) Create(title string) Task {

	task := Task{ID: s.nextID, Title: title} // Create a new task

	s.tasks[task.ID] = task // Save task in the map

	s.nextID++ // Increase ID for the next task

	return task // Return the created task
}

// ─────────────────────────────────────────
// 5. GET — retrieves one task using its ID
// ─────────────────────────────────────────

func (s *TaskStore) Get(id int) Task {

	return s.tasks[id] // Return task from the map
}

// ─────────────────────────────────────────
// 6. LIST — returns all tasks
// ─────────────────────────────────────────

func (s *TaskStore) List() []Task {

	result := []Task{} // Create an empty slice

	for _, task := range s.tasks { // Loop through all tasks

		result = append(result, task) // Add task to the slice
	}

	return result // Return all tasks
}

// ─────────────────────────────────────────
// 7. UPDATE — updates an existing task
// ─────────────────────────────────────────

func (s *TaskStore) Update(id int, title string) Task {

	task, exists := s.tasks[id] // Check whether task exists

	if !exists {
		return Task{} // Return empty task if not found
	}

	task.Title = title // Change the task title

	s.tasks[id] = task // Save the updated task

	return task // Return updated task
}

// ─────────────────────────────────────────
// 8. DELETE — removes a task
// ─────────────────────────────────────────

func (s *TaskStore) Delete(id int) {

	delete(s.tasks, id) // Remove task using its ID
}

// ─────────────────────────────────────────
// 9. MAIN — starts the HTTP server
// ─────────────────────────────────────────

func main() {

	store := NewTaskStore() // Create task storage

	mux := http.NewServeMux() // Create HTTP request router

	// ─────────────────────────────────────
	// GET /tasks — get all tasks
	// ─────────────────────────────────────

	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(store.List()) // Return all tasks as JSON
	})

	// ─────────────────────────────────────
	// POST /tasks — create a new task
	// ─────────────────────────────────────

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {

		var body struct {
			Title string `json:"title"` // Receive title from JSON
		}

		// Read JSON from the incoming request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		task := store.Create(body.Title) // Create and save the task

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated) // Send 201 Created status

		json.NewEncoder(w).Encode(task) // Return created task
	})

	// ─────────────────────────────────────
	// GET /tasks/{id} — get one task
	// ─────────────────────────────────────

	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {

		// Get ID from URL and convert string into integer
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		task := store.Get(id) // Find task using its ID

		if task.ID == 0 {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(task) // Return task as JSON
	})

	// ─────────────────────────────────────
	// PUT /tasks/{id} — update a task
	// ─────────────────────────────────────

	mux.HandleFunc("PUT /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {

		// Get ID from URL and convert it into integer
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var body struct {
			Title string `json:"title"` // Receive new title
		}

		// Read new title from request body
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		task := store.Update(id, body.Title) // Update the task

		if task.ID == 0 {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(task) // Return updated task
	})

	// ─────────────────────────────────────
	// DELETE /tasks/{id} — delete a task
	// ─────────────────────────────────────

	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {

		// Get ID from URL and convert it into integer
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		// Check whether task exists before deleting
		if store.Get(id).ID == 0 {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		store.Delete(id) // Delete task from the map

		w.WriteHeader(http.StatusNoContent) // Send 204 No Content
	})

	// ─────────────────────────────────────
	// START THE SERVER
	// ─────────────────────────────────────

	fmt.Println("Server running on http://localhost:8080")

	// Start HTTP server on port 8080
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Server error:", err)
	}
}
