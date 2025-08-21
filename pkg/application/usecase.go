package application

import (
	"anyzzapp/pkg/domain"
	"fmt"

	"github.com/rs/zerolog/log"
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
func (uc *WhatsAppUseCase) SendMessage(message domain.Message) (*domain.SendMessageResponse, error) {
	// Validate input
	if message.PhoneNumberID == "" {
		return nil, fmt.Errorf("phone number ID is required")
	}
	if message.To == "" {
		return nil, fmt.Errorf("recipient phone number is required")
	}
	if message.Content == "" {
		return nil, fmt.Errorf("message content is required")
	}
	// Send message through WhatsApp API
	response, err := uc.whatsappRepo.SendMessage(message)
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
			log.Debug().Msgf("message received: %v - metadata phone number id: %s", change.Value.Messages, change.Value.Metadata.PhoneNumberID)

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
			log.Warn().Msgf("failed to mark message as read: %v\n", err)
		}
		// Auto-reply
		if content != "" && msg.Type == "text" {
			replyMessage := ""
			// Send the question to LLM
			if replyMessage, err = uc.llmRepo.SendMessage(content); err != nil {
				log.Err(fmt.Errorf("failed to send message: %v\n", err))
				return err
			}
			// Send the reply
			if _, err = uc.whatsappRepo.SendMessage(domain.Message{
				PhoneNumberID: phoneNumberID,
				To:            removeNine(msg.From),
				Content:       replyMessage,
				MessageType:   "text",
			}); err != nil {
				log.Err(fmt.Errorf("failed to send auto-reply: %v\n", err))
				return err
			}
		}
	}
	return nil
}

// removeNine for Argentinian numbers it's necessary to remove the 9 from the reception phone number to send messages to it.
func removeNine(phoneNumber string) string {
	// example: phoneNumber = "5491112345678"
	if len(phoneNumber) == 13 && phoneNumber[:2] == "54" && phoneNumber[2] == '9' {
		return phoneNumber[:2] + phoneNumber[3:]
	}
	// result: "541112345678"
	return phoneNumber
}
