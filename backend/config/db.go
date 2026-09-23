package config

import (
	"fmt"
	"log"

	"backend-queue/models"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		GetEnv("DB_USER"),
		GetEnv("DB_PASS"),
		GetEnv("DB_HOST"),
		GetEnv("DB_PORT"),
		GetEnv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB not connected")
	}

	DB = db
	log.Println("DB connected")

	// 1. AutoMigrate Service terlebih dahulu
	if err := DB.AutoMigrate(&models.Service{}); err != nil {
		log.Fatal("AutoMigrate Service failed:", err)
	}

	// 2. Pastikan ada minimal 1 service default untuk data antrean existing
	var defaultService models.Service
	if err := DB.First(&defaultService).Error; err != nil {
		defaultService = models.Service{
			ID:               uuid.New().String(),
			Name:             "Servis Umum",
			EstimatedMinutes: 30,
		}
		if err := DB.Create(&defaultService).Error; err != nil {
			log.Println("Failed to create default service:", err)
		}
	}

	// 3. Jika kolom service_id sudah ada di queues, update baris lama yang service_id-nya kosong/invalid
	var hasColumn int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'queues' AND column_name = 'service_id'").Scan(&hasColumn)
	if hasColumn > 0 && defaultService.ID != "" {
		DB.Exec("UPDATE queues SET service_id = ? WHERE service_id = '' OR service_id IS NULL OR service_id NOT IN (SELECT id FROM services)", defaultService.ID)
	}

	// 4. AutoMigrate Queue (sekarang aman karena tidak ada data yatim/orphan FK)
	if err := DB.AutoMigrate(&models.Queue{}); err != nil {
		log.Fatal("AutoMigrate Queue failed:", err)
	}
	log.Println("AutoMigrate success")

	var count int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'queues' AND index_name = 'idx_queue_number_date'").Scan(&count)
	if count == 0 {
		DB.Exec("CREATE UNIQUE INDEX idx_queue_number_date ON queues(queue_number, queue_date)")
	}
}
