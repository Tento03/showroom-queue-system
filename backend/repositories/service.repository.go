package repositories

import (
	"backend-queue/config"
	"backend-queue/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepository struct{}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{}
}

func (r *ServiceRepository) FindAll() ([]models.Service, error) {
	var services []models.Service
	result := config.DB.Order("created_at ASC").Find(&services)
	if result.Error != nil {
		return nil, fmt.Errorf("FindAll: %w", result.Error)
	}
	return services, nil
}

func (r *ServiceRepository) Create(service *models.Service) error {
	service.ID = uuid.New().String()
	result := config.DB.Create(service)
	if result.Error != nil {
		return fmt.Errorf("Create: %w", result.Error)
	}
	return nil
}

func (r *ServiceRepository) FindByID(id string) (*models.Service, error) {
	var service models.Service
	result := config.DB.Where("id = ?", id).First(&service)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("FindByID: %w", result.Error)
	}
	return &service, nil
}

func (r *ServiceRepository) UpdateEstimatedMinutes(id string, minutes int) error {
	result := config.DB.Model(&models.Service{}).Where("id = ?", id).Update("estimated_minutes", minutes)
	if result.Error != nil {
		return fmt.Errorf("UpdateEstimatedMinutes: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
