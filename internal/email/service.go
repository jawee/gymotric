package email

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"weight-tracker/internal/utils"

	_ "github.com/joho/godotenv/autoload"
)

//go:embed emails/*.html
var embedEmails embed.FS

type ResetPasswordEmailData struct {
	Name      string
	ResetLink string
}

type SendEmailConfirmationData struct {
	Name string
	Link string
}

type SendAccountConfirmationData struct {
	Name string
	Link string
}

func SendPasswordReset(recipient string, data ResetPasswordEmailData) error {
	html, err := embedEmails.ReadFile("emails/reset-password.html")
	if err != nil {
		slog.Error("Failed to read HTML file", "error", err)
		return fmt.Errorf("Failed to read reset password HTML file: %w", err)
	}

	err = sendEmail(string(html), recipient, "Password Reset", data)
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return fmt.Errorf("Failed to send reset password email: %w", err)
	}

	return nil
}

func SendEmailConfirmation(recipient string, data SendEmailConfirmationData) error {
	html, err := embedEmails.ReadFile("emails/confirm-email.html")
	if err != nil {
		slog.Error("Failed to read HTML file", "error", err)
		return fmt.Errorf("Failed to read email confirmation HTML file: %w", err)
	}

	err = sendEmail(string(html), recipient, "Email Confirmation", data)
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return fmt.Errorf("Failed to send email confirmation email: %w", err)
	}

	return nil
}

func SendAccountConfirmation(recipient string, data SendAccountConfirmationData) error {
	html, err := embedEmails.ReadFile("emails/confirm-registration.html")
	if err != nil {
		slog.Error("Failed to read HTML file", "error", err)
		return fmt.Errorf("Failed to read account confirmation HTML file: %w", err)
	}

	err = sendEmail(string(html), recipient, "Confirm Account", data)
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return fmt.Errorf("Failed to send account confirmation email: %w", err)
	}

	return nil
}

func sendEmail(html string, recipient string, subject string, data any) error {
	tmpl, err := template.New("email").Parse(string(html))
	if err != nil {
		slog.Error("Failed to parse HTML template", "error", err)
		return fmt.Errorf("Failed to parse HTML template: %w", err)
	}

	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		return fmt.Errorf("Failed to execute template: %w", err)
	}

	SGKEY := os.Getenv(utils.EnvSendGridApiKey)
	if SGKEY != "" {
		res := sendSendGridEmail(emailContent, recipient, subject, SGKEY)
		if res != nil {
			return fmt.Errorf("Failed to send sendgrid email: %w", err)
		}
	}

	BREVOKEY := os.Getenv(utils.EnvBrevoApiKey)
	if BREVOKEY != "" {
		res := sendBrevoEmail(emailContent, recipient, subject, BREVOKEY)
		if res != nil {
			return fmt.Errorf("Failed to send brevo email: %w", err)
		}
	}
	return nil
}

func sendSendGridEmail(emailContent bytes.Buffer, recipient string, subject string, apiKey string) error {
	if apiKey == "" {
		slog.Error("SendGrid API key not set", "SGKEY", "")
		return errors.New("SendGrid API key not set")
	}

	client := &http.Client{}
	emailRequestBodyObj := &sendGridRequest{
		Personalizations: []personalization{
			{
				To: []from{
					{
						Email: recipient,
					},
				},
			},
		},
		From: from{
			Name:  "Gymotric",
			Email: "noreply@gymotric.anol.se",
		},
		Subject: subject + " - Gymotric",
		Content: []content{
			{
				Type:  "text/html",
				Value: emailContent.String(),
			},
		},
	}

	emailRequestBody, err := json.Marshal(emailRequestBodyObj)
	if err != nil {
		return fmt.Errorf("Failed to marshal sendgrid email request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(emailRequestBody))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send sendgrid email request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Failed to read sendgrid response body: %w", err)
		}
		return fmt.Errorf("Failed to send sendgrid email. Status code: %d. Message: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func sendBrevoEmail(emailContent bytes.Buffer, recipient string, subject string, apiKey string) error {
	if apiKey == "" {
		slog.Error("Brevo API key not set", "Brevo API key", "")
		return errors.New("Brevo API key not set")
	}

	client := &http.Client{}
	emailRequestBodyObj := &brevoRequest{
		Sender: brevoSender{
			Name:  "Gymotric",
			Email: "noreply@gymotric.anol.se",
		},
		To: []brevoTo{
			{
				Name: recipient,
				Email: recipient,
			},
		},
		Subject:     subject + " - Gymotric",
		HTMLContent: emailContent.String(),
	}

	emailRequestBody, err := json.Marshal(emailRequestBodyObj)
	if err != nil {
		return fmt.Errorf("Failed to marshal brevo email request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(emailRequestBody))
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send brevo email request", "error", err)
		return fmt.Errorf("Failed to send brevo email request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Failed to read brevo response body", "error", err)
			return fmt.Errorf("Failed to read brevo response body: %w", err)
		}
		slog.Error("Failed to send brevo email", "status_code", resp.StatusCode, "message", string(bodyBytes))
		return fmt.Errorf("Failed to send brevo email. Status code: %d. Message: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

type sendGridRequest struct {
	Personalizations []personalization `json:"personalizations"`
	From             from              `json:"from"`
	Subject          string            `json:"subject"`
	Content          []content         `json:"content"`
}

type content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type from struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type personalization struct {
	To []from `json:"to"`
}

type brevoRequest struct {
	Sender      brevoSender `json:"sender"`
	To          []brevoTo   `json:"to"`
	Subject     string      `json:"subject"`
	HTMLContent string      `json:"htmlContent"`
}
type brevoSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type brevoTo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
