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

func (r *QueueRepository) GetQueuesByDate(date string) ([]models.Queue, error) {
	var queues []models.Queue
	result := config.DB.Where("queue_date = ?", date).
		Order("created_at ASC").
		Find(&queues)
	if result.Error != nil {
		return nil, fmt.Errorf("GetQueuesByDate: %w", result.Error)
	}
	return queues, nil
}

func (r *QueueRepository) CreateQueueTx(tx *gorm.DB, queue *models.Queue) error {
	queue.ID = uuid.New().String()
	result := tx.Create(queue)
	if result.Error != nil {
		return fmt.Errorf("CreateQueue: %w", result.Error)
	}
	return nil
}
