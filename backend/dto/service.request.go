package dto

type CreateServiceRequest struct {
	Name             string `json:"name" binding:"required"`
	EstimatedMinutes int    `json:"estimated_minutes" binding:"required,min=1"`
}

type UpdateServiceRequest struct {
	EstimatedMinutes int `json:"estimated_minutes" binding:"required,min=1"`
}
