package routes

import (
	"github.com/gin-gonic/gin"
	"hinze.dev/home/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	v1.GET("/email", controllers.GetEmail)
	return r
}
