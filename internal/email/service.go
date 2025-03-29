package email

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

//go:embed emails/*.html
var embedEmails embed.FS

func SendEmail() error {
	SGKEY := os.Getenv("SENDGRID_KEY")

	html, err := embedEmails.ReadFile("emails/reset-password.html")
	if err != nil {
		slog.Error("Failed to read HTML file", "error", err)
		return err
	}

	tmpl, err := template.New("reset-password").Parse(string(html))
	if err != nil {
		slog.Error("Failed to parse HTML template", "error", err)
		return err
	}

	data := ResetPasswordEmailData{
		Name:      "John Doe",
		ResetLink: "https://gymotric.com/reset-password?token=abcd1234",
	}

	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		return err
	}

	client := &http.Client{}
	emailRequestBodyObj := &sendGridRequest{
		Personalizations: []personalization{
			{
				To: []from{
					{
						Email: "jawee.dev@gmail.com",
						// Email: "ld-ee3ef3a3af@dmarctester.com",
					},
				},
			},
		},
		From: from{
			Email: "noreply@gymotric.anol.se",
		},
		Subject: "Test",
		Content: []content{
			{
				Type:  "text/html",
				Value: emailContent.String(),
			},
		},
	}

	emailRequestBody, err := json.Marshal(emailRequestBodyObj)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(emailRequestBody))
	req.Header.Set("Authorization", "Bearer " + SGKEY)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to send email. Status code: %d. Message: %s", resp.StatusCode, string(bodyBytes))
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
	Email string `json:"email"`
}

type personalization struct {
	To []from `json:"to"`
}

type ResetPasswordEmailData struct {
	Name      string
	ResetLink string
}
