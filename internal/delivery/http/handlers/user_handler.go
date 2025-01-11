package handlers

import (
	"encoding/json"
	"net/http"

	"chat-be/internal/domain/entities"
	"chat-be/internal/usecases"
	"chat-be/package/helper"
	"chat-be/package/logging"
	"chat-be/package/middleware"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	UserUsecase usecases.UserUsecase
}

func NewUserHandler(userUsecase usecases.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user entities.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, helper.GetMessageValidator(validate, err), nil)
		return
	}

	// Call Usecase to register user
	err := h.UserUsecase.Register(&user)
	if err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	middleware.WriteResponse(w, http.StatusCreated, "User registered successfully", nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, helper.GetMessageValidator(validate, err), nil)
		return
	}

	token, err := h.UserUsecase.Login(req.Email, req.Password)
	if err != nil {
		logging.LogError(ctx, "Login error: %v", err)
		middleware.WriteResponse(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "", map[string]string{"token": token})
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteResponse(w, http.StatusUnauthorized, "User token not valid", nil)
		return
	}
	// Extract query parameter
	query := r.URL.Query().Get("q")
	if query == "" {
		middleware.WriteResponse(w, http.StatusBadRequest, "Query parameter 'q' is required", nil)
		return
	}

	// Search for users using the usecase
	users, err := h.UserUsecase.SearchUsers(query, user.UserID)
	if err != nil {
		logging.LogError(ctx, "Search error: %v", err)
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to search users", nil)
		return
	}

	if users == nil {
		users = []entities.UserResponse{} // Replace User with the actual type if necessary
	}

	// You can also handle the case where users are empty, if required
	if len(users) == 0 {
		// Optionally, return a message indicating no users were found
		middleware.WriteResponse(w, http.StatusOK, "No users found", users)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "Users fetched successfully", users)
}
