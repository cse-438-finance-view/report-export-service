package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/burakmike/report-export-service/docs"
)

// @title Report Export Service API
// @version 1.0
// @description This is a report export service API
// @host localhost:4444
// @BasePath /api/v1

func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/hello", HelloWorld)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":4444")
}

// HelloWorld godoc
// @Summary Hello world endpoint
// @Description Returns a hello world message
// @Tags hello
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /hello [get]
func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
