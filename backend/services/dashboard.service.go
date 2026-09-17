package services

import (
	"backend-queue/config"
	"backend-queue/dto"
	"backend-queue/repositories"
	"context"
	"encoding/json"
	"time"
)

type DashboardService struct {
	repo *repositories.QueueRepository
}

func NewDashboardService() *DashboardService {
	return &DashboardService{
		repo: repositories.NewQueueRepository(),
	}
}

func (s *DashboardService) GetDashboardStats(date string) (*dto.DashboardStats, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	ctx := context.Background()
	cacheKey := "dashboard:stats:" + date

	// Coba ambil dari Redis dulu
	cached, err := config.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var stats dto.DashboardStats
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			return &stats, nil
		}
	}

	// Cache miss — ambil dari DB
	raw, err := s.repo.GetStatsByDate(date)
	if err != nil {
		return nil, err
	}

	stats := &dto.DashboardStats{
		Total:     raw["total"],
		Waiting:   raw["waiting"],
		Processed: raw["processing"],
		Done:      raw["done"],
		Cancelled: raw["cancelled"],
	}

	// Simpan ke Redis — expire 30 detik
	// Stats hari ini: cache pendek biar tetap fresh
	// Stats hari lalu: cache lebih lama (data sudah tidak berubah)
	ttl := 30 * time.Second
	if date != time.Now().Format("2006-01-02") {
		ttl = 24 * time.Hour
	}

	data, _ := json.Marshal(stats)
	config.RDB.Set(ctx, cacheKey, data, ttl)

	return stats, nil
}

// Invalidate cache — dipanggil tiap ada queue baru atau status update
func InvalidateDashboardCache(date string) {
	ctx := context.Background()
	cacheKey := "dashboard:stats:" + date
	config.RDB.Del(ctx, cacheKey)
}
