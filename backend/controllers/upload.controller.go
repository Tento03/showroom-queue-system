package controllers

import (
	"backend-queue/config"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (ctrl *UploadController) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image required"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := fmt.Sprintf("uploads/%s", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	host := config.GetEnv("APP_HOST")
	port := config.GetEnv("APP_PORT")
	imageUrl := fmt.Sprintf("http://%s:%s/uploads/%s", host, port, filename)

	c.JSON(http.StatusOK, gin.H{"image_url": imageUrl})
}
