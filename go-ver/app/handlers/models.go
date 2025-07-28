package handlers

import "github.com/InfaFS/NLQBot/app/services"

type Handler struct {
	BotService  *services.BotService
	VerifyToken string
}
