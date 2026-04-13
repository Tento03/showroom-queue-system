package config

import (
	"fmt"
	"log"

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
}
