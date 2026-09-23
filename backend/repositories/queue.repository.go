package repositories

import (
	"backend-queue/config"
	"backend-queue/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QueueRepository struct{}

func NewQueueRepository() *QueueRepository {
	return &QueueRepository{}
}

func (r *QueueRepository) CountTodayQueues(tx *gorm.DB) (int, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	result := tx.Model(&models.Queue{}).
		Where("queue_date = ?", today).
		Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("CountTodayQueues: %w", result.Error)
	}
	return int(count), nil
}

// Ganti method ini
func (r *QueueRepository) GetQueuesByDate(date string, page, limit int) ([]models.Queue, int64, error) {
	var queues []models.Queue
	var total int64

	offset := (page - 1) * limit

	// Hitung total dulu
	config.DB.Model(&models.Queue{}).
		Where("queue_date = ?", date).
		Count(&total)

	// Ambil data dengan pagination
	result := config.DB.Preload("Service").Where("queue_date = ?", date).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&queues)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("GetQueuesByDate: %w", result.Error)
	}

	return queues, total, nil
}

func (r *QueueRepository) CreateQueueTx(tx *gorm.DB, queue *models.Queue) error {
	queue.ID = uuid.New().String()
	result := tx.Create(queue)
	if result.Error != nil {
		return fmt.Errorf("CreateQueue: %w", result.Error)
	}
	return nil
}

func (r *QueueRepository) GetQueueByID(id string) (*models.Queue, error) {
	var queue models.Queue
	result := config.DB.Preload("Service").Where("id = ?", id).First(&queue)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("GetQueueByID: %w", result.Error)
	}
	return &queue, nil
}

func (r *QueueRepository) UpdateStatus(id string, updates map[string]interface{}) error {
	result := config.DB.Model(&models.Queue{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

func (r *QueueRepository) GetActiveQueuesByDate(date string) ([]models.Queue, error) {
	var queues []models.Queue
	err := config.DB.Preload("Service").Where("queue_date = ? AND status IN ?", date, []models.QueueStatus{models.StatusWaiting, models.StatusProcessing}).
		Order("created_at ASC").
		Find(&queues).Error
	return queues, err
}

func (r *QueueRepository) GetAvgServiceMinutes(date string) (float64, error) {
	var avgMinutes *float64
	err := config.DB.Model(&models.Queue{}).
		Select("AVG(actual_minutes)").
		Where("queue_date = ? AND status = ? AND actual_minutes IS NOT NULL", date, models.StatusDone).
		Scan(&avgMinutes).Error

	if err != nil || avgMinutes == nil {
		return 0, err
	}
	return *avgMinutes, nil
}

func (r *QueueRepository) DeleteQueue(id string) error {
	result := config.DB.Where("id = ?", id).Delete(&models.Queue{})
	return result.Error
}

func (r *QueueRepository) GetStatsByDate(date string) (map[string]int, error) {
	var results []struct {
		Status string
		Count  int
	}

	err := config.DB.Model(&models.Queue{}).
		Select("status, count(*) as count").
		Where("queue_date = ?", date).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := map[string]int{
		"total":      0,
		"waiting":    0,
		"processing": 0,
		"done":       0,
		"cancelled":  0,
	}

	for _, r := range results {
		stats[r.Status] = r.Count
		stats["total"] += r.Count
	}

	return stats, nil
}
