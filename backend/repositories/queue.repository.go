package repositories

import (
	"backend-queue/config"
	"backend-queue/models"
	"fmt"

	"github.com/google/uuid"
)

type QueueRepository struct{}

func NewQueueRepository() *QueueRepository {
	return &QueueRepository{}
}

func (r *QueueRepository) CountTodayQueues() (int, error) {
	var count int64
	result := config.DB.Model(&models.Queue{}).
		Where("queue_date = CURDATE()").
		Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("CountTodayQueues: %w", result.Error)
	}
	return int(count), nil
}

func (r *QueueRepository) CreateQueue(queue *models.Queue) error {
	queue.ID = uuid.New().String()

	result := config.DB.Create(queue)
	if result.Error != nil {
		return fmt.Errorf("CreateQueue: %w", result.Error)
	}
	return nil
}
