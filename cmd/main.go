package main

import (
	"log"

	"github.com/InfaFS/NLQBot/app/handlers"
	"github.com/InfaFS/NLQBot/app/llm"
	"github.com/InfaFS/NLQBot/app/messaging"
	"github.com/InfaFS/NLQBot/app/services"
	"github.com/InfaFS/NLQBot/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	r := gin.Default()

	llmClient := llm.NewClient(cfg.GeminiAPIKey)
	whatsappClient := messaging.NewWhatsappClient(cfg.VerifyToken, cfg.WhatsappPhoneNumberID)
	slackClient := messaging.NewSlackClient(cfg.SlackToken)

	botService := services.NewBotService(llmClient, whatsappClient, slackClient)
	handler := &handlers.Handler{BotService: botService, VerifyToken: cfg.VerifyToken}

	r.GET("/webhook/whatsapp", handler.VerifyWhatsappWebhook)
	r.POST("/webhook/whatsapp", handler.HandleWhatsappWebhook)
	r.POST("/webhook/slack", handler.HandleSlackWebhook)
	r.POST("/webhook/bot", handler.ReceiveMessage)

	if err := r.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
