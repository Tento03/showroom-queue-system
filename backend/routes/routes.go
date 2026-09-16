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
		api.GET("/queue/:id", queueController.GetQueueByID)
		api.PATCH("/queue/:id/status", queueController.UpdateStatus)
		api.DELETE("/queue/:id", queueController.DeleteQueue)

		// Dashboard
		api.GET("/dashboard/stats", queueController.GetDashboardStats)
	}
}
