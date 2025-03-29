package main

import (
	"bytes"
	"embed"
	"log/slog"
	"os"
	"text/template"
	"weight-tracker/internal/email"
)

var embedMigrations embed.FS
func main() {
	html, err := os.ReadFile("emails/reset-password.html")
	if err != nil {
		slog.Error("Failed to read HTML file", "error", err)
		return
	}

	tmpl, err := template.New("reset-password").Parse(string(html))
	if err != nil {
		slog.Error("Failed to parse HTML template", "error", err)
		return
	}

	data := ResetPasswordEmailData{
		Name:      "John Doe",
		ResetLink: "https://gymotric.com/reset-password?token=abcd1234",
	}

	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		return
	}

	err = email.SendEmail(emailContent.String())
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return
	}
	slog.Info("Email sent successfully")
}

type ResetPasswordEmailData struct {
	Name      string
	ResetLink string
}
