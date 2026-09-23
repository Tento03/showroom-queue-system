package controllers

import (
	"backend-queue/dto"
	"backend-queue/services"
	"backend-queue/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {
	service *services.ServiceService
}

func NewServiceController() *ServiceController {
	return &ServiceController{
		service: services.NewServiceService(),
	}
}

func (ctrl *ServiceController) GetServices(c *gin.Context) {
	result, err := ctrl.service.GetAllServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get services"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ctrl *ServiceController) CreateService(c *gin.Context) {
	var req dto.CreateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := ctrl.service.CreateService(&req)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidEstimatedMinutes) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func (ctrl *ServiceController) UpdateEstimatedMinutes(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.service.UpdateEstimatedMinutes(id, req.EstimatedMinutes)
	if err != nil {
		if errors.Is(err, utils.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, utils.ErrInvalidEstimatedMinutes) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "service updated"})
}
