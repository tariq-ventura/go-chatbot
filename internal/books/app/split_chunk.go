package books_app

import "github.com/tariq-ventura/go-chatbot/internal/logs"

func SplitIntoChunks(text string, size int) []string {
	chunks := []string{}

	logs.LogInfo("Splitting text into chunks", map[string]any{"text_length": len(text), "chunk_size": size})

	for len(text) > size {
		chunk := text[:size]
		chunks = append(chunks, chunk)
		text = text[size:]
	}

	if len(text) > 0 {
		chunks = append(chunks, text)
	}

	return chunks
}
