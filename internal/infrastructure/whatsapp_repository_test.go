package infrastructure

import (
	"anyzzapp/internal/config"
	"anyzzapp/internal/infrastructure/entity"
	"anyzzapp/pkg/domain"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHttpClient is a mock implementation of HttpClient
type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Post(payload interface{}, url string) (*http.Response, error) {
	args := m.Called(payload, url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNewWhatsAppRepository(t *testing.T) {
	cfg := config.Config{
		WhatsAppAPIKey:  "test-api-key",
		WhatsAppBaseURL: "https://graph.facebook.com/v18.0",
	}
	mockClient := &MockHttpClient{}

	repo := NewWhatsAppRepository(cfg, mockClient)

	assert.NotNil(t, repo)
	assert.IsType(t, &WhatsAppRepository{}, repo)
}

func TestWhatsAppRepository_SendMessage_Success(t *testing.T) {
	cfg := config.Config{
		WhatsAppAPIKey:  "test-api-key",
		WhatsAppBaseURL: "https://graph.facebook.com/v18.0",
	}
	mockClient := &MockHttpClient{}
	repo := &WhatsAppRepository{
		apiKey:  cfg.WhatsAppAPIKey,
		baseURL: cfg.WhatsAppBaseURL,
		client:  mockClient,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		MessageType:   "text",
	}

	expectedPayload := entity.SendWhatsAppMessagePayload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               message.To,
		Type:             "text",
		Text: &entity.Text{
			PreviewURL: false,
			Body:       message.Content,
		},
	}

	responseBody := entity.Result{
		Messages: []struct {
			ID string `json:"id"`
		}{
			{ID: "msg_123"},
		},
	}
	responseJSON, _ := json.Marshal(responseBody)
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseJSON)),
	}

	expectedURL := "https://graph.facebook.com/v18.0/123456789/messages"
	mockClient.On("Post", expectedPayload, expectedURL).Return(mockResponse, nil)

	result, err := repo.SendMessage(message)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "msg_123", result.MessageID)
	assert.Equal(t, "sent", result.Status)
	mockClient.AssertExpectations(t)
}

func TestWhatsAppRepository_SendMessage_DefaultMessageType(t *testing.T) {
	cfg := config.Config{
		WhatsAppAPIKey:  "test-api-key",
		WhatsAppBaseURL: "https://graph.facebook.com/v18.0",
	}
	mockClient := &MockHttpClient{}
	repo := &WhatsAppRepository{
		apiKey:  cfg.WhatsAppAPIKey,
		baseURL: cfg.WhatsAppBaseURL,
		client:  mockClient,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		// MessageType is empty, should default to "text"
	}

	expectedPayload := entity.SendWhatsAppMessagePayload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               message.To,
		Type:             "text", // Should be defaulted to text
		Text: &entity.Text{
			PreviewURL: false,
			Body:       message.Content,
		},
	}

	responseBody := entity.Result{
		Messages: []struct {
			ID string `json:"id"`
		}{
			{ID: "msg_123"},
		},
	}
	responseJSON, _ := json.Marshal(responseBody)
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseJSON)),
	}

	expectedURL := "https://graph.facebook.com/v18.0/123456789/messages"
	mockClient.On("Post", expectedPayload, expectedURL).Return(mockResponse, nil)

	result, err := repo.SendMessage(message)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "msg_123", result.MessageID)
	assert.Equal(t, "sent", result.Status)
	mockClient.AssertExpectations(t)
}

func TestWhatsAppRepository_SendMessage_HttpClientError(t *testing.T) {
	cfg := config.Config{
		WhatsAppAPIKey:  "test-api-key",
		WhatsAppBaseURL: "https://graph.facebook.com/v18.0",
	}
	mockClient := &MockHttpClient{}
	repo := &WhatsAppRepository{
		apiKey:  cfg.WhatsAppAPIKey,
		baseURL: cfg.WhatsAppBaseURL,
		client:  mockClient,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		MessageType:   "text",
	}

	expectedError := errors.New("network error")
	mockClient.On("Post", mock.Anything, mock.Anything).Return((*http.Response)(nil), expectedError)

	result, err := repo.SendMessage(message)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockClient.AssertExpectations(t)
}

func TestWhatsAppRepository_SendMessage_APIError(t *testing.T) {
	cfg := config.Config{
		WhatsAppAPIKey:  "test-api-key",
		WhatsAppBaseURL: "https://graph.facebook.com/v18.0",
	}
	mockClient := &MockHttpClient{}
	repo := &WhatsAppRepository{
		apiKey:  cfg.WhatsAppAPIKey,
		baseURL: cfg.WhatsAppBaseURL,
		client:  mockClient,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		MessageType:   "text",
	}

	responseBody := entity.Result{
		Error: struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "Invalid phone number",
			Code:    400,
		},
	}
	responseJSON, _ := json.Marshal(responseBody)
	mockResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewReader(responseJSON)),
	}

	mockClient.On("Post", mock.Anything, mock.Anything).Return(mockResponse, nil)

	result, err := repo.SendMessage(message)

	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "failed", result.Status)
	assert.Equal(t, "Invalid phone number", result.Message)
	assert.Contains(t, err.Error(), "API error")
	mockClient.AssertExpectations(t)
}

func TestWhatsAppRepository_SendMessage_NoMessageID(t *testing.T) {
	cfg := config.Config{
		WhatsAppAPIKey:  "test-api-key",
		WhatsAppBaseURL: "https://graph.facebook.com/v18.0",
	}
	mockClient := &MockHttpClient{}
	repo := &WhatsAppRepository{
		apiKey:  cfg.WhatsAppAPIKey,
		baseURL: cfg.WhatsAppBaseURL,
		client:  mockClient,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		MessageType:   "text",
	}

	responseBody := entity.Result{
		Messages: []struct {
			ID string `json:"id"`
		}{}, // Empty messages array
	}
	responseJSON, _ := json.Marshal(responseBody)
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseJSON)),
	}

	mockClient.On("Post", mock.Anything, mock.Anything).Return(mockResponse, nil)

	result, err := repo.SendMessage(message)

	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "failed", result.Status)
	assert.Contains(t, err.Error(), "no message ID returned from API")
	mockClient.AssertExpectations(t)
}

func TestWhatsAppRepository_MarkAsRead_Success(t *testing.T) {
	cfg := config.Config{
		WhatsAppAPIKey:  "test-api-key",
		WhatsAppBaseURL: "https://graph.facebook.com/v18.0",
	}
	mockClient := &MockHttpClient{}
	repo := &WhatsAppRepository{
		apiKey:  cfg.WhatsAppAPIKey,
		baseURL: cfg.WhatsAppBaseURL,
		client:  mockClient,
	}

	phoneNumberID := "123456789"
	messageID := "msg_123"

	expectedPayload := markAsReadPayload{
		MessagingProduct: "whatsapp",
		Status:           "read",
		MessageID:        messageID,
	}

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
	}

	expectedURL := "https://graph.facebook.com/v18.0/123456789/messages"
	mockClient.On("Post", expectedPayload, expectedURL).Return(mockResponse, nil)

	err := repo.MarkAsRead(phoneNumberID, messageID)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWhatsAppRepository_MarkAsRead_HttpClientError(t *testing.T) {
	cfg := config.Config{
		WhatsAppAPIKey:  "test-api-key",
		WhatsAppBaseURL: "https://graph.facebook.com/v18.0",
	}
	mockClient := &MockHttpClient{}
	repo := &WhatsAppRepository{
		apiKey:  cfg.WhatsAppAPIKey,
		baseURL: cfg.WhatsAppBaseURL,
		client:  mockClient,
	}

	phoneNumberID := "123456789"
	messageID := "msg_123"

	expectedError := errors.New("network error")
	mockClient.On("Post", mock.Anything, mock.Anything).Return((*http.Response)(nil), expectedError)

	err := repo.MarkAsRead(phoneNumberID, messageID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockClient.AssertExpectations(t)
}

func TestWhatsAppRepository_MarkAsRead_APIError(t *testing.T) {
	cfg := config.Config{
		WhatsAppAPIKey:  "test-api-key",
		WhatsAppBaseURL: "https://graph.facebook.com/v18.0",
	}
	mockClient := &MockHttpClient{}
	repo := &WhatsAppRepository{
		apiKey:  cfg.WhatsAppAPIKey,
		baseURL: cfg.WhatsAppBaseURL,
		client:  mockClient,
	}

	phoneNumberID := "123456789"
	messageID := "msg_123"

	mockResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
	}

	mockClient.On("Post", mock.Anything, mock.Anything).Return(mockResponse, nil)

	err := repo.MarkAsRead(phoneNumberID, messageID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
	mockClient.AssertExpectations(t)
}
