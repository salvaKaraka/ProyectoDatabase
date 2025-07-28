package services

import (
	"context"
	"log"

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

func (b *BotService) ProcessMessage(message string) (string, error) {
	// Aquí enviás el texto al LLM y recibís respuesta
	ctx := context.Background()
	response, err := b.llmClient.GenerateContent(ctx, message)
	if err != nil {
		log.Printf("Error generando respuesta LLM: %v", err)
		return "", err
	}
	return response, nil
}

func (b *BotService) SendWhatsappMessage(to string, message string) error {
	return b.whatsappClient.SendMessage(to, message)
}

func (b *BotService) SendSlackMessage(channel string, message string) error {
	return b.slackClient.SendMessage(channel, message)
}
