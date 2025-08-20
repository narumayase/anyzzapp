package http

import (
	"anyzzapp/internal/config"
	"anyzzapp/internal/interfaces/http/handler"
	"anyzzapp/internal/interfaces/http/middleware"
	"anyzzapp/pkg/domain"
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the HTTP router
func NewRouter(config config.Config,
	whatsappUseCase domain.WhatsAppUseCaseInterface) *gin.Engine {
	router := gin.Default()

	// Add middlewares
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// Initialize handlers
	whatsappHandler := handler.NewWhatsAppHandler(config, whatsappUseCase)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "anyzzapp API is running",
		})
	})

	v1 := router.Group("/api/v1")

	whatsapp := v1.Group("/whatsapp")
	whatsapp.POST("/send", whatsappHandler.SendMessage)
	whatsapp.POST("/webhook", whatsappHandler.ReceiveWebhook)
	whatsapp.GET("/webhook", whatsappHandler.VerifyWebhook)

	return router
}
