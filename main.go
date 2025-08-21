package main

import (
	"anyzzapp/cmd/server"
	"anyzzapp/internal/config"
	"anyzzapp/internal/infrastructure"
	"anyzzapp/internal/infrastructure/client"
	"anyzzapp/pkg/application"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.Load()

	llmClient := &http.Client{}
	llmHttpClient := client.NewHttpClient(llmClient, cfg.LLMBearerToken)

	whatsAppClient := &http.Client{}
	whatsappHttpClient := client.NewHttpClient(whatsAppClient, cfg.WhatsAppAPIKey)

	// Initialize repository layers
	whatsappRepo := infrastructure.NewWhatsAppRepository(cfg, whatsappHttpClient)
	llmRepo := infrastructure.NewLLMRepository(cfg, llmHttpClient)

	server.Run(cfg, application.NewWhatsAppUseCase(whatsappRepo, llmRepo))
}
