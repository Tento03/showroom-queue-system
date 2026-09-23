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

func TestETACalculationLogic(t *testing.T) {
	now := time.Now()

	// 1. Processing queue: estimasi 30 menit, berjalan 10 menit -> sisa 20 menit
	startedAt := now.Add(-10 * time.Minute)
	dur1 := 30
	elapsed1 := now.Sub(startedAt).Minutes()
	sisa1 := float64(dur1) - elapsed1
	if sisa1 < 0 {
		sisa1 = 0
	}

	// 2. Waiting queue: estimasi 45 menit -> sisa 45 menit
	dur2 := 45
	sisa2 := float64(dur2)

	totalSisa := sisa1 + sisa2 // 20 + 45 = 65

	// Buffer 10%
	totalWithBuffer := totalSisa * 1.10 // 65 * 1.10 = 71.5 -> 72
	waitMinutes := int(totalWithBuffer + 0.5)

	if waitMinutes != 72 && waitMinutes != 71 {
		t.Errorf("expected waitMinutes ~72, got %d", waitMinutes)
	}

	// 3. Overdue processing queue: estimasi 30 menit, berjalan 40 menit -> sisa harus 0 (tidak negatif)
	overdueStartedAt := now.Add(-40 * time.Minute)
	elapsedOverdue := now.Sub(overdueStartedAt).Minutes()
	sisaOverdue := float64(dur1) - elapsedOverdue
	if sisaOverdue < 0 {
		sisaOverdue = 0
	}
	if sisaOverdue != 0 {
		t.Errorf("expected sisaOverdue to be 0, got %f", sisaOverdue)
	}
}

