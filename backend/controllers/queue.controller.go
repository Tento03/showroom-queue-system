package controllers

import (
	"backend-queue/dto"
	"backend-queue/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueueController struct {
	services *services.QueueService
}

func NewQueueController() *QueueController {
	return &QueueController{
		services: services.NewQueueService(),
	}
}

func (ctrl *QueueController) CreateQueue(c *gin.Context) {
	var req dto.CreateQueueRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queueNumber, err := ctrl.services.CreateQueue(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create queue"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"queue_number": queueNumber})
}
