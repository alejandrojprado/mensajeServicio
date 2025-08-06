package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"mensajesService/components/config"
	"mensajesService/components/logger"
	"mensajesService/components/metrics"
	"mensajesService/message-api/model"
	"mensajesService/message-api/service"

	"github.com/go-chi/chi/v5"
)

type MessageController struct {
	messageService  service.MessageServiceInterface
	timelineService service.TimelineServiceInterface
	config          *config.AppConfig
}

func NewMessageController(messageService service.MessageServiceInterface, timelineService service.TimelineServiceInterface, cfg *config.AppConfig) *MessageController {
	return &MessageController{
		messageService:  messageService,
		timelineService: timelineService,
		config:          cfg,
	}
}

func (c *MessageController) MountIn(r chi.Router) {
	r.Route("/messages", func(r chi.Router) {
		r.Post("/", c.CreateMessage)
		r.Get("/", c.GetUserMessages)
	})
}

func (c *MessageController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		metrics.PutCountMetric(metrics.MetricMessageError, 1)
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	var message model.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		metrics.PutCountMetric(metrics.MetricMessageError, 1)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if message.Content == "" {
		metrics.PutCountMetric(metrics.MetricMessageError, 1)
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	if len(message.Content) > c.config.MaxMessageLength {
		metrics.PutCountMetric(metrics.MetricMessageError, 1)
		http.Error(w, "Content too long", http.StatusBadRequest)
		return
	}

	createdMessage, err := c.messageService.CreateMessage(r.Context(), userID, message.Content)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricMessageError, 1)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.LogError("CreateMessage error", "error", err, "user_id", userID)
		return
	}

	go func() {
		if err := c.timelineService.UpdateFollowersTimeline(context.Background(), createdMessage); err != nil {
			logger.LogError("Error updating followers timeline", "error", err, "message_id", createdMessage.ID)
		}
	}()

	metrics.PutCountMetric(metrics.MetricMessageSuccess, 1)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMessage)
}

func (c *MessageController) GetUserMessages(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		metrics.PutCountMetric(metrics.MetricMessageError, 1)
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	messages, err := c.messageService.GetUserMessages(r.Context(), userID, c.config.DefaultLimit)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricUserMessagesError, 1)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.LogError("GetUserMessages error", "error", err, "user_id", userID)
		return
	}
	messagesAmount := float64(len(messages))
	metrics.PutCountMetric(metrics.MetricUserMessagesSuccess, 1)
	metrics.PutCountMetric(metrics.MetricUserMessagesCount, messagesAmount)
	logger.LogInfo("GetUserMessages success", "user_id", userID, "messages_amount", messagesAmount)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
