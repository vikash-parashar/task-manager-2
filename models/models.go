package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Task represents a task with its details
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	DueDateTime time.Time  `json:"dueDateTime"`
	Reminders   []Reminder `json:"reminders"`

	// Notification fields
	NotifyMethod  string `json:"notifyMethod"`  // e.g., "email", "push"
	NotifyStatus  string `json:"notifyStatus"`  // e.g., "pending", "sent", "failed"
	NotifyMessage string `json:"notifyMessage"` // Additional information about the notification
}

// Reminder represents a reminder associated with a task
type Reminder struct {
	ID     string `json:"id"`
	Date   string `json:"date"`
	TaskID string `json:"taskID"`
}

// @Summary Create tables
// @Description Creates the Task and Reminder tables in the database
// @ID create-tables
// @Success 200 {string} string "Tables created successfully"
// @Failure 500 {object} string "Internal server error"
// @Router /createTables [post]
func CreateTables(db *sql.DB) error {
	// Create Task table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id VARCHAR(36) PRIMARY KEY,
			title VARCHAR(255),
			description VARCHAR(255),
			priority VARCHAR(50),
			due_date_time TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create Task table: %v", err)
	}

	// Create Reminder table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reminders (
			id VARCHAR(36) PRIMARY KEY,
			date VARCHAR(20),
			task_id VARCHAR(36) REFERENCES tasks(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create Reminder table: %v", err)
	}

	fmt.Println("Tables created successfully")
	return nil
}
