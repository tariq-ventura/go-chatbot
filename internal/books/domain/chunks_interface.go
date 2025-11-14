package books_domain

import "github.com/gin-gonic/gin"

type ChunkInterface interface {
	Insert(c *gin.Context)
}
