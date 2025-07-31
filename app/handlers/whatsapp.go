package handlers

import (
	"log"
	"net/http"

	"github.com/InfaFS/NLQBot/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) VerifyWhatsappWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == h.VerifyToken {
		c.String(http.StatusOK, challenge)
		return
	}
	c.String(http.StatusForbidden, "Verification failed")
}

func (h *Handler) HandleWhatsappWebhook(c *gin.Context) {
	var body map[string]interface{}

	if err := c.BindJSON(&body); err != nil {
		log.Println("Error al parsear JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	log.Println(body)

	// Responder rápido para evitar reintentos
	c.Status(http.StatusOK)

	go func() {
		entries, ok := body["entry"].([]interface{})
		if !ok || len(entries) == 0 {
			log.Println("No se encontró 'entry' o está vacío")
			return
		}

		entry, ok := entries[0].(map[string]interface{})
		if !ok {
			log.Println("Formato de 'entry' inválido")
			return
		}

		changes, ok := entry["changes"].([]interface{})
		if !ok || len(changes) == 0 {
			log.Println("No se encontró 'changes' o está vacío")
			return
		}

		change, ok := changes[0].(map[string]interface{})
		if !ok {
			log.Println("Formato de 'change' inválido")
			return
		}

		val, ok := change["value"].(map[string]interface{})
		if !ok {
			log.Println("Formato de 'value' inválido")
			return
		}

		msgRaw, exists := val["messages"]
		if !exists {
			log.Println("No hay 'messages'")
			return
		}

		messages, ok := msgRaw.([]interface{})
		if !ok || len(messages) == 0 {
			log.Println("'messages' vacío o inválido")
			return
		}

		message, ok := messages[0].(map[string]interface{})
		if !ok {
			log.Println("Formato de mensaje inválido")
			return
		}

		from, _ := message["from"].(string)
		textObj, _ := message["text"].(map[string]interface{})
		text, _ := textObj["body"].(string)

		from = utils.ParseNumber(from)

		//despues hay que volver a usar esto bien
		response, err := h.BotService.ProcessMessage(text)
		if err != nil {
			response = "No se pudo realizar tu consulta, lo sentimos"
		}
		h.BotService.SendWhatsappMessage(from, response)
	}()
}
