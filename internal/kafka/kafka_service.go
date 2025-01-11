package kafka

import (
	"chat-be/internal/config"
	"chat-be/internal/domain/entities"
	"chat-be/internal/usecases"
	"chat-be/package/logging"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaService struct {
	Reader         *kafka.Reader
	MessageUsecase usecases.MessageUsecase
}

func NewKafkaService(messageUsecase usecases.MessageUsecase) *KafkaService {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.GetEnv("KAFKA_HOST", "localhost:9092")},
		Topic:   config.GetEnv("KAFKA_TOPIC", "chat"),
		GroupID: "chat-be-group",
	})
	return &KafkaService{Reader: reader, MessageUsecase: messageUsecase}
}

func (k *KafkaService) checkKafkaConnection(ctx context.Context) error {
	// Attempt to read metadata to ensure connectivity
	_, err := k.Reader.FetchMessage(ctx)
	if err != nil {
		logging.LogError(ctx, "Kafka connection failed: %v", err)
		return err
	}
	return nil
}

func (k *KafkaService) ConsumeMessage() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Check Kafka connection before starting
	if err := k.checkKafkaConnection(ctx); err != nil {
		logging.LogError(ctx, "Failed to connect to Kafka. Exiting...")
		return
	}

	for {
		requestID := uuid.New().String()
		loopCtx := context.WithValue(ctx, logging.RequestIDKey, requestID)

		msg, err := k.Reader.ReadMessage(loopCtx)
		if err != nil {
			logging.LogError(loopCtx, "Error while reading message: %v", err)
			logging.LogError(loopCtx, "Stopping message consumption due to Kafka error.")
			return
		}

		logging.LogInfo(loopCtx, "Incoming Kafka message: %v", string(msg.Value))

		key := string(msg.Key)

		switch key {
		case "message":

			var messageModel entities.Message
			if err := json.Unmarshal(msg.Value, &messageModel); err != nil {
				logging.LogError(loopCtx, "Failed to parse message: %v", err)
				continue
			}

			logging.LogInfo(loopCtx, "Incoming Kafka message model: %v", messageModel)

			// Save the message

			if err := k.MessageUsecase.SaveMessage(&messageModel); err != nil {
				logging.LogError(loopCtx, "Error while saving message: %v", err)
				continue
			}
		case "update_status":
			var messageStatusModel entities.MessageStatus
			if err := json.Unmarshal(msg.Value, &messageStatusModel); err != nil {
				logging.LogError(loopCtx, "Failed to parse message status: %v", err)
				continue
			}

			logging.LogInfo(loopCtx, "Incoming Kafka message status model: %v", messageStatusModel)
			if err := k.MessageUsecase.UpdateStatusMessage(messageStatusModel.MessageID, messageStatusModel.ReceiverID, messageStatusModel.Status); err != nil {
				logging.LogError(loopCtx, "Error while updating message status: %v", err)
				continue
			}
		default:
			logging.LogError(loopCtx, "Unknown message key: %v", key)
			continue
		}
		// Acknowledge that the message has been successfully saved
		logging.LogInfo(loopCtx, "Message saved and processed successfully. ACK.")

		// Commit the offset to mark the message as consumed
		if err := k.Reader.CommitMessages(loopCtx, msg); err != nil {
			logging.LogError(loopCtx, "Failed to commit Kafka message offset: %v", err)
		} else {
			logging.LogInfo(loopCtx, "Message offset committed successfully.")
		}
	}
}
