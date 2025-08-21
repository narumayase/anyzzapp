package infrastructure

import (
	"anyzzapp/internal/config"
	client2 "anyzzapp/internal/infrastructure/client"
	"anyzzapp/internal/infrastructure/entity"
	"anyzzapp/pkg/domain"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// WhatsAppRepository implements WhatsAppRepository
type WhatsAppRepository struct {
	apiKey  string
	baseURL string
	client  client2.HttpClient
}

// NewWhatsAppRepository creates a new instance of WhatsAppRepository
func NewWhatsAppRepository(config config.Config, client client2.HttpClient) domain.WhatsAppRepository {
	return &WhatsAppRepository{
		apiKey:  config.WhatsAppAPIKey,
		baseURL: config.WhatsAppBaseURL,
		client:  client,
	}
}

// SendMessage sends a message through WhatsApp API
func (r *WhatsAppRepository) SendMessage(message domain.Message) (*domain.SendMessageResponse, error) {
	// Default message type to text if not specified
	if message.MessageType == "" {
		message.MessageType = "text"
	}
	// TODO add more types - add audio type
	// Prepare the payload
	payload := entity.SendWhatsAppMessagePayload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               message.To,
		Type:             message.MessageType,
	}
	// Currently only supporting text messages //TODO - add audio type
	if message.MessageType == "text" {
		payload.Text = &entity.Text{
			PreviewURL: false,
			Body:       message.Content,
		}
	}
	url := fmt.Sprintf("%s/%s/messages", r.baseURL, message.PhoneNumberID)
	// Execute POST
	resp, err := r.client.Post(payload, url)
	if err != nil {
		return nil, err
	}
	// Parse entity
	result := entity.Result{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode entity: %w", err)
	}
	log.Debug().Msgf("whatsapp response: %v", resp)

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return &domain.SendMessageResponse{
			Status:  "failed",
			Message: result.Error.Message,
		}, fmt.Errorf("API error: %s (code: %d)", result.Error.Message, result.Error.Code)
	}
	defer resp.Body.Close()

	// Check if we got a message ID
	if len(result.Messages) == 0 {
		return &domain.SendMessageResponse{
			Status:  "failed",
			Message: "No message ID returned from API",
		}, fmt.Errorf("no message ID returned from API")
	}
	return &domain.SendMessageResponse{
		MessageID: result.Messages[0].ID,
		Status:    "sent",
	}, nil
}

// markAsReadPayload represents the payload for marking messages as read
type markAsReadPayload struct {
	MessagingProduct string `json:"messaging_product"`
	Status           string `json:"status"`
	MessageID        string `json:"message_id"`
}

// MarkAsRead marks a message as read
func (r *WhatsAppRepository) MarkAsRead(phoneNumberID, messageID string) error {
	payload := markAsReadPayload{
		MessagingProduct: "whatsapp",
		Status:           "read",
		MessageID:        messageID,
	}
	url := fmt.Sprintf("%s/%s/messages", r.baseURL, phoneNumberID)
	// Execute POST
	resp, err := r.client.Post(payload, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: status code %d", resp.StatusCode)
	}
	return nil
}
