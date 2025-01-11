// internal/delivery/http/message_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"chat-be/internal/delivery/http/models"
	"chat-be/internal/usecases"
	"chat-be/package/middleware"
)

type MessageHandler struct {
	MessageUsecase usecases.MessageUsecase
}

func NewMessageHandler(messageUsecase usecases.MessageUsecase) *MessageHandler {
	return &MessageHandler{MessageUsecase: messageUsecase}
}

func (h *MessageHandler) GetMessageHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteResponse(w, http.StatusUnauthorized, "User token not valid", nil)
		return
	}

	roomID := r.URL.Query().Get("room_id")
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	// Validate sender_id and receiver_id
	if roomID == "" {
		middleware.WriteResponse(w, http.StatusBadRequest, "sender_id and receiver_id are required", nil)
		return
	}

	// Parse limit and offset
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default value
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default to page 1
	}

	// Fetch messages
	messages, total, err := h.MessageUsecase.GetMessageHistory(user.UserID, roomID, limit, page)
	if err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to fetch messages", nil)
		return
	}

	if messages == nil {
		messages = []models.Message{} // Replace User with the actual type if necessary
	}

	response := map[string]interface{}{
		"messages": messages,
		"total":    total,
		"page":     page,
		"limit":    limit,
	}

	middleware.WriteResponse(w, http.StatusOK, "Chat history fetched", response)
}
