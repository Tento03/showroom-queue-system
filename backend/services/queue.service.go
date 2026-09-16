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
