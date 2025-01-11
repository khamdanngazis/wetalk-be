package repository_test

import (
	"chat-be/internal/config"
	"chat-be/internal/database"
	"chat-be/internal/domain/repositories"
	"chat-be/package/logging"
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	userRepo     repositories.UserRepository
	chatRoomRepo repositories.ChatRoomRepository
	ctx          context.Context
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {

	config.LoadEnv()

	// Initialize Database
	db := database.InitDB()

	chatRoomRepo = repositories.NewChatRoomRepository(db)
	userRepo = repositories.NewUserRepository(db)

	requestID := uuid.New().String()
	ctx = context.WithValue(context.Background(), logging.RequestIDKey, requestID)
}
