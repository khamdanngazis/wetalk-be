package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         string         `gorm:"type:uuid;primaryKey"`
	Username   string         `gorm:"unique;not null" json:"username" validate:"required"`
	Email      string         `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password   string         `gorm:"not null" json:"password" validate:"required,min=8"`
	SocketID   string         `gorm:"type:uuid" json:"socket_id"`
	SocketPath SocketPath     `gorm:"foreignKey:SocketID;references:ID"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type UserResponse struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
