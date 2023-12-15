package controllers

import (
	"fmt"
	"time"

	"github.com/vikash-parashar/task-manager-2/config"
	"github.com/vikash-parashar/task-manager-2/models"
)

// @Summary Create a new task
// @Description Adds a new task to the database
// @ID create-task
// @Accept json
// @Produce json
// @Param task body models.Task true "Task details"
// @Success 200 {string} string "Successfully created task"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/create [post]
// CreateTask creates a new task in the database
func CreateTask(task models.Task) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Insert task
	_, err = tx.Exec("INSERT INTO tasks (id, title, description, priority, due_date_time, email) VALUES ($1, $2, $3, $4, $5, $6)",
		task.ID, task.Title, task.Description, task.Priority, task.DueDateTime, task.Email)
	if err != nil {
		return fmt.Errorf("failed to insert task: %v", err)
	}

	// Insert reminders
	for _, reminder := range task.Reminders {
		_, err = tx.Exec("INSERT INTO reminders (id, date, task_id) VALUES ($1, $2, $3)", reminder.ID, reminder.Date, task.ID)
		if err != nil {
			return fmt.Errorf("failed to insert reminder: %v", err)
		}
	}

	return nil
}

// @Summary Get a task by ID
// @Description Retrieves a task from the database by ID
// @ID get-task
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task "Successfully retrieved task"
// @Failure 404 {object} string "Task not found"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/get/{id} [get]
func GetTask(id string) (models.Task, error) {
	var task models.Task

	row := config.DB.QueryRow("SELECT id, title, description, priority, due_date_time , email FROM tasks WHERE id = $1", id)
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDateTime, &task.Email)
	if err != nil {
		return models.Task{}, err
	}

	rows, err := config.DB.Query("SELECT id, date FROM reminders WHERE task_id = $1", id)
	if err != nil {
		return models.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var reminder models.Reminder
		err := rows.Scan(&reminder.ID, &reminder.Date)
		if err != nil {
			return models.Task{}, err
		}
		task.Reminders = append(task.Reminders, reminder)
	}

	return task, nil
}

// @Summary Update a task by ID
// @Description Updates an existing task in the database
// @ID update-task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.Task true "Updated task details"
// @Success 200 {string} string "Successfully updated task"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/update/{id} [put]
// UpdateTask updates an existing task in the database
func UpdateTask(id string, task models.Task) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Update task
	_, err = tx.Exec("UPDATE tasks SET title = $1, description = $2, priority = $3, due_date_time = $4, email = $5 WHERE id = $6",
		task.Title, task.Description, task.Priority, task.DueDateTime, task.Email, id)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	// Delete existing reminders
	_, err = tx.Exec("DELETE FROM reminders WHERE task_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete existing reminders: %v", err)
	}

	// Insert updated reminders
	for _, reminder := range task.Reminders {
		_, err = tx.Exec("INSERT INTO reminders (id, date, task_id) VALUES ($1, $2, $3)", reminder.ID, reminder.Date, id)
		if err != nil {
			return fmt.Errorf("failed to insert updated reminder: %v", err)
		}
	}

	return nil
}

// @Summary Delete a task by ID
// @Description Removes a task from the database by ID
// @ID delete-task
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {string} string "Successfully deleted task"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/delete/{id} [delete]
// DeleteTask deletes a task from the database by ID
func DeleteTask(id string) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Delete associated reminders
	_, err = tx.Exec("DELETE FROM reminders WHERE task_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete associated reminders: %v", err)
	}

	// Delete task
	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

// @Summary Get all tasks
// @Description Retrieves a list of all tasks from the database
// @ID get-all-tasks
// @Produce json
// @Success 200 {array} models.Task "Successfully retrieved tasks"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/getAll [get]
func GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task

	rows, err := config.DB.Query("SELECT id, title, description, priority, due_date_time ,email FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDateTime, &task.Email)
		if err != nil {
			return nil, err
		}

		remindersRows, err := config.DB.Query("SELECT id, date FROM reminders WHERE task_id = $1", task.ID)
		if err != nil {
			return nil, err
		}
		defer remindersRows.Close()

		for remindersRows.Next() {
			var reminder models.Reminder
			err := remindersRows.Scan(&reminder.ID, &reminder.Date)
			if err != nil {
				return nil, err
			}
			task.Reminders = append(task.Reminders, reminder)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// @Summary Get tasks with due reminders
// @Description Retrieves a list of tasks with due reminders from the database
// @ID get-tasks-with-due-reminders
// @Produce json
// @Param currentTime query string true "Current time in RFC3339 format"
// @Success 200 {array} models.Task "Successfully retrieved tasks with due reminders"
// @Failure 500 {object} string "Internal server error"
// @Router /tasks/dueReminders [get]
func GetTasksWithDueReminders(currentTime time.Time) ([]models.Task, error) {

	rows, err := config.DB.Query("SELECT id, title, description, priority, due_date_time,email FROM tasks WHERE due_date_time <= $1", currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDateTime, &task.Email)
		if err != nil {
			return nil, err
		}

		// Query reminders for the task
		remindersRows, err := config.DB.Query("SELECT id, date FROM reminders WHERE task_id = $1", task.ID)
		if err != nil {
			return nil, err
		}
		defer remindersRows.Close()

		for remindersRows.Next() {
			var reminder models.Reminder
			err := remindersRows.Scan(&reminder.ID, &reminder.Date)
			if err != nil {
				return nil, err
			}
			task.Reminders = append(task.Reminders, reminder)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
