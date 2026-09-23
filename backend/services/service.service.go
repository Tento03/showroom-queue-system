package services

import (
	"backend-queue/dto"
	"backend-queue/models"
	"backend-queue/repositories"
	"backend-queue/utils"
	"errors"

	"gorm.io/gorm"
)

type ServiceService struct {
	repo *repositories.ServiceRepository
}

func NewServiceService() *ServiceService {
	return &ServiceService{
		repo: repositories.NewServiceRepository(),
	}
}

func (s *ServiceService) GetAllServices() ([]models.Service, error) {
	return s.repo.FindAll()
}

func (s *ServiceService) CreateService(req *dto.CreateServiceRequest) (*models.Service, error) {
	if req.EstimatedMinutes <= 0 {
		return nil, utils.ErrInvalidEstimatedMinutes
	}

	service := &models.Service{
		Name:             req.Name,
		EstimatedMinutes: req.EstimatedMinutes,
	}

	if err := s.repo.Create(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceService) UpdateEstimatedMinutes(id string, minutes int) error {
	if minutes <= 0 {
		return utils.ErrInvalidEstimatedMinutes
	}

	err := s.repo.UpdateEstimatedMinutes(id, minutes)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrServiceNotFound
		}
		return err
	}

	return nil
}
