package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mensajesService/components/config"
	"mensajesService/components/metrics"
	"mensajesService/message-api/service"

	"github.com/go-chi/chi/v5"
)

type TimelineController struct {
	messageService service.MessageServiceInterface
	config         *config.AppConfig
}

func NewTimelineController(messageService service.MessageServiceInterface, cfg *config.AppConfig) *TimelineController {
	return &TimelineController{
		messageService: messageService,
		config:         cfg,
	}
}

func (c *TimelineController) MountIn(r chi.Router) {
	r.Route("/timeline", func(r chi.Router) {
		r.Get("/", c.GetTimeline)
	})
}

func (c *TimelineController) GetTimeline(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		metrics.PutCountMetric(metrics.MetricTimelineError, 1)
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

	timeline, err := c.messageService.GetUserTimeline(r.Context(), userID, limit)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricTimelineError, 1)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	metrics.PutCountMetric(metrics.MetricTimelineSuccess, 1)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(timeline)
}
