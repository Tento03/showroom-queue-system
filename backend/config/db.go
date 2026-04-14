package config

import (
	"fmt"
	"log"

	"backend-queue/models"

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

	// AutoMigrate: buat/update tabel sesuai struct model
	if err := DB.AutoMigrate(&models.Queue{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}
	log.Println("AutoMigrate success")

	// Composite unique index: queue_number + queue_date
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_queue_number_date ON queues(queue_number, queue_date)")
}
