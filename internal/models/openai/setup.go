package models_openai

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

type OpenaiClient struct {
	client *openai.Client
	ctx    context.Context
}

var SetupOpenai = func(ctx context.Context) *OpenaiClient {
	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")

	if !exists {
		logs.LogError("OPENAI_API_KEY not set in environment variables", nil)
		panic("OPENAI_API_KEY not set in environment variables")
	}

	client := openai.NewClient(apiKey)

	return &OpenaiClient{
		client: client,
		ctx:    ctx,
	}

}
