package models

type CreateChatRoomRequest struct {
	UserIDs  []string `json:"user_ids" validate:"required"`
	IsGroup  bool     `json:"is_group"`
	RoomName string   `json:"room_name"`
}

type GetChatRoomResponse struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	LastMessage     string         `json:"last_message"`
	LastMessageTime string         `json:"last_message_time"`
	Participants    []Participants `json:"participants"`
}

type Participants struct {
	UserID     string `json:"user_id"`
	SocketPath string `json:"socket_path"`
}
