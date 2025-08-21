package server

import (
	"anyzzapp/internal/config"
	"anyzzapp/internal/interfaces/http"
	"anyzzapp/pkg/domain"
	"log"
)

func Run(config config.Config, whatsAppUsecase domain.WhatsAppUseCaseInterface) {
	// Initialize HTTP handlers
	router := http.NewRouter(config, whatsAppUsecase)

	// Start server
	log.Printf("Server starting on port %s", config.ServerPort)
	if err := router.Run(":" + config.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
