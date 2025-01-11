package repositories

import (
	"chat-be/internal/domain/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id string) (*entities.User, error)
	CountUsersBySocketID(socketID string) (int64, error)
	SearchByUsernameOrEmail(query string) ([]*entities.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Preload("SocketPath").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*entities.User, error) {
	var user entities.User
	err := r.db.Preload("SocketPath").Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CountUsersBySocketID(socketID string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.User{}).Where("socket_id = ?", socketID).Count(&count).Error
	return count, err
}

func (r *userRepository) SearchByUsernameOrEmail(query string) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.Preload("SocketPath").
		Where("username LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
