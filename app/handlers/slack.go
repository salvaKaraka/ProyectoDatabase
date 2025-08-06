package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleSlackWebhook(c *gin.Context) {
	//json que viene de slack
	var data map[string]interface{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	//Identifica si quiere verificar o quiere mandar un mensaje

	//verificar
	if data["type"] == "url_verification" {
		challenge := data["challenge"].(string)
		c.JSON(http.StatusOK, gin.H{"challenge": challenge})
		return
	}

	//mensaje que llegó
	if data["type"] == "event_callback" {
		//usamos una go func porque tenemos que devolver rapido el status OK al server (sino manda de vuelta mensajes)
		go func() {
			event, ok := data["event"].(map[string]interface{})
			if !ok {
				fmt.Println("event is missing or not a map")
				return
			}

			// Asegurate que sea mensaje y no venga de un bot
			eventType, _ := event["type"].(string)
			_, hasBotID := event["bot_id"]
			text, hasText := event["text"].(string)
			channel, hasChannel := event["channel"].(string)

			if eventType == "message" && !hasBotID && hasText && hasChannel {
				response, err := h.BotService.ProcessMessage(text)
				if err != nil {
					response = "No se pudo realizar tu consulta, lo sentimos"
				}

				if err := h.BotService.SendSlackMessage(channel, response); err != nil {
					fmt.Println("error al enviar mensaje:", err)
				}
			}
		}()

	}

	c.Status(http.StatusOK)
}
