package entities

import (
	"time"

	"gorm.io/gorm"
)

type ChatRoom struct {
	ID            string                `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string                `gorm:"type:varchar(255)" json:"name"`
	IsGroup       bool                  `gorm:"not null;default:false" json:"is_group"`
	LastMessageID *string               `gorm:"type:uuid;null" json:"last_message_id"`
	Message       Message               `gorm:"foreignKey:LastMessageID;references:ID"`
	Participants  []ChatRoomParticipant `gorm:"foreignKey:ChatRoomID" json:"participants"`
	CreatedAt     time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt        `gorm:"index" json:"deleted_at"`
}

type ChatRoomParticipant struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	ChatRoomID string    `gorm:"type:uuid;not null" json:"chat_room_id"`
	UserID     string    `gorm:"type:uuid;not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID;references:ID"`
	JoinedAt   time.Time `gorm:"not null" json:"joined_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
