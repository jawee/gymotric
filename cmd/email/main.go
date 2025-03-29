package main

import (
	"embed"
	"log/slog"
	"weight-tracker/internal/email"
)

var embedMigrations embed.FS
func main() {
	err := email.SendEmail()
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return
	}
	slog.Info("Email sent successfully")
}
