package services

import (
	"backend-queue/config"
	"backend-queue/dto"
	"backend-queue/models"
	"backend-queue/repositories"
	"backend-queue/utils"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

type QueueService struct {
	repo        *repositories.QueueRepository
	serviceRepo *repositories.ServiceRepository
}

func NewQueueService() *QueueService {
	return &QueueService{
		repo:        repositories.NewQueueRepository(),
		serviceRepo: repositories.NewServiceRepository(),
	}
}

func (s *QueueService) GetQueues(date string, page, limit int) ([]models.Queue, int64, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, 0, utils.ErrInvalidDateFormat
		}
	}

	// Default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.repo.GetQueuesByDate(date, page, limit)
}

func (s *QueueService) CreateQueue(req *dto.CreateQueueRequest) (string, error) {
	// Validasi service_id ada
	service, err := s.serviceRepo.FindByID(req.ServiceID)
	if err != nil {
		return "", err
	}
	if service == nil {
		return "", utils.ErrServiceNotFound
	}

	var queueNumber string

	err = config.DB.Transaction(func(tx *gorm.DB) error {
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
			ServiceID:       req.ServiceID,
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

	// 🆕 Invalidate cache setelah queue baru dibuat
	InvalidateDashboardCache(time.Now().Format("2006-01-02"))

	msg, _ := json.Marshal(map[string]string{
		"event":        "queue_created",
		"queue_number": queueNumber,
		"date":         time.Now().Format("2006-01-02"),
	})
	Hub.Broadcast(msg)

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

	nextStatus := models.QueueStatus(status)
	if err := validateStatusTransition(queue.Status, nextStatus); err != nil {
		return err
	}

	updates := map[string]interface{}{
		"status": nextStatus,
	}

	now := time.Now()
	if nextStatus == models.StatusProcessing {
		updates["started_at"] = &now
	} else if nextStatus == models.StatusDone {
		updates["completed_at"] = &now
		var startTime time.Time
		if queue.StartedAt != nil {
			startTime = *queue.StartedAt
		} else {
			startTime = queue.CreatedAt
		}
		minutes := int(now.Sub(startTime).Minutes())
		if minutes < 0 {
			minutes = 0
		}
		updates["actual_minutes"] = minutes
	}

	if err := s.repo.UpdateStatus(id, updates); err != nil {
		return err
	}

	// 🆕 Invalidate cache setelah status berubah
	InvalidateDashboardCache(queue.QueueDate)

	msg, _ := json.Marshal(map[string]string{
		"event":  "status_updated",
		"id":     id,
		"status": status,
	})
	Hub.Broadcast(msg)

	return nil
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

func (s *QueueService) GetEstimates(date string) (*dto.ETAResponse, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	avgMinutes, err := s.repo.GetAvgServiceMinutes(date)
	if err != nil || avgMinutes <= 0 {
		avgMinutes = 30.0 // Default baseline Service ETA = 30 mins
	}

	avgInt := int(avgMinutes)

	activeQueues, err := s.repo.GetActiveQueuesByDate(date)
	if err != nil {
		return nil, err
	}

	var estimates []dto.QueueETA
	now := time.Now()

	for i, q := range activeQueues {
		queuesAhead := i
		estWait := queuesAhead * avgInt
		estDoneAt := now.Add(time.Duration(estWait) * time.Minute)

		estimates = append(estimates, dto.QueueETA{
			ID:                   q.ID,
			QueueNumber:          q.QueueNumber,
			VehiclePlate:         q.VehiclePlate,
			Status:               string(q.Status),
			EstimatedWaitMinutes: estWait,
			EstimatedDoneAt:      estDoneAt.Format(time.RFC3339),
			QueuesAhead:          queuesAhead,
		})
	}

	return &dto.ETAResponse{
		Date:              date,
		AvgServiceMinutes: avgInt,
		Estimates:         estimates,
	}, nil
}

func (s *QueueService) GetQueueETA(id string) (*dto.SingleQueueETAResponse, error) {
	queue, err := s.repo.GetQueueByID(id)
	if err != nil {
		return nil, err
	}
	if queue == nil {
		return nil, utils.ErrQueueNotFound
	}

	serviceName := ""
	if queue.Service != nil {
		serviceName = queue.Service.Name
	}

	now := time.Now()

	// Jika antrean sudah selesai atau dibatalkan
	if queue.Status == models.StatusDone || queue.Status == models.StatusCancelled {
		return &dto.SingleQueueETAResponse{
			QueueNumber:          queue.QueueNumber,
			Service:              serviceName,
			EstimatedWaitMinutes: 0,
			EstimatedReadyAt:     now.Format(time.RFC3339),
			QueuesAhead:          0,
		}, nil
	}

	// 1. Ambil semua antrian dengan status "waiting" atau "processing" yang created_at lebih awal dari queue yang diminta
	queuesAhead, err := s.repo.GetQueuesAhead(queue.QueueDate, queue.CreatedAt)
	if err != nil {
		return nil, err
	}

	queuesAheadCount := len(queuesAhead)

	// Helper untuk mendapatkan durasi servis efektif (cek historis >= 5 data)
	durationCache := make(map[string]int)
	getEffectiveDuration := func(q *models.Queue) int {
		if q.Service == nil {
			return 30
		}
		if dur, ok := durationCache[q.ServiceID]; ok {
			return dur
		}

		histAvg, err := s.repo.GetHistoricalAvgByService(q.ServiceID)
		if err == nil && histAvg != nil && *histAvg > 0 {
			dur := int(math.Round(*histAvg))
			durationCache[q.ServiceID] = dur
			return dur
		}

		dur := q.Service.EstimatedMinutes
		if dur <= 0 {
			dur = 30
		}
		durationCache[q.ServiceID] = dur
		return dur
	}

	var totalSisaWaktu float64 = 0

	// Jika queue yang diminta sendiri sedang processing
	if queue.Status == models.StatusProcessing {
		dur := getEffectiveDuration(queue)
		elapsed := 0.0
		if queue.StartedAt != nil {
			elapsed = now.Sub(*queue.StartedAt).Minutes()
		}
		sisa := float64(dur) - elapsed
		if sisa < 0 {
			sisa = 0
		}
		totalSisaWaktu = sisa
	} else {
		// Queue berstatus waiting
		for _, q := range queuesAhead {
			dur := getEffectiveDuration(&q)
			if q.Status == models.StatusProcessing {
				// 2. Untuk yang "processing": hitung sisa waktu = service.estimated_minutes - (now() - started_at dalam menit)
				elapsed := 0.0
				if q.StartedAt != nil {
					elapsed = now.Sub(*q.StartedAt).Minutes()
				}
				sisa := float64(dur) - elapsed
				if sisa < 0 {
					sisa = 0
				}
				totalSisaWaktu += sisa
			} else if q.Status == models.StatusWaiting {
				// 3. Untuk yang "waiting": ambil service.estimated_minutes
				totalSisaWaktu += float64(dur)
			}
		}
	}

	// 4. Total ETA = jumlah semua sisa waktu di atas
	// 5. Tambahkan buffer 10% untuk variasi
	totalWithBuffer := totalSisaWaktu * 1.10
	estimatedWaitMinutes := int(math.Round(totalWithBuffer))
	if estimatedWaitMinutes < 0 {
		estimatedWaitMinutes = 0
	}

	estimatedReadyAt := now.Add(time.Duration(estimatedWaitMinutes) * time.Minute)

	return &dto.SingleQueueETAResponse{
		QueueNumber:          queue.QueueNumber,
		Service:              serviceName,
		EstimatedWaitMinutes: estimatedWaitMinutes,
		EstimatedReadyAt:     estimatedReadyAt.Format(time.RFC3339),
		QueuesAhead:          queuesAheadCount,
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
