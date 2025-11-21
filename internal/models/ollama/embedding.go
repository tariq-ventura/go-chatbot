package models_ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (oc *OllamaClient) GetEmbedding(text string) ([]float32, error) {
	ollamaUrl, exists := os.LookupEnv("OLLAMA_URL")

	if !exists {
		ollamaUrl = "http://localhost:11434"
	}

	reqBody := map[string]any{
		"model":  "qwen3-embedding:0.6b",
		"prompt": text,
	}

	jsonData, err := json.Marshal(reqBody)

	if err != nil {
		logs.LogError("Failed to marshal request body for Ollama embedding", map[string]any{"error": err.Error()})
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/embeddings", ollamaUrl),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		logs.LogError("Failed to make request to Ollama embedding API", map[string]any{"error": err.Error()})
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logs.LogError("Ollama embedding API returned non-200 status", map[string]any{"error": err.Error()})
		return nil, fmt.Errorf("ollama API error: %s", string(body))
	}

	var result OllamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logs.LogError("Failed to decode Ollama embedding response", map[string]any{"error": err.Error()})
		return nil, err
	}

	return result.Embedding, nil
}
