package routes

import (
	"backend-queue/controllers"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine) {
	queueController := controllers.NewQueueController()
	uploadController := controllers.NewUploadController()

	r.Static("/uploads", "./uploads")

	api := r.Group("/")
	{
		api.POST("/upload", uploadController.UploadImage)
		api.GET("/queues", queueController.GetQueues)
		api.POST("/queue", queueController.CreateQueue)
	}
}
