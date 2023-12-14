package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"github.com/vikash-parashar/task-manager-2/config"
	"github.com/vikash-parashar/task-manager-2/handlers"
	"github.com/vikash-parashar/task-manager-2/models"
)

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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/tasks/create", handlers.CreateTaskHandler)
	r.Get("/tasks/get/{id}", handlers.GetTaskHandler)
	r.Put("/tasks/update/{id}", handlers.UpdateTaskHandler)
	r.Delete("/tasks/delete/{id}", handlers.DeleteTaskHandler)
	r.Get("/tasks/getAll", handlers.GetAllTasksHandler)
	r.Get("/tasks/dueReminders", handlers.GetTasksWithDueReminder)

	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
