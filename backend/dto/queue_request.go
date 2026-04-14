package dto

type CreateQueueRequest struct {
	VehiclePlate    string `json:"vehicle_plate" binding:"required"`
	VehicleImageURL string `json:"vehicle_image_url" binding:"required"`
	OwnerName       string `json:"owner_name"`
	OwnerPhone      string `json:"owner_phone"`
}
