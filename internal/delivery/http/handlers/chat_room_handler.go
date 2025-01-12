// internal/delivery/http/message_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"chat-be/internal/delivery/http/models"
	"chat-be/internal/usecases"
	"chat-be/package/helper"
	"chat-be/package/middleware"

	"github.com/go-playground/validator/v10"
)

type ChatRoomHandler struct {
	ChatRoomUsecase usecases.ChatRoomUsecase
}

func NewChatRoomHandler(chatRoomUsecase usecases.ChatRoomUsecase) *ChatRoomHandler {
	return &ChatRoomHandler{ChatRoomUsecase: chatRoomUsecase}
}

func (h *ChatRoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteResponse(w, http.StatusUnauthorized, "User token not valid", nil)
		return
	}

	var request models.CreateChatRoomRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, helper.GetMessageValidator(validate, err), nil)
		return
	}

	room, err := h.ChatRoomUsecase.CreateRoom(user.UserID, request.UserIDs, request.IsGroup, request.RoomName)

	if err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, err.Error(), nil)
	}

	middleware.WriteResponse(w, http.StatusCreated, "Success create room", room)
}

func (h *ChatRoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteResponse(w, http.StatusUnauthorized, "User token not valid", nil)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default to page 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default value
	}

	rooms, total, err := h.ChatRoomUsecase.GetRoomsForUser(user.UserID, page, limit)

	if err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, err.Error(), nil)
	}

	if rooms == nil {
		rooms = []models.GetChatRoomResponse{} // Assuming Room is the type of the elements in the rooms slice
	}

	response := map[string]interface{}{
		"rooms": rooms,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	middleware.WriteResponse(w, http.StatusOK, "Success get rooms", response)
}
