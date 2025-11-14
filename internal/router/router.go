package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	Routes  *gin.Engine
	Context *gin.Context
}

func (r *Routes) SetupRouter() *gin.Engine {
	r.Routes = gin.Default()

	r.SetupCors()
	r.Routes.GET("/healthz", r.Print)

	r.BooksRoutes(r.Routes)

	return r.Routes
}

func (r *Routes) Run() {
	r.Routes.Run(":3000")
}

func (r *Routes) Print(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to my API with Golang"})
}
