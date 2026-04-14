package services

import (
	"backend-queue/dto"
	"backend-queue/models"
	"backend-queue/repositories"
	"fmt"
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
	// Step 1: Count how many queues were created today
	count, err := s.repo.CountTodayQueues()
	if err != nil {
		return "", err
	}

	// Step 2: Generate queue number — resets every day, format A001..A999
	queueNumber := fmt.Sprintf("D%03d", count+1)

	// Step 3: Build the model with data from request
	queue := &models.Queue{
		QueueNumber:     queueNumber,
		VehiclePlate:    req.VehiclePlate,
		VehicleImageURL: req.VehicleImageURL,
		OwnerName:       req.OwnerName,
		OwnerPhone:      req.OwnerPhone,
	}

	// Step 4: Persist to database
	if err := s.repo.CreateQueue(queue); err != nil {
		return "", err
	}

	// Step 5: Return queue number to controller
	return queueNumber, nil
}
