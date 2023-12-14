package helpers

import (
	"fmt"
	"log"
	"time"

	"github.com/vikash-parashar/task-manager-2/controllers"
	"github.com/vikash-parashar/task-manager-2/models"
)

func CheckReminders() {
	for {
		// Get current time
		currentTime := time.Now()

		// Query tasks with reminders due
		tasks, err := controllers.GetTasksWithDueReminders(currentTime)
		if err != nil {
			log.Println("Error querying tasks with due reminders:", err)
		}

		// Check and send notifications for each task
		for _, task := range tasks {
			sendNotification(task)
		}

		// Sleep for a specific duration before checking again
		time.Sleep(time.Minute)
	}
}

// sendNotification checks the notification method and sends the notification
func sendNotification(task models.Task) {
	switch task.NotifyMethod {
	case "email":
		sendEmailNotification(task)
	case "push":
		sendPushNotification(task)
	default:
		log.Printf("Unknown notification method for task %s: %s", task.ID, task.NotifyMethod)
	}
}

// sendEmailNotification sends an email notification for the task
func sendEmailNotification(task models.Task) {
	//TODO:
	//FIXME:
	fmt.Printf("Sending email notification for task %s\n", task.ID)
}

// sendPushNotification sends a push notification for the task
func sendPushNotification(task models.Task) {
	//TODO:
	//FIXME:
	fmt.Printf("Sending push notification for task %s\n", task.ID)
}
