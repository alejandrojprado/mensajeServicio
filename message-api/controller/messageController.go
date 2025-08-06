package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mensajesService/components/config"
	"mensajesService/components/logger"
	"mensajesService/components/metrics"
	"mensajesService/message-api/model"
	"mensajesService/message-api/service"

	"github.com/go-chi/chi/v5"
)

type MessageController struct {
	messageService service.MessageServiceInterface
	config         *config.AppConfig
}

func NewMessageController(messageService service.MessageServiceInterface, cfg *config.AppConfig) *MessageController {
	return &MessageController{
		messageService: messageService,
		config:         cfg,
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
		http.Error(w, "Content too long (max "+strconv.Itoa(c.config.MaxMessageLength)+" characters)", http.StatusBadRequest)
		return
	}

	createdMessage, err := c.messageService.CreateMessage(r.Context(), userID, message.Content)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricMessageError, 1)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.LogError("CreateMessage error", "error", err, "user_id", userID)
		return
	}

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

	limitStr := r.URL.Query().Get("limit")
	limit := c.config.DefaultLimit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	messages, err := c.messageService.GetUserMessages(r.Context(), userID, limit)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricUserMessagesError, 1)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.LogError("GetUserMessages error", "error", err, "user_id", userID)
		return
	}

	metrics.PutCountMetric(metrics.MetricUserMessagesSuccess, 1)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
