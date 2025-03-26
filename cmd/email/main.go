package main

import (
	"log/slog"
	"weight-tracker/internal/email"
)

func main() {
	err := email.SendEmail()
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return
	}
	slog.Info("Email sent successfully")
}
