package helpers

import (
	"fmt"
	"net/smtp"

	"github.com/vikash-parashar/task-manager-2/models"
)

// sendEmailNotification sends an email notification for a task
func SendEmailNotification(task models.Task, message string) error {
	// Replace these with your email sending configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := "gowithvikash@gmail.com"
	smtpPassword := "pfzucjdducunohbd"

	// Set up authentication information
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Compose the email message
	to := []string{task.Email}
	subject := "Task Notification"
	body := fmt.Sprintf("Subject: %s\n\n%s - %s", subject, message, task.Title)

	// Send the email
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, smtpUsername, to, []byte(body))
	if err != nil {
		return err
	}

	fmt.Printf("Sent email to %s: %s - %s\n", task.Email, message, task.Title)
	return nil
}
