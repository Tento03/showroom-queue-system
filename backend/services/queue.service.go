package services

import (
	"backend-queue/config"
	"backend-queue/dto"
	"backend-queue/models"
	"backend-queue/repositories"
	"backend-queue/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type QueueService struct {
	repo *repositories.QueueRepository
}

func NewQueueService() *QueueService {
	return &QueueService{
		repo: repositories.NewQueueRepository(),
	}
}

func (s *QueueService) GetQueues(date string) ([]models.Queue, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			// Fix #5 — return typed error, bukan string
			return nil, utils.ErrInvalidDateFormat
		}
	}
	return s.repo.GetQueuesByDate(date)
}

func (s *QueueService) CreateQueue(req *dto.CreateQueueRequest) (string, error) {
	var queueNumber string

	// Fix #1 — wrap dalam DB transaction + lock
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Lock row agar request paralel antri di sini
		var dummy models.Queue
		tx.Raw("SELECT id FROM queues WHERE queue_date = ? LIMIT 1 FOR UPDATE",
			time.Now().Format("2006-01-02")).
			Scan(&dummy)

		count, err := s.repo.CountTodayQueues(tx)
		if err != nil {
			return err
		}

		queueNumber = fmt.Sprintf("D%03d", count+1)
		today := time.Now().Format("2006-01-02")

		queue := &models.Queue{
			QueueNumber:     queueNumber,
			QueueDate:       today,
			VehiclePlate:    req.VehiclePlate,
			VehicleImageURL: req.VehicleImageURL,
			OwnerName:       req.OwnerName,
			OwnerPhone:      req.OwnerPhone,
		}

		return s.repo.CreateQueueTx(tx, queue)
	})

	if err != nil {
		return "", err
	}

	return queueNumber, nil
}

func (s *QueueService) GetQueueByID(id string) (*models.Queue, error) {
	queue, err := s.repo.GetQueueByID(id)
	if err != nil {
		return nil, err
	}
	if queue == nil {
		return nil, utils.ErrQueueNotFound
	}
	return queue, nil
}

func (s *QueueService) UpdateStatus(id string, status string) error {
	queue, err := s.repo.GetQueueByID(id)
	if err != nil {
		return err
	}
	if queue == nil {
		return utils.ErrQueueNotFound
	}

	if err := validateStatusTransition(queue.Status, models.QueueStatus(status)); err != nil {
		return err
	}

	return s.repo.UpdateStatus(id, models.QueueStatus(status))
}

func (s *QueueService) DeleteQueue(id string) error {
	queue, err := s.repo.GetQueueByID(id)
	if err != nil {
		return err
	}
	if queue == nil {
		return utils.ErrQueueNotFound
	}

	if queue.Status != models.StatusWaiting && queue.Status != models.StatusCancelled {
		return utils.ErrCannotDelete
	}

	return s.repo.DeleteQueue(id)
}

func (s *QueueService) GetDashboardStats(date string) (*dto.DashboardStats, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	stats, err := s.repo.GetStatsByDate(date)
	if err != nil {
		return nil, err
	}

	return &dto.DashboardStats{
		Total:     stats["total"],
		Waiting:   stats["waiting"],
		Processed: stats["processing"],
		Done:      stats["done"],
		Cancelled: stats["cancelled"],
	}, nil
}

func validateStatusTransition(current, next models.QueueStatus) error {
	allowed := map[models.QueueStatus][]models.QueueStatus{
		models.StatusWaiting:    {models.StatusProcessing, models.StatusCancelled},
		models.StatusProcessing: {models.StatusDone, models.StatusCancelled},
		models.StatusDone:       {},
		models.StatusCancelled:  {},
	}

	for _, s := range allowed[current] {
		if s == next {
			return nil
		}
	}
	return utils.ErrInvalidStatusTransition
}
