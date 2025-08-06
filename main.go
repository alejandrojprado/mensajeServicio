package main

import (
	"context"
	"mensajesService/components/config"
	"mensajesService/components/database"
	"mensajesService/components/logger"
	"mensajesService/components/metrics"
	"mensajesService/message-api/controller"
	"mensajesService/message-api/service"
	"mensajesService/message-api/web"
	"net/http"
	"os"
)

func main() {
	logger.Init()
	ctx := context.Background()
	cfg := config.LoadConfig()
	metrics.Init(cfg.Region)

	dbClient, err := database.NewDDBClient(ctx, cfg)
	if err != nil {
		logger.LogError("Error initializing database", "error", err)
		os.Exit(1)
	}

	messageService := service.NewMessageService(dbClient)
	followService := service.NewFollowService(dbClient)
	timelineService := service.NewTimelineService(dbClient)

	messageController := controller.NewMessageController(messageService, timelineService, cfg)
	followController := controller.NewFollowController(followService, cfg)
	timelineController := controller.NewTimelineController(timelineService, cfg)

	router := web.NewHttpHandler("v1")

	messageController.MountIn(router)
	followController.MountIn(router)
	timelineController.MountIn(router)

	port := cfg.Port
	logger.LogInfo("Service started on port: " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.LogError("Error starting service: ", "error", err)
	}
}
