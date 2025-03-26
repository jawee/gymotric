package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func SendEmail() error {
	SGKEY := os.Getenv("SENDGRID_KEY")
	client := &http.Client{}
	emailRequestBodyObj := &sendGridRequest{
		Personalizations: []personalization{
			{
				To: []from{
					{
						// Email: "jawee.dev@gmail.com",
						Email: "ld-ee3ef3a3af@dmarctester.com",
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
				Type:  "text/plain",
				Value: "Test",
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
