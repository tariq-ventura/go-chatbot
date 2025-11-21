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

	prompt := "Eres un asistente que SOLO responde basándose en el contexto proporcionado.\n\n" +
		"INSTRUCCIONES ESTRICTAS:\n" +
		"- Responde ÚNICAMENTE en español\n" +
		"- Si la respuesta NO está en el contexto, di: 'Lo siento, no tengo información sobre eso en el contexto proporcionado.'\n" +
		"- NO inventes ni agregues información que no esté en el contexto\n" +
		"- Sé conciso y directo\n" +
		"- NO incluyas referencias a páginas o títulos en tu respuesta, eso se agregará automáticamente\n\n" +
		"CONTEXTO:\n" +
		contextText +
		"\n\nPREGUNTA: " + question +
		"\n\nRESPUESTA:"

	reqBody := map[string]any{
		"model":  "thirdeyeai/Qwen2.5-1.5B-Instruct-uncensored:Q4_0",
		"prompt": prompt,
		"stream": false,
		"options": map[string]any{
			"temperature": 0.1,
			"top_p":       0.3,
		},
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
