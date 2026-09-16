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

var allowedMIME = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

const maxFileSize = 5 << 20 // 5MB

func (ctrl *UploadController) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image required"})
		return
	}

	// Fix — validasi ukuran file
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 5MB"})
		return
	}

	// Fix — validasi MIME type asli
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	mimeType := http.DetectContentType(buffer)
	if !allowedMIME[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only JPEG, PNG, and WebP are allowed"})
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
