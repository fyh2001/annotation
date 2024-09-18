package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controllers"
	"server/app/middlewares/cors"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./static")

	r.Use(cors.Cors())

	root := r.Group("/api")
	{
		root.POST("/get-task", controllers.Task.GetTask)
		root.POST("/submit", controllers.Task.Submit)
	}

	return r
}
