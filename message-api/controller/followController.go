package controller

import (
	"encoding/json"
	"net/http"

	"mensajesService/components/config"
	"mensajesService/components/metrics"
	"mensajesService/message-api/model"
	"mensajesService/message-api/service"

	"github.com/go-chi/chi/v5"
)

type FollowController struct {
	messageService service.MessageServiceInterface
	config         *config.AppConfig
}

func NewFollowController(messageService service.MessageServiceInterface, cfg *config.AppConfig) *FollowController {
	return &FollowController{
		messageService: messageService,
		config:         cfg,
	}
}

func (c *FollowController) MountIn(r chi.Router) {
	r.Route("/follows", func(r chi.Router) {
		r.Post("/", c.FollowUser)
	})
}

func (c *FollowController) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		metrics.PutCountMetric(metrics.MetricValidationError, 1)
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	var followRequest model.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&followRequest); err != nil {
		metrics.PutCountMetric(metrics.MetricValidationError, 1)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if followRequest.FollowingID == "" {
		metrics.PutCountMetric(metrics.MetricValidationError, 1)
		http.Error(w, "Following ID is required", http.StatusBadRequest)
		return
	}

	if userID == followRequest.FollowingID {
		metrics.PutCountMetric(metrics.MetricValidationError, 1)
		http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
		return
	}

	// TODO: Implementar la lógica de follow en el service
	// err := c.messageService.FollowUser(r.Context(), userID, followRequest.FollowingID)
	// if err != nil {
	//     metrics.PutCountMetric(metrics.MetricFollowError, 1)
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     logger.LogError("FollowUser error", "error", err, "user_id", userID, "following_id", followRequest.FollowingID)
	//     return
	// }

	metrics.PutCountMetric(metrics.MetricFollowSuccess, 1)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User followed successfully",
	})
}
