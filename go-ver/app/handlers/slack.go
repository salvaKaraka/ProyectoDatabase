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
			event := data["event"].(map[string]interface{})
			if event["type"] == "message" && event["bot_id"] == nil {
				text := event["text"].(string)
				channel := event["channel"].(string)

				response := "Estoy procesando tu consulta..."
				err := h.BotService.ProcessMessage("slack", channel, text)
				if err != nil {
					response = "No se pudo realizar tu consulta"
				}

				err = h.BotService.SendSlackMessage(channel, response)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}

	c.Status(http.StatusOK)
}
