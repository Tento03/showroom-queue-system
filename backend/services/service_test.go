package services

import (
	"backend-queue/dto"
	"backend-queue/models"
	"backend-queue/utils"
	"testing"
	"time"
)

func TestValidateStatusTransition(t *testing.T) {
	tests := []struct {
		current models.QueueStatus
		next    models.QueueStatus
		wantErr bool
	}{
		{models.StatusWaiting, models.StatusProcessing, false},
		{models.StatusWaiting, models.StatusCancelled, false},
		{models.StatusWaiting, models.StatusDone, true},
		{models.StatusProcessing, models.StatusDone, false},
		{models.StatusProcessing, models.StatusCancelled, false},
		{models.StatusProcessing, models.StatusWaiting, true},
		{models.StatusDone, models.StatusWaiting, true},
		{models.StatusCancelled, models.StatusProcessing, true},
	}

	for _, tt := range tests {
		err := validateStatusTransition(tt.current, tt.next)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateStatusTransition(%s, %s) error = %v, wantErr %v", tt.current, tt.next, err, tt.wantErr)
		}
	}
}

func TestActualMinutesCalculation(t *testing.T) {
	startedAt := time.Now().Add(-45 * time.Minute)
	completedAt := time.Now()

	minutes := int(completedAt.Sub(startedAt).Minutes())
	if minutes != 45 {
		t.Errorf("expected 45 minutes, got %d", minutes)
	}
}

func TestServiceValidation(t *testing.T) {
	svc := NewServiceService()
	_, err := svc.CreateService(&dto.CreateServiceRequest{
		Name:             "Ganti Oli",
		EstimatedMinutes: 0,
	})
	if err != utils.ErrInvalidEstimatedMinutes {
		t.Errorf("expected ErrInvalidEstimatedMinutes, got %v", err)
	}

	err = svc.UpdateEstimatedMinutes("some-id", -5)
	if err != utils.ErrInvalidEstimatedMinutes {
		t.Errorf("expected ErrInvalidEstimatedMinutes, got %v", err)
	}
}
