package entities

import (
	"time"

	"gorm.io/gorm"
)

type SocketPath struct {
	ID        string         `gorm:"type:uuid;primaryKey"`
	Path      string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
