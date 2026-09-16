package dto

type CreateQueueRequest struct {
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
