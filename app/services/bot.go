package services

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/InfaFS/NLQBot/app/llm"
	"github.com/InfaFS/NLQBot/app/messaging"
)

type BotService struct {
	llmClient      *llm.Client
	whatsappClient *messaging.WhatsappClient
	slackClient    *messaging.SlackClient
}

func NewBotService(llm *llm.Client, whatsapp *messaging.WhatsappClient, slack *messaging.SlackClient) *BotService {
	return &BotService{
		llmClient:      llm,
		whatsappClient: whatsapp,
		slackClient:    slack,
	}
}

// esta funcion luego tiene que abstraerse mas, por el momento voy a hacer toda la funcion aca
func (b *BotService) ProcessMessage(messanger string, recipient string, message string) error {

	type requestPayload struct {
		Messanger string `json:"messanger"`
		Recipient string `json:"recipient"`
		Message   string `json:"message"`
		SessionID string `json:"session_id"`
	}

	payload := requestPayload{
		Messanger: messanger,
		Recipient: recipient,
		Message:   message,
		SessionID: "prueba123",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:3000/input/bot", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil

}

func (b *BotService) SendWhatsappMessage(to string, message string) error {
	return b.whatsappClient.SendMessage(to, message)
}

func (b *BotService) SendSlackMessage(channel string, message string) error {
	return b.slackClient.SendMessage(channel, message)
}
