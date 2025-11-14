package books_infra

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	books_app "github.com/tariq-ventura/go-chatbot/internal/books/app"
	books_db "github.com/tariq-ventura/go-chatbot/internal/books/db"
	books_domain "github.com/tariq-ventura/go-chatbot/internal/books/domain"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

func (ch *ChunkHandler) Insert(c *gin.Context) {
	ctx := c.Request.Context()
	title := c.PostForm("title")
	author := c.PostForm("author")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required (field 'file'), accept plain txt or markdown"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot open file"})
		return
	}
	defer file.Close()

	var extracted string
	switch filepath.Ext(fileHeader.Filename) {
	case ".pdf":
		extracted, err = books_app.ExtractTextFromPDF(file)
	default:
		c.JSON(400, gin.H{"error": "unsupported file type"})
		return
	}

	database, err := books_db.NewDatabase(ctx)
	if err != nil {
		logs.LogError("Database connection error", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal setup error"})
		return
	}

	namespace := uuid.New()
	name := []byte("books")
	id := uuid.NewSHA1(namespace, name)

	book := books_domain.Book{ID: id.String(), Title: title, Author: author}
	if err := database.InsertBook(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	chunks := books_app.SplitIntoChunks(extracted, 1500)
	logs.LogInfo("Book split into chunks", map[string]any{"book_id": book.ID, "total_chunks": len(chunks)})

	done := make(chan error, len(chunks))
	semaphore := make(chan struct{}, 200)

	for i, ch := range chunks {
		semaphore <- struct{}{}
		go func(page int, content string) {
			defer func() { <-semaphore }()
			embedding, error := books_app.GetEmbedding(content)

			if error != nil {
				logs.LogError("Error getting embedding", map[string]any{"page": page, "error": error.Error()})
				done <- error
				return
			}

			err := database.InsertChunks(books_domain.Chunk{
				BookID:    book.ID,
				Page:      page,
				Content:   content,
				Embedding: pgvector.NewVector(embedding),
			})

			if err != nil {
				logs.LogError("Error inserting chunk", map[string]any{"page": page, "error": err.Error()})
			} else {
				logs.LogInfo("Chunk inserted successfully", map[string]any{"page": page, "book_id": book.ID})
			}

			done <- err
		}(i+1, ch)
	}

	var errors []error
	for i := 0; i < len(chunks); i++ {
		if err := <-done; err != nil {
			errors = append(errors, err)
		}
	}

	close(done)

	if len(errors) > 0 {
		logs.LogError("Some chunks failed to insert", map[string]any{"failed_count": len(errors), "total": len(chunks)})
		c.JSON(http.StatusPartialContent, gin.H{
			"status":       "partially_ingested",
			"book_id":      book.ID,
			"title":        title,
			"total_chunks": len(chunks),
			"failed":       len(errors),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "ingested",
		"book_id":      book.ID,
		"title":        title,
		"total_chunks": len(chunks),
	})
}
