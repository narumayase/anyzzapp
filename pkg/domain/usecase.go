package domain

// WhatsAppUseCaseInterface defines the business logic operations
type WhatsAppUseCaseInterface interface {
	SendMessage(phoneNumberID, to, content, messageType string) (*SendMessageResponse, error)
	ProcessIncomingWebhook(webhook *WebhookRequest) error
}
