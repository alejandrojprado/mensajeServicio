package controller

import (
	"encoding/json"
	"net/http"

	"mensajesService/components/config"
	"mensajesService/components/logger"
	"mensajesService/components/metrics"
	"mensajesService/message-api/model"
	"mensajesService/message-api/service"

	"github.com/go-chi/chi/v5"
)

type FollowController struct {
	followService service.FollowServiceInterface
	config        *config.AppConfig
}

func NewFollowController(followService service.FollowServiceInterface, cfg *config.AppConfig) *FollowController {
	return &FollowController{
		followService: followService,
		config:        cfg,
	}
}

func (c *FollowController) MountIn(r chi.Router) {
	r.Route("/follow", func(r chi.Router) {
		r.Post("/", c.FollowUser)
	})
}

func (c *FollowController) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		metrics.PutCountMetric(metrics.MetricFollowError, 1)
		logger.LogError("FollowUser error", "error", "User ID required in X-User-ID header", "user_id", userID)
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	var followRequest model.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&followRequest); err != nil {
		metrics.PutCountMetric(metrics.MetricFollowError, 1)
		logger.LogError("FollowUser error", "error", "Invalid request body", "user_id", userID)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if followRequest.FollowingID == "" {
		metrics.PutCountMetric(metrics.MetricFollowError, 1)
		logger.LogError("FollowUser error", "error", "Following ID is required", "user_id", userID)
		http.Error(w, "Following ID is required", http.StatusBadRequest)
		return
	}

	if userID == followRequest.FollowingID {
		metrics.PutCountMetric(metrics.MetricFollowError, 1)
		logger.LogError("FollowUser error", "error", "Cannot follow yourself", "user_id", userID, "following_id", followRequest.FollowingID)
		http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
		return
	}

	err := c.followService.FollowUser(r.Context(), userID, followRequest.FollowingID)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricFollowError, 1)
		logger.LogError("FollowUser error", "error", err, "user_id", userID, "following_id", followRequest.FollowingID)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	metrics.PutCountMetric(metrics.MetricFollowSuccess, 1)
	logger.LogInfo("FollowUser success", "user_id", userID, "following_id", followRequest.FollowingID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User followed successfully",
	})
}
