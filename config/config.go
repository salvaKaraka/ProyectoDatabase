package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	VerifyToken           string
	WhatsappPhoneNumberID string
	GeminiAPIKey          string
	SlackToken            string
	Port                  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: no .env file found, falling back to env variables")
	}

	return &Config{
		VerifyToken:           os.Getenv("VERIFY_TOKEN"),
		WhatsappPhoneNumberID: os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
		GeminiAPIKey:          os.Getenv("GEMINI_API_KEY"),
		SlackToken:            os.Getenv("SLACK_TOKEN"),
		Port:                  os.Getenv("PORT"),
	}
}
