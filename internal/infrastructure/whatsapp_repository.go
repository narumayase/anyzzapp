package infrastructure

import (
	"anyzzapp/internal/config"
	"anyzzapp/internal/infrastructure/entity"
	"anyzzapp/pkg/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WhatsAppRepository implements WhatsAppRepository
type WhatsAppRepository struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewWhatsAppRepository creates a new instance of WhatsAppRepository
func NewWhatsAppRepository(config config.Config) domain.WhatsAppRepository {
	return &WhatsAppRepository{
		apiKey:  config.WhatsAppAPIKey,
		baseURL: config.WhatsAppBaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage sends a message through WhatsApp API
func (r *WhatsAppRepository) SendMessage(phoneNumberID, to, content, messageType string) (*domain.SendMessageResponse, error) {
	// Default message type to text if not specified
	if messageType == "" {
		messageType = "text"
	}
	// TODO add more types?

	// Prepare the payload
	payload := entity.SendMessagePayload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             messageType,
	}

	// Currently only supporting text messages //TODO
	if messageType == "text" {
		payload.Text = &entity.Text{
			PreviewURL: false,
			Body:       content,
		}
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/%s/messages", r.baseURL, phoneNumberID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Parse entity
	result := entity.Result{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode entity: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return &domain.SendMessageResponse{
			Status:  "failed",
			Message: result.Error.Message,
		}, fmt.Errorf("API error: %s (code: %d)", result.Error.Message, result.Error.Code)
	}

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

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", r.baseURL, phoneNumberID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	return nil
}
