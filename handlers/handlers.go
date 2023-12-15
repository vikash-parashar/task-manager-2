package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/vikash-parashar/task-manager-2/controllers"
	"github.com/vikash-parashar/task-manager-2/models"
)

// @Summary Create a new task
// @Description Creates a new task with the specified details
// @ID create-task
// @Produce json
// @Param task body models.Task true "models.Task details"
// @Success 201 {object} models.Task "Successfully created task"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/create [post]
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = controllers.CreateTask(newTask)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Get a task by ID
// @Description Retrieves a task by its unique identifier
// @ID get-task
// @Produce json
// @Param id path string true "models.Task ID"
// @Success 200 {object} models.Task "Successfully retrieved task"
// @Failure 404 {object} string "models.Task not found"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/get/{id} [get]
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	task, err := controllers.GetTask(taskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// @Summary Update a task by ID
// @Description Updates a task with the specified details
// @ID update-task
// @Produce json
// @Param id path string true "models.Task ID"
// @Param task body models.Task true "Updated task details"
// @Success 200 {object} models.Task "Successfully updated task"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "models.Task not found"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/update/{id} [put]
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	var updatedTask models.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = controllers.UpdateTask(taskID, updatedTask)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a task by ID
// @Description Deletes a task by its unique identifier
// @ID delete-task
// @Produce json
// @Param id path string true "models.Task ID"
// @Success 200 {object} string "Successfully deleted task"
// @Failure 404 {object} string "models.Task not found"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/delete/{id} [delete]
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	err := controllers.DeleteTask(taskID)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get all tasks
// @Description Retrieves a list of all tasks
// @ID get-all-tasks
// @Produce json
// @Success 200 {array} models.Task "Successfully retrieved tasks"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/getAll [get]
func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := controllers.GetAllTasks()
	if err != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

// @Summary Get tasks with due reminders
// @Description Retrieves a list of tasks with due reminders
// @ID get-tasks-with-due-reminders
// @Produce json
// @Success 200 {array} models.Task "Successfully retrieved tasks with due reminders"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/dueReminders [get]
func GetTasksWithDueReminder(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	tasks, err := controllers.GetTasksWithDueReminders(currentTime)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(tasks)
}
