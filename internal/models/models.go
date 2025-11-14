package models

import (
	"context"
	"errors"
	"os"

	"github.com/tariq-ventura/go-chatbot/internal/logs"
	models_ollama "github.com/tariq-ventura/go-chatbot/internal/models/ollama"
	models_openai "github.com/tariq-ventura/go-chatbot/internal/models/openai"
)

type Embedding interface {
	GetEmbedding(text string) ([]float32, error)
	AskQuestion(question string, contextTexts []string) (string, error)
}

var NewEmbedding = func(ctx context.Context) (Embedding, error) {
	et := os.Getenv("EMBEDDING_TYPE")

	switch et {
	case "openai":
		logs.LogInfo("Using OpenAI Embedding Service", nil)
		return models_openai.SetupOpenai(ctx), nil
	case "ollama":
		logs.LogInfo("Using Ollama Embedding Service", nil)
		return models_ollama.SetupOllama(ctx), nil
	default:
		return nil, errors.New("unsupported EMBEDDING_TYPE: " + et)
	}
}
