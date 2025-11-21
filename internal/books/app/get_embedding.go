package books_app

import (
	"context"

	"github.com/tariq-ventura/go-chatbot/internal/logs"
	"github.com/tariq-ventura/go-chatbot/internal/models"
)

func GetEmbedding(text string) ([]float32, error) {
	client, err := models.NewEmbedding(context.Background())

	if err != nil {
		logs.LogError("Failed to initialize embedding client", map[string]any{"error": err.Error()})
		return nil, err
	}

	embedding, err := client.GetEmbedding(text)

	if err != nil {
		logs.LogError("Failed to get embedding", map[string]any{"error": err.Error()})
		return nil, err
	}

	return embedding, nil
}
