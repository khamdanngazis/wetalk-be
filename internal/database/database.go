package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"chat-be/internal/config"
	"chat-be/internal/domain/entities"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_USER", "postgres"),
		config.GetEnv("DB_PASSWORD", "postgres"),
		config.GetEnv("DB_NAME", "chatdb"),
		config.GetEnv("DB_PORT", "5432"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Database connected")

	return DB
}

func InitMigration(db *gorm.DB) {
	db.AutoMigrate(&entities.User{}, &entities.Message{}, &entities.MessageStatus{}, &entities.ChatRoom{}, &entities.ChatRoomParticipant{}, &entities.SocketPath{})
}
