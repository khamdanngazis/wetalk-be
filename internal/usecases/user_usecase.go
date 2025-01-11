package usecases

import (
	"errors"
	"fmt"

	"chat-be/internal/domain/entities"
	"chat-be/internal/domain/repositories"
	"chat-be/package/helper"

	"github.com/google/uuid"
)

type UserUsecase interface {
	Register(user *entities.User) error
	Login(email, password string) (string, error)
	SearchUsers(query string, userID string) ([]entities.UserResponse, error)
}

type userUsecase struct {
	userRepo       repositories.UserRepository
	socketPathRepo repositories.SocketPathRepository
}

func NewUserUsecase(userRepo repositories.UserRepository, socketPathRepo repositories.SocketPathRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo, socketPathRepo: socketPathRepo}
}

func (u *userUsecase) Register(user *entities.User) error {

	// Check if email or username already exists
	existingUser, _ := u.userRepo.FindByEmail(user.Email)

	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Generate unique user ID
	user.ID = uuid.New().String()

	// Hash the password
	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	// Assign available socket path ID
	socketID, err := u.getAvailableSocketPathID()
	if err != nil {
		return errors.New("failed to assign socket path: " + err.Error())
	}
	user.SocketID = socketID

	// Save to database
	err = u.userRepo.Create(user)
	if err != nil {
		return errors.New("failed to register user: " + err.Error())
	}

	return nil
}

func (u *userUsecase) Login(email, password string) (string, error) {
	// Find user by email
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Check password
	err = helper.CompareHashAndPassword(user.Password, password)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token, err := helper.GenerateToken(user.ID, user.Email, user.Username, user.SocketPath.Path)
	if err != nil {
		return "", errors.New("internal server error")
	}
	return token, nil
}

func (u *userUsecase) SearchUsers(query string, userID string) ([]entities.UserResponse, error) {
	// Use repository to search for users by username or email
	users, err := u.userRepo.SearchByUsernameOrEmail(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	var userResponses []entities.UserResponse
	for _, user := range users {
		if user.ID != userID {
			userResponses = append(userResponses, entities.UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			})
		}
	}
	return userResponses, nil

}

func (u *userUsecase) getAvailableSocketPathID() (string, error) {
	socketPaths, err := u.socketPathRepo.FindAll()
	if err != nil {
		return "", fmt.Errorf("failed to fetch socket paths: %w", err)
	}

	for _, socketPath := range socketPaths {
		userCount, err := u.userRepo.CountUsersBySocketID(socketPath.ID)
		if err != nil {
			return "", fmt.Errorf("failed to count users for socket ID %s: %w", socketPath.ID, err)
		}
		if userCount < 1000 {
			return socketPath.ID, nil
		}
	}

	newSocketPathID, err := u.createNewSocketPath()
	if err != nil {
		return "", fmt.Errorf("failed to create a new socket path: %w", err)
	}

	return newSocketPathID, nil
}

func (u *userUsecase) createNewSocketPath() (string, error) {
	// Generate a new UUID for the socket path
	newSocketPathID := uuid.New().String()

	// Create a new socket path entry
	newSocketPath := &entities.SocketPath{
		ID:   newSocketPathID,
		Path: fmt.Sprintf("/ws/%s", newSocketPathID), // Example path format
	}

	// Save the new socket path to the database
	err := u.socketPathRepo.Create(newSocketPath)
	if err != nil {
		return "", fmt.Errorf("failed to save new socket path: %w", err)
	}

	return newSocketPathID, nil
}
