package dto

type CreateQueueRequest struct {
	ServiceID       string `json:"service_id" binding:"required"`
	VehiclePlate    string `json:"vehicle_plate" binding:"required"`
	VehicleImageURL string `json:"vehicle_image_url" binding:"required"`
	OwnerName       string `json:"owner_name"`
	OwnerPhone      string `json:"owner_phone"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=waiting processing done cancelled"`
}

type DashboardStats struct {
	Total     int `json:"total"`
	Waiting   int `json:"waiting"`
	Processed int `json:"processed"`
	Done      int `json:"done"`
	Cancelled int `json:"cancelled"`
}

type PaginationQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type PaginatedQueues struct {
	Date       string        `json:"date"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
	Queues     []interface{} `json:"queues"`
}

type QueueETA struct {
	ID                   string `json:"id"`
	QueueNumber          string `json:"queue_number"`
	VehiclePlate         string `json:"vehicle_plate"`
	Status               string `json:"status"`
	EstimatedWaitMinutes int    `json:"estimated_wait_minutes"`
	EstimatedDoneAt      string `json:"estimated_done_at"` // RFC3339 string
	QueuesAhead          int    `json:"queues_ahead"`
}

type ETAResponse struct {
	Date              string     `json:"date"`
	AvgServiceMinutes int        `json:"avg_service_minutes"`
	Estimates         []QueueETA `json:"estimates"`
}

