package models

import "time"

type Queue struct {
	ID              string    `gorm:"type:varchar(20);primaryKey" json:"id"`
	QueueNumber     string    `gorm:"type:varchar(10);uniqueIndex" json:"queue_number"`
	VehiclePlate    string    `gorm:"type:varchar(20);not null" json:"vehicle_plate"`
	VehicleImageURL string    `gorm:"type:text;not null" json:"vehicle_image_url"`
	OwnerName       string    `gorm:"type:varchar(100)" json:"owner_name"`
	OwnerPhone      string    `gorm:"type:varchar(20)" json:"owner_phone"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}
