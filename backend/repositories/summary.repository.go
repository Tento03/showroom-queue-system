package repositories

import (
	"backend-queue/config"
	"backend-queue/models"
	"fmt"
)

type SummaryRepository struct{}

func NewSummaryRepository() *SummaryRepository {
	return &SummaryRepository{}
}

// GetServiceBreakdown returns count of queues per service name for a given date
func (r *SummaryRepository) GetServiceBreakdown(date string) ([]struct {
	Name  string
	Count int
}, error) {
	var results []struct {
		Name  string
		Count int
	}

	err := config.DB.Model(&models.Queue{}).
		Select("services.name as name, count(*) as count").
		Joins("JOIN services ON queues.service_id = services.id").
		Where("queues.queue_date = ?", date).
		Group("services.id, services.name").
		Order("count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("GetServiceBreakdown: %w", err)
	}
	return results, nil
}

// GetPeakHour returns the hour (0-23) with the most queues and the count
func (r *SummaryRepository) GetPeakHour(date string) (hour int, count int, err error) {
	var result struct {
		Hour  int
		Count int
	}

	err = config.DB.Model(&models.Queue{}).
		Select("HOUR(created_at) as hour, count(*) as count").
		Where("queue_date = ?", date).
		Group("HOUR(created_at)").
		Order("count DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		return 0, 0, fmt.Errorf("GetPeakHour: %w", err)
	}
	return result.Hour, result.Count, nil
}
