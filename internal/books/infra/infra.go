package books_infra

import (
	"github.com/gin-gonic/gin"
	books_domain "github.com/tariq-ventura/go-chatbot/internal/books/domain"
)

type ChunkHandler struct{}

func NewChunkHandler(server *gin.Context) books_domain.ChunkInterface {
	return &ChunkHandler{}
}
