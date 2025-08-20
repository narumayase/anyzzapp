package server

import (
	"anyzzapp/internal/config"
	"anyzzapp/internal/infrastructure"
	"anyzzapp/internal/interfaces/http"
	"anyzzapp/pkg/application"
	"log"
)

func Run() {
	// Load configuration
	cfg := config.Load()

	// Initialize repository layers
	whatsappRepo := infrastructure.NewWhatsAppRepository(*cfg)
	llmRepo := infrastructure.NewLLMRepository(*cfg)

	// Initialize use case layers
	whatsappUseCase := application.NewWhatsAppUseCase(whatsappRepo, llmRepo)

	// Initialize HTTP handlers
	router := http.NewRouter(*cfg, whatsappUseCase)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
