package models

import "time"

type Service struct {
	ID               string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name             string    `gorm:"type:varchar(100);not null"  json:"name"`
	EstimatedMinutes int       `gorm:"not null"                   json:"estimated_minutes"`
	CreatedAt        time.Time `gorm:"autoCreateTime"             json:"created_at"`
}

func (Service) TableName() string { return "services" }
