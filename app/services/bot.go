package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
func (b *BotService) ProcessMessage(message string) (string, error) {
	type requestPayload struct {
		Question string `json:"question"`
	}

	payload := requestPayload{
		Question: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/query/meinlup/f1", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "qRkHi9L046ehpfoydFDmQCxz_KdCOf7aUGqhiUQHn7s")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Leemos el body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parseamos el JSON
	type responsePayload struct {
		Status      string   `json:"status"`
		Result      string   `json:"result"`      // suponemos que es string
		Explicacion string   `json:"explicacion"` // no lo usamos, pero lo dejamos por claridad
		Questions   []string `json:"questions"`
	}

	var parsedResponse responsePayload
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return "", err
	}

	// Validamos que la respuesta haya sido exitosa
	if (parsedResponse.Status != "success") && (parsedResponse.Status != "clarification") {
		return "", fmt.Errorf("error del servidor: %s", parsedResponse.Status)
	}

	//este parseo se realiza ya que del backend de python recibimos un array de strings
	if parsedResponse.Questions != nil {
		var res string
		for i := 0; i < len(parsedResponse.Questions); i++ {
			res += parsedResponse.Questions[i] + " "
		}

		return res, nil
	}
	// Retornamos solo el resultado
	return parsedResponse.Result, nil
}

func (b *BotService) SendWhatsappMessage(to string, message string) error {
	return b.whatsappClient.SendMessage(to, message)
}

func (b *BotService) SendSlackMessage(channel string, message string) error {
	return b.slackClient.SendMessage(channel, message)
}
