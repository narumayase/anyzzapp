package main

import (
	"anyzzapp/cmd/server"
	"anyzzapp/internal/config"
	"anyzzapp/internal/infrastructure"
	"anyzzapp/pkg/application"
	"github.com/rs/zerolog"
)

func main() {
	// Load configuration
	cfg := config.Load()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Initialize repository layers
	whatsappRepo := infrastructure.NewWhatsAppRepository(*cfg)
	llmRepo := infrastructure.NewLLMRepository(*cfg)

	server.Run(*cfg, application.NewWhatsAppUseCase(whatsappRepo, llmRepo))
}
