package router

import (
	"github.com/gin-gonic/gin"
	books_infra "github.com/tariq-ventura/go-chatbot/internal/books/infra"
)

func (ro *Routes) BooksRoutes(r *gin.Engine) {
	br := books_infra.NewChunkHandler(ro.Context)

	routes := r.Group("/api/v1/books")
	{
		routes.POST("/ingest", br.Insert)
	}
}
