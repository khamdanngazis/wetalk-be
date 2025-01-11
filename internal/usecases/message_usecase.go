// internal/usecases/message_usecase.go
package usecases

import (
	"chat-be/internal/delivery/http/models"
	"chat-be/internal/domain/entities"
	"chat-be/internal/domain/repositories"
	"errors"

	"github.com/google/uuid"
)

type MessageUsecase interface {
	GetMessageHistory(senderID, roomId string, limit, page int) ([]models.Message, int64, error)
	SaveMessage(message *entities.Message) error
	UpdateStatusMessage(messageID, receiverID string, status int) error
}

type messageUsecase struct {
	chatRoom    repositories.ChatRoomRepository
	messageRepo repositories.MessageRepository
	userRepo    repositories.UserRepository
}

func NewMessageUsecase(chatRoom repositories.ChatRoomRepository, messageRepo repositories.MessageRepository, userRepo repositories.UserRepository) MessageUsecase {
	return &messageUsecase{
		chatRoom:    chatRoom,
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

func (m *messageUsecase) GetMessageHistory(senderID, roomId string, limit, page int) ([]models.Message, int64, error) {
	participants, err := m.chatRoom.FindUsersByRoomID(roomId)
	if err != nil {
		return nil, 0, err
	}
	validSender := false
	for _, v := range participants {
		if v.UserID == senderID {
			validSender = true
			break
		}
	}

	if !validSender {
		return nil, 0, errors.New("invalid sender")
	}

	offset := (page - 1) * limit

	messageHistories, total, err := m.messageRepo.GetMessagesByRoomID(roomId, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	var message models.Message
	var messages []models.Message
	for _, v := range messageHistories {
		message.ID = v.ID
		message.RoomID = v.ChatRoomID
		if len(v.MessageStatus) == 1 {
			message.Status = v.MessageStatus[0].Status
		}

		if v.SenderID == senderID {
			message.Type = "outgoing"
		} else {
			message.Type = "incoming"
		}
		message.Time = v.CreatedAt.Format("2006-01-02 15:04")
		message.Text = v.Content
		messages = append(messages, message)
	}

	return messages, total, nil
}

func (m *messageUsecase) SaveMessage(message *entities.Message) error {
	sender, err := m.userRepo.FindByID(message.SenderID)
	if err != nil || sender == nil {
		return errors.New("invalid sender")
	}

	receiver, err := m.chatRoom.FindRoomByID(message.ChatRoomID)
	if err != nil || receiver == nil {
		return errors.New("invalid receiver")
	}
	err = m.messageRepo.SaveMessage(message)
	if err != nil {
		return err
	}

	for _, v := range receiver.Participants {
		if v.UserID != message.SenderID {
			messageStatus := entities.MessageStatus{
				ID:         uuid.New().String(),
				MessageID:  message.ID,
				ReceiverID: v.UserID,
				Status:     entities.StatusSend,
			}
			err = m.messageRepo.CreateMessageStatus(&messageStatus)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *messageUsecase) UpdateStatusMessage(messageID, receiverID string, status int) error {
	return m.messageRepo.UpdateMessageStatus(messageID, receiverID, status)
}
