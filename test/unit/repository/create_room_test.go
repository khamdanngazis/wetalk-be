package repository_test

import (
	"chat-be/internal/domain/entities"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateRoom(t *testing.T) {
	userIDs := []string{"ccffa72d-8a9f-463d-8f76-aa2edbba7b8b", "e4714647-38ce-4cef-8934-2407f9a4f923"}
	room := &entities.ChatRoom{
		ID:      uuid.New().String(),
		Name:    "a",
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

	err := chatRoomRepo.CreateRoom(room, participants)
	fmt.Println(err)
	fmt.Println(room)
}
