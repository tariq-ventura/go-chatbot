package models_ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

func buildContext(chunks []string) string {
	return strings.Join(chunks, "\n---\n")
}

func (oc *OllamaClient) AskQuestion(question string, contextTexts []string) (string, error) {
	ollamaUrl, exists := os.LookupEnv("OLLAMA_URL")

	if !exists {
		ollamaUrl = "http://localhost:11434"
	}

	contextText := buildContext(contextTexts)

	prompt := "Usa estrictamente este contexto para responder:\n\n" +
		contextText +
		"\n\nPregunta: " + question +
		"\nRespuesta:"

	reqBody := map[string]any{
		"model":  "tinyllama:1.1b",
		"prompt": prompt,
		"stream": false,
	}

	jsonData, err := json.Marshal(reqBody)

	if err != nil {
		logs.LogError("Failed to marshal request body for Ollama", map[string]any{"error": err.Error()})
		return "", err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/generate", ollamaUrl),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		logs.LogError("Failed to make request to Ollama API", map[string]any{"error": err.Error()})
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logs.LogError("Ollama API returned error", map[string]any{
			"status": resp.StatusCode,
			"body":   string(bodyBytes),
		})
		return "", fmt.Errorf("ollama API error: status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logs.LogError("Failed to decode Ollama response", map[string]any{"error": err.Error()})
		return "", err
	}

	return result.Response, nil
}
