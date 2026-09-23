package models

import "time"

type QueueStatus string

const (
	StatusWaiting    QueueStatus = "waiting"
	StatusProcessing QueueStatus = "processing"
	StatusDone       QueueStatus = "done"
	StatusCancelled  QueueStatus = "cancelled"
)

type Queue struct {
	ID              string      `gorm:"type:varchar(36);primaryKey"                json:"id"`
	QueueNumber     string      `gorm:"type:varchar(10)"                           json:"queue_number"`
	QueueDate       string      `gorm:"type:date;index"                            json:"queue_date"`
	VehiclePlate    string      `gorm:"type:varchar(20);not null"                  json:"vehicle_plate"`
	VehicleImageURL string      `gorm:"type:text;not null"                         json:"vehicle_image_url"`
	OwnerName       string      `gorm:"type:varchar(100)"                          json:"owner_name"`
	OwnerPhone      string      `gorm:"type:varchar(20)"                           json:"owner_phone"`
	ServiceID       string      `gorm:"type:varchar(36);index"                     json:"service_id"`
	Service         *Service    `gorm:"foreignKey:ServiceID"                       json:"service,omitempty"`
	Status          QueueStatus `gorm:"type:varchar(20);default:waiting"           json:"status"`
	StartedAt       *time.Time  `gorm:"column:started_at"                          json:"started_at,omitempty"`
	CompletedAt     *time.Time  `gorm:"column:completed_at"                        json:"completed_at,omitempty"`
	ActualMinutes   *int        `gorm:"column:actual_minutes"                      json:"actual_minutes,omitempty"`
	CreatedAt       time.Time   `gorm:"autoCreateTime"                             json:"created_at"`
}

// Composite unique: queue_number A001 boleh ada lagi, asal beda tanggal
func (Queue) TableName() string { return "queues" }
