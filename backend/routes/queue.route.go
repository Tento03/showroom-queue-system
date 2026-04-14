package routes

import (
	"backend-queue/controllers"

	"github.com/gin-gonic/gin"
)

func QueueRoute(r *gin.Engine) {
	queueController := controllers.NewQueueController()

	queue := r.Group("/queue")
	{
		queue.POST("", queueController.CreateQueue)
	}
}
