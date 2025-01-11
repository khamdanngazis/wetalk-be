// internal/domain/repositories/message_repository.go
package repositories

import (
	"chat-be/internal/domain/entities"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *entities.Message) error
	SaveMessage(message *entities.Message) error
	CreateMessageStatus(messageStatus *entities.MessageStatus) error
	GetMessagesByRoomID(chatRoomID string, offset int, limit int) ([]entities.Message, int64, error)
	UpdateMessageStatus(messageID string, receiverID string, status int) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) Create(message *entities.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) SaveMessage(message *entities.Message) error {
	var existingMessage entities.Message
	err := r.db.Where("id = ?", message.ID).First(&existingMessage).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Message does not exist, create a new one
			return r.db.Create(message).Error
		}
		// Other database error
		return err
	}

	// Message exists, update its status
	existingMessage.Status = message.Status
	return r.db.Save(&existingMessage).Error
}

func (r *messageRepository) CreateMessageStatus(messageStatus *entities.MessageStatus) error {
	return r.db.Create(messageStatus).Error
}

func (r *messageRepository) GetMessageStatus(messageID string, receiverID string) (*entities.MessageStatus, error) {
	var messageStatus entities.MessageStatus
	err := r.db.Where("message_id = ? AND receiver_id = ?", messageID, receiverID).First(&messageStatus).Error
	if err != nil {
		return nil, err
	}

	return &messageStatus, nil
}

func (r *messageRepository) UpdateMessageStatus(messageID string, receiverID string, status int) error {
	return r.db.Model(&entities.MessageStatus{}).
		Where("message_id = ? AND receiver_id = ?", messageID, receiverID).
		Update("status", status).Error
}

func (r *messageRepository) GetMessagesByRoomID(chatRoomID string, offset int, limit int) ([]entities.Message, int64, error) {
	var messages []entities.Message
	var totalRows int64

	// Count total rows
	err := r.db.Model(&entities.Message{}).
		Where("chat_room_id = ?", chatRoomID).
		Count(&totalRows).Error
	if err != nil {
		return nil, 0, err
	}

	// Query messages with offset, limit, and preload MessageStatus
	err = r.db.Preload("MessageStatus").
		Where("chat_room_id = ?", chatRoomID).
		Offset(offset).
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, totalRows, nil
}
