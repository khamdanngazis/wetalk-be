package usecase_test

import (
	"fmt"
	"testing"
)

func TestCreateRoom(t *testing.T) {
	roomName := "a|b"
	isGroup := false
	userIDs := []string{"ccffa72d-8a9f-463d-8f76-aa2edbba7b8b", "e4714647-38ce-4cef-8934-2407f9a4f923"}

	room, err := chatRoomUsecase.CreateRoom("ccffa72d-8a9f-463d-8f76-aa2edbba7b8b", userIDs, isGroup, roomName)
	fmt.Println(err)
	fmt.Println(room)
}
