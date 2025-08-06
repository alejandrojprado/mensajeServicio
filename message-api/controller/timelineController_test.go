package controller

import (
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

type MockTimelineService struct {
	mock.Mock
}

var _ service.TimelineServiceInterface = (*MockTimelineService)(nil)

func (m *MockTimelineService) GetUserTimeline(ctx context.Context, userID string, limit int) ([]*model.TimelineItem, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]*model.TimelineItem), args.Error(1)
}

func TestGetTimeline_Success(t *testing.T) {
	logger.Init()

	mockService := &MockTimelineService{}
	mockConfig := &config.AppConfig{DefaultLimit: 10}

	expectedTimeline := []*model.TimelineItem{
		{
			MessageID: "msg1",
			UserID:    "user123",
			AuthorID:  "author1",
			Content:   "Test message 1",
			CreatedAt: time.Now(),
		},
		{
			MessageID: "msg2",
			UserID:    "user123",
			AuthorID:  "author2",
			Content:   "Test message 2",
			CreatedAt: time.Now(),
		},
	}

	mockService.On("GetUserTimeline", mock.Anything, "user123", 10).Return(expectedTimeline, nil)

	controller := NewTimelineController(mockService, mockConfig)

	req := httptest.NewRequest("GET", "/timeline", nil)
	req.Header.Set("X-User-ID", "user123")

	response := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var timelineResponse []*model.TimelineItem
	err := json.Unmarshal(response.Body.Bytes(), &timelineResponse)
	assert.NoError(t, err)
	assert.Len(t, timelineResponse, 2)
	assert.Equal(t, expectedTimeline[0].MessageID, timelineResponse[0].MessageID)
	assert.Equal(t, expectedTimeline[1].MessageID, timelineResponse[1].MessageID)

	mockService.AssertExpectations(t)
}
