package models

type Message struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	Type   string `json:"type"`
	Text   string `json:"text"`
	Time   string `json:"time"`
	Status int    `json:"status"`
}
