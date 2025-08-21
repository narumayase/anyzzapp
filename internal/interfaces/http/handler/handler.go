package handler

import (
	"anyzzapp/internal/config"
	"anyzzapp/pkg/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

// WhatsAppHandler handles HTTP requests related to WhatsApp operations
type WhatsAppHandler struct {
	whatsappUseCase domain.WhatsAppUseCaseInterface
	config          config.Config
}

// NewWhatsAppHandler creates a new instance of WhatsAppHandler
func NewWhatsAppHandler(
	config config.Config,
	whatsappUseCase domain.WhatsAppUseCaseInterface) *WhatsAppHandler {
	return &WhatsAppHandler{
		whatsappUseCase: whatsappUseCase,
		config:          config,
	}
}

// SendMessage handles POST /api/v1/whatsapp/send
func (h *WhatsAppHandler) SendMessage(c *gin.Context) {
	var req domain.Message

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Default message type to text if not provided
	if req.MessageType == "" {
		req.MessageType = "text"
	}

	// Call use case
	response, err := h.whatsappUseCase.SendMessage(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "send_failed",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ReceiveWebhook handles POST /api/v1/whatsapp/webhook
func (h *WhatsAppHandler) ReceiveWebhook(c *gin.Context) {
	var webhook domain.WebhookRequest

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "invalid_webhook",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Process the webhook
	if err := h.whatsappUseCase.ProcessIncomingWebhook(&webhook); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "webhook_processing_failed",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// WhatsApp expects a 200 status code
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// VerifyWebhook handles GET /api/v1/whatsapp/webhook for webhook verification
func (h *WhatsAppHandler) VerifyWebhook(c *gin.Context) {
	// Get verification parameters
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	// Verify token
	expectedToken := h.config.WebhookVerifyToken

	if mode == "subscribe" && token == expectedToken {
		// Return the challenge to verify the webhook
		c.String(http.StatusOK, challenge)
		return
	}

	c.JSON(http.StatusForbidden, domain.ErrorResponse{
		Error:   "verification_failed",
		Message: "Webhook verification failed",
		Code:    http.StatusForbidden,
	})
}
