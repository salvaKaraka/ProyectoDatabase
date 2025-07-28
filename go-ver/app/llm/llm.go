package llm

// internal/llm/client.go

import (
	"context"
	"log"

	"google.golang.org/genai"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) GenerateContent(ctx context.Context, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)

	return result.Text(), nil

}
