package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vikash-parashar/task-manager-2/config"
	"github.com/vikash-parashar/task-manager-2/handlers"
	"github.com/vikash-parashar/task-manager-2/helpers"
	"github.com/vikash-parashar/task-manager-2/models"
)

var done = make(chan bool)

// @title Task API
// @version 1.0
// @description API for managing tasks and reminders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@task-api.com
// @host localhost:8080
// @BasePath /v1
func main() {
	// Initialize the database connection
	err := config.InitDB()
	if err != nil {
		panic(err)
	}

	// Create tables
	err = models.CreateTables(config.DB)
	if err != nil {
		panic(err)
	}

	// PeriodicTask checks for due tasks in the background and sends email notifications.
	// @title Background Task
	// @description This background task runs periodically to find and notify due tasks.
	// @version 1.0
	// @termsOfService http://swagger.io/terms/
	// @contact.name API Support
	// @contact.email gowithvikash@gmail.com
	// @host localhost:8080
	// @BasePath /v1
	go func() {
		ticker := time.NewTicker(30 * time.Minute) // Adjust the interval as needed
		defer ticker.Stop()

		for range ticker.C {
			// Execute the function to find and notify due tasks
			go func() {
				err := FindAndNotifyDueTasks()
				if err != nil {
					fmt.Println("Error while checking and notifying due tasks:", err)
				}
			}()
		}
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Allow all CORS requests
	corsHandler := cors.AllowAll()
	// Use CORS middleware
	r.Use(corsHandler.Handler)

	// Serve the Swagger UI at /swagger/index.html
	r.Get("/swagger/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))).ServeHTTP(w, r)
	}))

	// Swagger configuration
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	// @Summary Create a new task
	// @Description Creates a new task with the specified details
	// @ID create-task
	// @Produce json
	// @Param task body Task true "Task details"
	// @Success 201 {object} Task "Successfully created task"
	// @Failure 400 {object} ErrorResponse "Bad request"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	// @Router /tasks/create [post]
	r.Post("/tasks/create", handlers.CreateTaskHandler)

	// @Summary Get a task by ID
	// @Description Retrieves a task by its unique identifier
	// @ID get-task
	// @Produce json
	// @Param id path string true "Task ID"
	// @Success 200 {object} Task "Successfully retrieved task"
	// @Failure 404 {object} ErrorResponse "Task not found"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	// @Router /tasks/get/{id} [get]
	r.Get("/tasks/get/{id}", handlers.GetTaskHandler)

	// @Summary Update a task by ID
	// @Description Updates a task with the specified details
	// @ID update-task
	// @Produce json
	// @Param id path string true "Task ID"
	// @Param task body Task true "Task details"
	// @Success 200 {object} Task "Successfully updated task"
	// @Failure 400 {object} ErrorResponse "Bad request"
	// @Failure 404 {object} ErrorResponse "Task not found"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	// @Router /tasks/update/{id} [put]
	r.Put("/tasks/update/{id}", handlers.UpdateTaskHandler)

	// @Summary Delete a task by ID
	// @Description Deletes a task by its unique identifier
	// @ID delete-task
	// @Produce json
	// @Param id path string true "Task ID"
	// @Success 200 {object} SuccessResponse "Successfully deleted task"
	// @Failure 404 {object} ErrorResponse "Task not found"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	// @Router /tasks/delete/{id} [delete]
	r.Delete("/tasks/delete/{id}", handlers.DeleteTaskHandler)

	// @Summary Get all tasks
	// @Description Retrieves a list of all tasks
	// @ID get-all-tasks
	// @Produce json
	// @Success 200 {array} Task "Successfully retrieved tasks"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	// @Router /tasks/getAll [get]
	r.Get("/tasks/getAll", handlers.GetAllTasksHandler)

	// @Summary Get tasks with due reminders
	// @Description Retrieves a list of tasks with due reminders
	// @ID get-tasks-with-due-reminders
	// @Produce json
	// @Success 200 {array} Task "Successfully retrieved tasks with due reminders"
	// @Failure 500 {object} ErrorResponse "Internal server error"
	// @Router /tasks/dueReminders [get]
	r.Get("/tasks/dueReminders", handlers.GetTasksWithDueReminderHandler)

	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// FindAndNotifyDueTasks checks for tasks that are due based on reminders and sends email notifications
func FindAndNotifyDueTasks() error {
	defer func() {
		done <- true
	}()
	currentTime := time.Now()

	// Retrieve tasks with associated reminders
	rows, err := config.DB.Query("SELECT tasks.id, tasks.title, tasks.description, tasks.priority, tasks.due_date_time, tasks.email, reminders.id, reminders.date FROM tasks LEFT JOIN reminders ON tasks.id = reminders.task_id")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		var reminder models.Reminder
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDateTime, &task.Email, &reminder.ID, &reminder.Date)
		if err != nil {
			return err
		}

		// Check if the reminder is due
		if !reminder.Date.IsZero() && reminder.Date.Before(currentTime) {
			// Reminder is due
			helpers.SendEmailNotification(task, "Reminder due")
		}
	}

	return nil
}
