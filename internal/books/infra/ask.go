package books_infra

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	books_app "github.com/tariq-ventura/go-chatbot/internal/books/app"
	books_db "github.com/tariq-ventura/go-chatbot/internal/books/db"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
	"github.com/tariq-ventura/go-chatbot/internal/models"
)

func (ch *ChunkHandler) AskQuestion(c *gin.Context) {
	ctx := c.Request.Context()
	question := c.PostForm("question")

	logs.LogInfo("Received question", map[string]any{"question": question})

	emb, err := books_app.GetEmbedding(question)

	if err != nil {
		logs.LogError("Failed to get embedding", map[string]any{"error": err.Error()})
		c.JSON(500, gin.H{"error": "Failed to get embedding"})
		return
	}

	database, err := books_db.NewDatabase(ctx)
	if err != nil {
		logs.LogError("Database connection error", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal setup error"})
		return
	}

	chunks, err := database.SearchChunks(emb, 3)

	if err != nil {
		logs.LogError("Failed to search chunks", map[string]any{"error": err.Error()})
		c.JSON(500, gin.H{"error": "Failed to search chunks"})
		return
	}

	contextTexts := make([]string, len(chunks))
	for i, chunk := range chunks {
		contextTexts[i] = chunk.Content
	}

	client, err := models.NewEmbedding(context.Background())

	if err != nil {
		logs.LogError("Failed to initialize embedding client", map[string]any{"error": err.Error()})
		c.JSON(500, gin.H{"error": "Failed to initialize embedding client"})
		return
	}

	response, err := client.AskQuestion(question, contextTexts)

	if err != nil {
		logs.LogError("Failed to get answer", map[string]any{"error": err.Error()})
		c.JSON(500, gin.H{"error": "Failed to get answer"})
		return
	}

	var sources []map[string]interface{}
	for _, chunk := range chunks {
		sources = append(sources, map[string]interface{}{
			"title": chunk.Title,
			"page":  chunk.Page,
		})
	}

	var sourceText strings.Builder
	sourceText.WriteString("\n\n---\nFuentes:\n")
	for _, chunk := range chunks {
		sourceText.WriteString(fmt.Sprintf("- %s (Página %d)\n", chunk.Title, chunk.Page))
	}

	c.JSON(200, gin.H{
		"answer":      response + sourceText.String(),
		"chunks_used": len(chunks),
		"sources":     sources,
	})
}
