package repositories

import (
	"chat-be/internal/domain/entities"
	"errors"

	"gorm.io/gorm"
)

type chatRoomRepository struct {
	db *gorm.DB
}

type ChatRoomRepository interface {
	CreateRoom(room *entities.ChatRoom, participants []entities.ChatRoomParticipant) error
	FindRoomByParticipants(userIDs []string) (*entities.ChatRoom, error)
	FindRoomsByUser(userID string, offset int, limit int) ([]entities.ChatRoom, int64, error)
	FindUsersByRoomID(roomID string) ([]entities.ChatRoomParticipant, error)
	FindRoomByID(ID string) (*entities.ChatRoom, error)
}

func NewChatRoomRepository(db *gorm.DB) ChatRoomRepository {
	return &chatRoomRepository{db: db}
}

func (r *chatRoomRepository) CreateRoom(room *entities.ChatRoom, participants []entities.ChatRoomParticipant) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(room).Error; err != nil {
			return err
		}
		if len(participants) > 0 {
			if err := tx.Create(&participants).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *chatRoomRepository) FindRoomByParticipants(userIDs []string) (*entities.ChatRoom, error) {
	if len(userIDs) != 2 {
		return nil, errors.New("exactly two user IDs are required")
	}

	var room entities.ChatRoom
	subQuery := r.db.Table("chat_room_participants").
		Select("chat_room_id").
		Where("user_id IN ?", userIDs).
		Group("chat_room_id").
		Having("COUNT(DISTINCT user_id) = ?", len(userIDs))

	err := r.db.Table("chat_rooms").
		Where("id IN (?) AND is_group = ?", subQuery, false).
		First(&room).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *chatRoomRepository) FindRoomsByUser(userID string, offset int, limit int) ([]entities.ChatRoom, int64, error) {
	var rooms []entities.ChatRoom
	var totalRows int64

	// Count total rows
	err := r.db.Table("chat_rooms").
		Joins("JOIN chat_room_participants c ON c.chat_room_id = chat_rooms.id").
		Where("c.user_id = ?", userID).
		Count(&totalRows).Error
	if err != nil {
		return nil, 0, err
	}

	// Query rooms with offset and limit
	err = r.db.Joins("JOIN chat_room_participants c ON c.chat_room_id = chat_rooms.id").
		Where("c.user_id = ?", userID).
		Preload("Participants.User.SocketPath"). // Preload participants, users, and socket paths
		Offset(offset).
		Limit(limit).
		Find(&rooms).Error
	if err != nil {
		return nil, 0, err
	}

	return rooms, totalRows, nil
}

func (r *chatRoomRepository) FindRoomByID(ID string) (*entities.ChatRoom, error) {
	var room entities.ChatRoom

	// Query the chat room by ID
	err := r.db.Where("id = ?", ID).Preload("Participants.User.SocketPath").Take(&room).Error
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *chatRoomRepository) FindUsersByRoomID(roomID string) ([]entities.ChatRoomParticipant, error) {
	var participants []entities.ChatRoomParticipant

	err := r.db.Preload("User").Preload("User.SocketPath").Where("chat_room_id = ?", roomID).Find(&participants).Error
	if err != nil {
		return nil, err
	}

	return participants, nil
}
