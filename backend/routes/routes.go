package routes

import (
	"backend-queue/controllers"
	"backend-queue/middleware"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine) {
	queueController := controllers.NewQueueController()
	uploadController := controllers.NewUploadController()
	wsController := controllers.NewWebSocketController()
	serviceController := controllers.NewServiceController()

	r.Static("/uploads", "./uploads")

	api := r.Group("/")
	{
		api.POST("/upload", middleware.RateLimitUpload(), uploadController.UploadImage)

		// Services
		api.GET("/services", serviceController.GetServices)
		api.POST("/services", serviceController.CreateService)
		api.PATCH("/services/:id", serviceController.UpdateEstimatedMinutes)

		// Queue
		api.GET("/queues", queueController.GetQueues)
		api.POST("/queue", queueController.CreateQueue)
		api.GET("/queue/:id", queueController.GetQueueByID)
		api.PATCH("/queue/:id/status", queueController.UpdateStatus)
		api.DELETE("/queue/:id", queueController.DeleteQueue)

		// Dashboard
		api.GET("/dashboard/stats", queueController.GetDashboardStats)
		api.GET("/dashboard/estimates", queueController.GetEstimates)

		api.GET("/ws", wsController.Handle)
	}
}
