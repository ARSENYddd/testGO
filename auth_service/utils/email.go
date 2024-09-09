package utils

import (
	"log"
	"net/smtp"
)

func SendEmailWarning(userId, oldIP, newIP string) error {
	from := "your-email@example.com"
	to := "user-email@example.com"
	subject := "IP Address Change Warning"
	body := "Hello,\n\nYour IP address has changed from " + oldIP + " to " + newIP + ".\nIf this was not you, please take appropriate action.\n\nBest regards,\nYour App"

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail("smtp.example.com:587",
		smtp.PlainAuth("", from, "your-email-password", "smtp.example.com"),
		from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}
