package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ReceiveMessage(c *gin.Context) {

	type requestPayload struct {
		Messanger string `json:"messanger"`
		Recipient string `json:"recipient"`
		Message   string `json:"message"`
	}

	var body requestPayload
	if err := c.BindJSON(&body); err != nil {
		log.Println("Error al parsear JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	c.Status(http.StatusOK)

	if body.Messanger == "whatsapp" {
		h.BotService.SendWhatsappMessage(body.Recipient, body.Message)
	} else {
		h.BotService.SendSlackMessage(body.Recipient, body.Message)
	}

}
