package services

import (
	"backend-queue/dto"
	"backend-queue/models"
	"backend-queue/repositories"
	"fmt"
	"time"
)

type QueueService struct {
	repo *repositories.QueueRepository
}

func NewQueueService() *QueueService {
	return &QueueService{
		repo: repositories.NewQueueRepository(),
	}
}

func (s *QueueService) CreateQueue(req *dto.CreateQueueRequest) (string, error) {
	count, err := s.repo.CountTodayQueues()
	if err != nil {
		return "", err
	}

	queueNumber := fmt.Sprintf("D%03d", count+1)
	today := time.Now().Format("2006-01-02")

	queue := &models.Queue{
		QueueNumber:     queueNumber,
		QueueDate:       today,
		VehiclePlate:    req.VehiclePlate,
		VehicleImageURL: req.VehicleImageURL,
		OwnerName:       req.OwnerName,
		OwnerPhone:      req.OwnerPhone,
	}

	if err := s.repo.CreateQueue(queue); err != nil {
		return "", err
	}

	return queueNumber, nil
}
