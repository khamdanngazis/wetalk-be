package usecases

import (
	"chat-be/internal/delivery/http/models"
	"chat-be/internal/domain/entities"
	"chat-be/internal/domain/repositories"
	"chat-be/package/helper"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ChatRoomUsecase interface {
	CreateRoom(userCreator string, userIDs []string, isGroup bool, roomName string) (*models.GetChatRoomResponse, error)
	FindUsersByRoomID(roomID string) ([]entities.ChatRoomParticipant, error)
	GetRoomsForUser(userID string, page int, limit int) ([]models.GetChatRoomResponse, int64, error)
}

type chatRoomUsecase struct {
	chatRoomRepo repositories.ChatRoomRepository
	userRepo     repositories.UserRepository
}

func NewChatRoomUsecase(chatRoomRepo repositories.ChatRoomRepository, userRepo repositories.UserRepository) ChatRoomUsecase {
	return &chatRoomUsecase{
		chatRoomRepo: chatRoomRepo,
		userRepo:     userRepo,
	}
}

func (u *chatRoomUsecase) CreateRoom(userIdCreator string, userIDs []string, isGroup bool, roomName string) (*models.GetChatRoomResponse, error) {
	if len(userIDs) < 2 {
		return nil, errors.New("at least two participants are required to create a chat room")
	}
	var participantList []models.Participants
	var users []entities.User
	for _, v := range userIDs {
		user, err := u.userRepo.FindByID(v)
		if err != nil || user == nil {
			return nil, errors.New("invalid sender")
		}
		users = append(users, *user)
		participant := models.Participants{
			UserID:     user.ID,
			SocketPath: user.SocketPath.Path,
		}
		participantList = append(participantList, participant)
	}

	if len(userIDs) == 2 && !isGroup {

		for _, v := range users {
			if v.ID != userIdCreator {
				roomName = v.Username
			}
		}

		existingRoom, err := u.chatRoomRepo.FindRoomByParticipants(userIDs)
		if err != nil {
			return nil, err
		}

		if existingRoom != nil {

			chatRoom := mappingChatRoom(*existingRoom, userIdCreator)
			chatRoom.Participants = participantList
			chatRoom.Name = roomName
			return &chatRoom, nil
		}
	}

	room := &entities.ChatRoom{
		ID:      uuid.New().String(),
		Name:    roomName,
		IsGroup: len(userIDs) > 2,
	}

	var participants []entities.ChatRoomParticipant
	for _, userID := range userIDs {
		participants = append(participants, entities.ChatRoomParticipant{
			ID:         uuid.New().String(),
			UserID:     userID,
			JoinedAt:   time.Now(),
			ChatRoomID: room.ID,
		})
	}

	err := u.chatRoomRepo.CreateRoom(room, participants)
	if err != nil {
		return nil, err
	}

	chatRoom := mappingChatRoom(*room, userIdCreator)
	chatRoom.Participants = participantList
	chatRoom.Name = roomName
	return &chatRoom, nil
}

func (u *chatRoomUsecase) FindUsersByRoomID(roomID string) ([]entities.ChatRoomParticipant, error) {
	return u.chatRoomRepo.FindUsersByRoomID(roomID)
}

func (u *chatRoomUsecase) GetRoomsForUser(userID string, page int, limit int) ([]models.GetChatRoomResponse, int64, error) {
	offset := (page - 1) * limit

	rooms, total, err := u.chatRoomRepo.FindRoomsByUser(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var chatRooms []models.GetChatRoomResponse
	for _, v := range rooms {
		chatRoom := mappingChatRoom(v, userID)
		chatRooms = append(chatRooms, chatRoom)
	}

	return chatRooms, total, nil
}

func mappingChatRoom(room entities.ChatRoom, userID string) models.GetChatRoomResponse {
	var chatRoom models.GetChatRoomResponse
	chatRoom.ID = room.ID
	name := ""
	for _, a := range room.Participants {
		if a.User.ID != userID {
			name = a.User.Username
		}
	}

	chatRoom.Name = name
	if room.LastMessageID != nil {
		chatRoom.LastMessage = room.Message.Content
		chatRoom.LastMessageTime = helper.FormatMessageTime(room.Message.CreatedAt)
	}

	for _, v := range room.Participants {
		participant := models.Participants{
			UserID:     v.User.ID,
			SocketPath: v.User.SocketPath.Path,
		}
		chatRoom.Participants = append(chatRoom.Participants, participant)
	}

	return chatRoom
}
