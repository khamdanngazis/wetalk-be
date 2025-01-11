package usecase_test

import (
	"chat-be/internal/config"
	"chat-be/internal/database"
	"chat-be/internal/domain/repositories"
	"chat-be/internal/usecases"
	"chat-be/package/logging"
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	userRepo        repositories.UserRepository
	chatRoomRepo    repositories.ChatRoomRepository
	chatRoomUsecase usecases.ChatRoomUsecase
	ctx             context.Context
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

	chatRoomUsecase = usecases.NewChatRoomUsecase(chatRoomRepo, userRepo)
	requestID := uuid.New().String()
	ctx = context.WithValue(context.Background(), logging.RequestIDKey, requestID)
}
