package main

import (
	"context"
	"mensajesService/components/config"
	"mensajesService/components/database"
	"mensajesService/components/logger"
	"mensajesService/components/metrics"
	messageController "mensajesService/message-api/controller"
	messageService "mensajesService/message-api/service"
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
		logger.LogError("No se pudo inicializar la base de datos", "error", err)
		os.Exit(1)
	}

	messageService := messageService.NewMessageService(dbClient)
	messageController := messageController.NewMessageController(messageService, cfg)

	router := web.NewHttpHandler("v1")
	messageController.MountIn(router)

	port := cfg.Port
	logger.LogInfo("Iniciando servicio en el puerto: " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.LogError("Error al iniciar el servicio: ", "error", err)
	}
}
