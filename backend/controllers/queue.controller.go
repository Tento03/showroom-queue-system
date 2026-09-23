package controllers

import (
	"backend-queue/dto"
	"backend-queue/services"
	"backend-queue/utils"
	"errors"
	"math"
	"net/http"
	"strconv"

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

// Ganti method ini
func (ctrl *QueueController) GetQueues(c *gin.Context) {
	date := c.Query("date")

	// Parse pagination params
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			page = val
		}
	}
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	queues, total, err := ctrl.services.GetQueues(date, page, limit)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidDateFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get queues"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"date":        date,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"queues":      queues,
	})
}

func (ctrl *QueueController) CreateQueue(c *gin.Context) {
	var req dto.CreateQueueRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queueNumber, err := ctrl.services.CreateQueue(&req)
	if err != nil {
		if errors.Is(err, utils.ErrServiceNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create queue"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"queue_number": queueNumber})
}

func (ctrl *QueueController) GetQueueByID(c *gin.Context) {
	id := c.Param("id")

	queue, err := ctrl.services.GetQueueByID(id)
	if err != nil {
		if errors.Is(err, utils.ErrQueueNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get queue"})
		return
	}

	c.JSON(http.StatusOK, queue)
}

func (ctrl *QueueController) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.services.UpdateStatus(id, req.Status)
	if err != nil {
		if errors.Is(err, utils.ErrQueueNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrInvalidStatusTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (ctrl *QueueController) DeleteQueue(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.services.DeleteQueue(id)
	if err != nil {
		if errors.Is(err, utils.ErrQueueNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrCannotDelete) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete queue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "queue deleted"})
}

// Ganti GetDashboardStats di controller
func (ctrl *QueueController) GetDashboardStats(c *gin.Context) {
	date := c.Query("date")

	svc := services.NewDashboardService()
	stats, err := svc.GetDashboardStats(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (ctrl *QueueController) GetEstimates(c *gin.Context) {
	date := c.Query("date")

	res, err := ctrl.services.GetEstimates(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get estimates"})
		return
	}

	c.JSON(http.StatusOK, res)
}

