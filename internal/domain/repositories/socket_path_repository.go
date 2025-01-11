package repositories

import (
	"chat-be/internal/domain/entities"

	"gorm.io/gorm"
)

type SocketPathRepository interface {
	Create(socketPath *entities.SocketPath) error
	FindByID(id string) (*entities.SocketPath, error)
	Delete(id string) error
	FindAll() ([]entities.SocketPath, error)
}

type socketPathRepository struct {
	db *gorm.DB
}

func NewSocketPathRepository(db *gorm.DB) SocketPathRepository {
	return &socketPathRepository{db}
}
func (r *socketPathRepository) Create(socketPath *entities.SocketPath) error {
	return r.db.Create(socketPath).Error
}

func (r *socketPathRepository) FindByID(id string) (*entities.SocketPath, error) {
	var socketPath entities.SocketPath
	err := r.db.Where("id = ?", id).First(&socketPath).Error
	if err != nil {
		return nil, err
	}
	return &socketPath, nil
}

func (r *socketPathRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.SocketPath{}).Error
}

func (r *socketPathRepository) FindAll() ([]entities.SocketPath, error) {
	var socketPaths []entities.SocketPath
	err := r.db.Find(&socketPaths).Error
	return socketPaths, err
}
