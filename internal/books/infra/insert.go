package books_infra

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	done := make(chan error)

	for i, ch := range chunks {
		go func(page int, content string) {
			err := database.InsertChunks(books_domain.Chunk{
				BookID:  book.ID,
				Page:    page,
				Content: content,
			})
			done <- err
		}(i+1, ch)
	}

	for i := 0; i < len(chunks); i++ {
		if err := <-done; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "ingested",
		"book_id":      book.ID,
		"title":        title,
		"total_chunks": len(chunks),
	})
}
