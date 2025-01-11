package entities

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID            string          `gorm:"type:uuid;primaryKey" json:"id"`
	ChatRoomID    string          `gorm:"type:uuid;not null" json:"chat_room_id"` // Referensi ke ChatRoom
	SenderID      string          `gorm:"type:uuid;not null" json:"sender_id"`    // ID pengirim pesan
	Content       string          `gorm:"not null" json:"content"`
	Status        int             `gorm:"not null" json:"status"`
	MessageStatus []MessageStatus `gorm:"foreignKey:MessageID;references:ID" json:"message_status"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt  `gorm:"index"`
}

type MessageStatus struct {
	ID         string         `gorm:"type:uuid;primaryKey" json:"id"`
	MessageID  string         `gorm:"type:uuid;not null" json:"message_id"`
	ReceiverID string         `gorm:"type:uuid;not null" json:"receiver_id"`
	Status     int            `gorm:"not null" json:"status"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

const (
	StatusSend      = 1
	StatusDelivered = 2
	StatusRead      = 3
)
