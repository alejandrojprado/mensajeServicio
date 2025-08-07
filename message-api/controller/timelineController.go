package controller

import (
	"encoding/json"
	"net/http"

	"mensajesService/components/config"
	"mensajesService/components/logger"
	"mensajesService/components/metrics"
	"mensajesService/message-api/service"

	"github.com/go-chi/chi/v5"
)

type TimelineController struct {
	timelineService service.TimelineServiceInterface
	config          *config.AppConfig
}

func NewTimelineController(timelineService service.TimelineServiceInterface, cfg *config.AppConfig) *TimelineController {
	return &TimelineController{
		timelineService: timelineService,
		config:          cfg,
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
		logger.LogError("GetTimeline error", "error", "User ID required in X-User-ID header", "user_id", userID)
		http.Error(w, "User ID required in X-User-ID header", http.StatusBadRequest)
		return
	}

	timeline, err := c.timelineService.GetUserTimeline(r.Context(), userID, c.config.DefaultLimit)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricTimelineError, 1)
		logger.LogError("GetTimeline error", "error", err, "user_id", userID)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(timeline) == 0 {
		metrics.PutCountMetric(metrics.MetricTimelineError, 1)
		logger.LogError("Get Timeline error", "error", "User not found", "user_id", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	timelineMessagesAmount := float64(len(timeline))
	metrics.PutCountMetric(metrics.MetricTimelineSuccess, 1)
	metrics.PutCountMetric(metrics.MetricTimelineCount, timelineMessagesAmount)
	logger.LogInfo("GetTimeline success", "user_id", userID, "timeline_count", timelineMessagesAmount)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(timeline)
}
