package messaging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type WhatsappClient struct {
	AccessToken   string
	PhoneNumberID string
}

func NewWhatsappClient(token, phoneID string) *WhatsappClient {
	return &WhatsappClient{AccessToken: token, PhoneNumberID: phoneID}
}

func (w *WhatsappClient) SendMessage(to, message string) error {

	url := fmt.Sprintf("https://graph.facebook.com/v23.0/%s/messages", w.PhoneNumberID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text":              map[string]string{"body": message},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println(resp)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("error enviando mensaje whatsapp, status: %d", resp.StatusCode)
	}
	return nil
}
