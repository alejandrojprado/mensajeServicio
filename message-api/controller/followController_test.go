package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mensajesService/components/config"
	"mensajesService/message-api/model"
	"mensajesService/message-api/service"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFollowService struct {
	mock.Mock
}

var _ service.FollowServiceInterface = (*MockFollowService)(nil)

func (m *MockFollowService) FollowUser(ctx context.Context, userID, followingID string) error {
	args := m.Called(ctx, userID, followingID)
	return args.Error(0)
}

func TestFollowUser_Success(t *testing.T) {
	mockService := &MockFollowService{}
	mockConfig := &config.AppConfig{}

	mockService.On("FollowUser", mock.Anything, "user123", "user456").Return(nil)

	controller := NewFollowController(mockService, mockConfig)

	followRequest := model.FollowRequest{
		FollowingID: "user456",
	}

	body, _ := json.Marshal(followRequest)
	req := httptest.NewRequest("POST", "/follows", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "user123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	controller.MountIn(router)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User followed successfully", response["message"])

	mockService.AssertExpectations(t)
}
