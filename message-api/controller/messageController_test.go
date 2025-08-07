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
	"mensajesService/components/logger"
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

func (m *MockMessageService) GetUserMessages(ctx context.Context, userID string, limit int) ([]*model.Message, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Message), args.Error(1)
}

func TestNewMessageController(t *testing.T) {
	logger.Init()

	mockService := &MockMessageService{}
	mockTimelineService := &MockTimelineService{}
	mockConfig := &config.AppConfig{}

	controller := NewMessageController(mockService, mockTimelineService, mockConfig)

	assert.NotNil(t, controller)
	assert.Equal(t, mockService, controller.messageService)
	assert.Equal(t, mockTimelineService, controller.timelineService)
	assert.Equal(t, mockConfig, controller.config)
}

func TestCreateMessage_Success(t *testing.T) {
	mockService := &MockMessageService{}
	mockTimelineService := &MockTimelineService{}
	mockConfig := &config.AppConfig{
		MaxMessageLength: 280,
	}

	controller := NewMessageController(mockService, mockTimelineService, mockConfig)

	message := &model.Message{
		ID:        "test-id",
		UserID:    "user123",
		Content:   "Test message",
		CreatedAt: time.Now(),
	}

	mockService.On("CreateMessage", mock.Anything, "user123", "Test message").Return(message, nil)
	mockTimelineService.On("UpdateFollowersTimeline", mock.Anything, message).Return(nil)

	body, _ := json.Marshal(map[string]string{"content": "Test message"})
	req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "user123")
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusCreated, response.Code)

	var messageResponse model.Message
	err := json.Unmarshal(response.Body.Bytes(), &messageResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Test message", messageResponse.Content)

	// Wait a bit for the goroutine to complete. It could be better
	time.Sleep(100 * time.Millisecond)

	mockService.AssertExpectations(t)
	mockTimelineService.AssertExpectations(t)
}

func TestCreateMessage_MissingUserID(t *testing.T) {
	mockService := &MockMessageService{}
	mockTimelineService := &MockTimelineService{}
	mockConfig := &config.AppConfig{
		MaxMessageLength: 280,
	}

	controller := NewMessageController(mockService, mockTimelineService, mockConfig)

	body, _ := json.Marshal(map[string]string{"content": "Test message"})
	req := httptest.NewRequest("POST", "/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "User ID required")
}

func TestGetUserMessages_Success(t *testing.T) {
	mockService := &MockMessageService{}
	mockTimelineService := &MockTimelineService{}
	mockConfig := &config.AppConfig{
		DefaultLimit: 20,
	}

	controller := NewMessageController(mockService, mockTimelineService, mockConfig)

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

	response := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var messagesResponse []*model.Message
	err := json.Unmarshal(response.Body.Bytes(), &messagesResponse)
	assert.NoError(t, err)
	assert.Len(t, messagesResponse, 2)
	assert.Equal(t, "msg1", messagesResponse[0].ID)

	mockService.AssertExpectations(t)
}

func TestGetUserMessages_MissingUserID(t *testing.T) {
	mockService := &MockMessageService{}
	mockTimelineService := &MockTimelineService{}
	mockConfig := &config.AppConfig{
		DefaultLimit: 20,
	}

	controller := NewMessageController(mockService, mockTimelineService, mockConfig)

	req := httptest.NewRequest("GET", "/messages", nil)

	response := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "User ID required")
}
