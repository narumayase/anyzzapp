package domain

// WhatsAppUseCaseInterface defines the business logic operations
type WhatsAppUseCaseInterface interface {
	SendMessage(message Message) (*SendMessageResponse, error)
	ProcessIncomingWebhook(webhook *WebhookRequest) error
}
