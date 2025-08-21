package main

import (
	"anyzzapp/cmd/server"
	"anyzzapp/internal/config"
	"anyzzapp/internal/infrastructure"
	"anyzzapp/pkg/application"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.Load()

	llmClient := &http.Client{}
	llmHttpClient := infrastructure.NewHttpClient(llmClient, cfg.LLMBearerToken)

	whatsAppClient := &http.Client{}
	whatsappHttpClient := infrastructure.NewHttpClient(whatsAppClient, cfg.WhatsAppAPIKey)

	// Initialize repository layers
	whatsappRepo := infrastructure.NewWhatsAppRepository(cfg, whatsappHttpClient)
	llmRepo := infrastructure.NewLLMRepository(cfg, llmHttpClient)

	server.Run(cfg, application.NewWhatsAppUseCase(whatsappRepo, llmRepo))
}
