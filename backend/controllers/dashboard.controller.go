package controllers

import (
	"backend-queue/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	summaryService *services.AISummaryService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		summaryService: services.NewAISummaryService(),
	}
}

// GetAISummary godoc
// GET /dashboard/summary?date=YYYY-MM-DD
// Returns an AI-generated daily summary for the given date.
// Defaults to today if no date is provided.
func (ctrl *DashboardController) GetAISummary(c *gin.Context) {
	date := c.Query("date")

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "format date tidak valid, gunakan YYYY-MM-DD",
		})
		return
	}

	result, err := ctrl.summaryService.GetDailySummary(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "gagal generate summary: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
