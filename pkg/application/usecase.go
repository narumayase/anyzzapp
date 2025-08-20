package application

import (
	"anyzzapp/pkg/domain"
	"fmt"
)

// WhatsAppUseCase implements WhatsAppUseCaseInterface
type WhatsAppUseCase struct {
	whatsappRepo domain.WhatsAppRepository
	llmRepo      domain.LLMRepository
}

// NewWhatsAppUseCase creates a new instance of WhatsAppUseCase
func NewWhatsAppUseCase(whatsappRepo domain.WhatsAppRepository,
	llmRepo domain.LLMRepository) domain.WhatsAppUseCaseInterface {
	return &WhatsAppUseCase{
		whatsappRepo: whatsappRepo,
		llmRepo:      llmRepo,
	}
}

// SendMessage handles the business logic for sending a message
func (uc *WhatsAppUseCase) SendMessage(phoneNumberID, to, content, messageType string) (*domain.SendMessageResponse, error) {
	// Validate input
	if phoneNumberID == "" {
		return nil, fmt.Errorf("phone number ID is required")
	}
	if to == "" {
		return nil, fmt.Errorf("recipient phone number is required")
	}
	if content == "" {
		return nil, fmt.Errorf("message content is required")
	}

	// Send message through WhatsApp API
	response, err := uc.whatsappRepo.SendMessage(phoneNumberID, to, content, messageType)
	if err != nil {
		return response, fmt.Errorf("failed to send message: %w", err)
	}
	return response, nil
}

// ProcessIncomingWebhook processes incoming webhook data from WhatsApp
func (uc *WhatsAppUseCase) ProcessIncomingWebhook(webhook *domain.WebhookRequest) error {
	if webhook == nil {
		return fmt.Errorf("webhook data cannot be nil")
	}

	// Process each entry in the webhook
	for _, entry := range webhook.Entry {
		for _, change := range entry.Changes {
			// Process incoming messages
			if err := uc.processMessages(change.Value.Messages, change.Value.Metadata.PhoneNumberID); err != nil {
				return fmt.Errorf("failed to process messages: %w", err)
			}
		}
	}
	return nil
}

// processMessages handles incoming messages
func (uc *WhatsAppUseCase) processMessages(messages []domain.WebhookMessage, phoneNumberID string) error {
	for _, msg := range messages {

		// Extract message content based on type
		var content string
		var err error

		if msg.Text != nil {
			content = msg.Text.Body
		}
		// Future: handle other message types (image, audio, etc.)TODO

		// Mark message as read
		if err = uc.whatsappRepo.MarkAsRead(phoneNumberID, msg.ID); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to mark message as read: %v\n", err)
		}

		// Auto-reply
		if content != "" && msg.Type == "text" {
			replyMessage := ""

			// Send the question to LLM
			if replyMessage, err = uc.llmRepo.SendMessage(content); err != nil {
				fmt.Printf("Warning: failed to send message: %v\n", err)
				//TODO por ahora no manejamos el error para poder devolver una respuesta "mockeada"
			} else {
				replyMessage = fmt.Sprintf("Me dijiste: %s", content)
			}
			// Send the reply
			_, err = uc.whatsappRepo.SendMessage(phoneNumberID, msg.From, replyMessage, "text")
			if err != nil {
				fmt.Printf("Warning: failed to send auto-reply: %v\n", err)
				return err
			}
		}
	}
	return nil
}
