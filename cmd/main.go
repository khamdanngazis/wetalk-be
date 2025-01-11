package main

import (
	"chat-be/internal/config"
	"chat-be/internal/database"
	"chat-be/internal/delivery/http/handlers"
	"chat-be/internal/delivery/http/router"
	"chat-be/internal/domain/repositories"
	"chat-be/internal/kafka"
	"chat-be/internal/usecases"
	"chat-be/package/middleware"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize Database
	db := database.InitDB()

	// Initialize Migration
	database.InitMigration(db)

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	socketPathRepo := repositories.NewSocketPathRepository(db)
	chatRoomRepo := repositories.NewChatRoomRepository(db)

	// Initialize Usecases
	userUsecase := usecases.NewUserUsecase(userRepo, socketPathRepo)
	messageUsecase := usecases.NewMessageUsecase(chatRoomRepo, messageRepo, userRepo)
	chatRoomUsecase := usecases.NewChatRoomUsecase(chatRoomRepo, userRepo)

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userUsecase)
	messageHandler := handlers.NewMessageHandler(messageUsecase)
	chatRoomHandler := handlers.NewChatRoomHandler(chatRoomUsecase)

	kafkaService := kafka.NewKafkaService(messageUsecase)

	go kafkaService.ConsumeMessage()

	httpRouter := router.NewMuxRouter()
	httpRouter.POST("/api/users/login", userHandler.Login)
	httpRouter.OPTIONS("/api/users/login")
	httpRouter.POST("/api/users/register", userHandler.Register)
	httpRouter.OPTIONS("/api/users/register")
	httpRouter.GETWithMiddleware("/api/users/search", userHandler.SearchUsers, middleware.AuthMiddleware)
	httpRouter.OPTIONS("/api/users/search")
	//message
	httpRouter.GETWithMiddleware("/api/messages/history", messageHandler.GetMessageHistory, middleware.AuthMiddleware)
	httpRouter.OPTIONS("/api/messages/history")

	//room
	httpRouter.GETWithMiddleware("/api/rooms", chatRoomHandler.GetRooms, middleware.AuthMiddleware)
	httpRouter.OPTIONS("/api/rooms")
	httpRouter.POSTWithMiddleware("/api/rooms", chatRoomHandler.CreateRoom, middleware.AuthMiddleware)
	httpRouter.OPTIONS("/api/rooms")
	// Start Server
	port := config.GetEnv("APP_PORT", ":8080")
	httpRouter.SERVE(port)
}
