package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mensajesService/components/config"
	"mensajesService/message-api/model"
	"mensajesService/message-api/service"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageService struct {
	mock.Mock
}

var _ service.MessageServiceInterface = (*MockMessageService)(nil)

func (m *MockMessageService) CreateMessage(ctx context.Context, userID, content string) (*model.Message, error) {
	args := m.Called(ctx, userID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *MockMessageService) GetUserTimeline(ctx context.Context, userID string, limit int) ([]*model.TimelineItem, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.TimelineItem), args.Error(1)
}

func (m *MockMessageService) GetUserMessages(ctx context.Context, userID string, limit int) ([]*model.Message, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Message), args.Error(1)
}

func TestNewMessageController(t *testing.T) {
	mockService := &MockMessageService{}
	mockConfig := &config.AppConfig{}

	controller := NewMessageController(mockService, mockConfig)

	assert.NotNil(t, controller)
	assert.Equal(t, mockService, controller.messageService)
	assert.Equal(t, mockConfig, controller.config)
}

func TestCreateMessage_Success(t *testing.T) {
	mockService := &MockMessageService{}
	mockConfig := &config.AppConfig{
		MaxMessageLength: 280,
	}

	controller := NewMessageController(mockService, mockConfig)

	message := &model.Message{
		ID:        "test-id",
		UserID:    "user123",
		Content:   "Test message",
		CreatedAt: time.Now(),
	}

	mockService.On("CreateMessage", mock.Anything, "user123", "Test message").Return(message, nil)

	body, _ := json.Marshal(map[string]string{"content": "Test message"})
	req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "user123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test message", response.Content)

	mockService.AssertExpectations(t)
}

func TestCreateMessage_MissingUserID(t *testing.T) {
	mockService := &MockMessageService{}
	mockConfig := &config.AppConfig{
		MaxMessageLength: 280,
	}

	controller := NewMessageController(mockService, mockConfig)

	body, _ := json.Marshal(map[string]string{"content": "Test message"})
	req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "User ID required")
}

func TestCreateMessage_EmptyContent(t *testing.T) {
	mockService := &MockMessageService{}
	mockConfig := &config.AppConfig{
		MaxMessageLength: 280,
	}

	controller := NewMessageController(mockService, mockConfig)

	body, _ := json.Marshal(map[string]string{"content": ""})
	req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "user123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Content is required")
}

func TestCreateMessage_ContentTooLong(t *testing.T) {
	mockService := &MockMessageService{}
	mockConfig := &config.AppConfig{
		MaxMessageLength: 10,
	}

	controller := NewMessageController(mockService, mockConfig)

	body, _ := json.Marshal(map[string]string{"content": "Un mensaje muyyyyy largo que supera el limite de 10 caracteres"})
	req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "user123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Content too long")
}

func TestGetUserMessages_Success(t *testing.T) {
	mockService := &MockMessageService{}
	mockConfig := &config.AppConfig{
		DefaultLimit: 20,
	}

	controller := NewMessageController(mockService, mockConfig)

	messages := []*model.Message{
		{
			ID:        "msg1",
			UserID:    "user123",
			Content:   "Test content 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        "msg2",
			UserID:    "user123",
			Content:   "Test content 2",
			CreatedAt: time.Now(),
		},
	}

	mockService.On("GetUserMessages", mock.Anything, "user123", 20).Return(messages, nil)

	req := httptest.NewRequest("GET", "/messages", nil)
	req.Header.Set("X-User-ID", "user123")

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*model.Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "msg1", response[0].ID)

	mockService.AssertExpectations(t)
}

func TestGetUserMessages_MissingUserID(t *testing.T) {
	mockService := &MockMessageService{}
	mockConfig := &config.AppConfig{
		DefaultLimit: 20,
	}

	controller := NewMessageController(mockService, mockConfig)

	req := httptest.NewRequest("GET", "/messages", nil)

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "User ID required")
}
