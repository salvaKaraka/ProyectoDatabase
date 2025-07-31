// internal/messaging/slack.go
package messaging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackClient struct {
	WebhookURL string
}

func NewSlackClient(webhookURL string) *SlackClient {
	return &SlackClient{WebhookURL: webhookURL}
}

func (s *SlackClient) SendMessage(channel, message string) error {
	payload := map[string]any{
		"text":    message,
		"channel": channel,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error enviando mensaje slack, status: %d", resp.StatusCode)
	}
	return nil
}
