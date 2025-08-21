package handler

import (
	"anyzzapp/internal/config"
	"anyzzapp/pkg/domain"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWhatsAppUseCase is a mock implementation of WhatsAppUseCaseInterface
type MockWhatsAppUseCase struct {
	mock.Mock
}

func (m *MockWhatsAppUseCase) SendMessage(message domain.Message) (*domain.SendMessageResponse, error) {
	args := m.Called(message)
	return args.Get(0).(*domain.SendMessageResponse), args.Error(1)
}

func (m *MockWhatsAppUseCase) ProcessIncomingWebhook(webhook *domain.WebhookRequest) error {
	args := m.Called(webhook)
	return args.Error(0)
}

func TestNewWhatsAppHandler(t *testing.T) {
	cfg := config.Config{
		WebhookVerifyToken: "test-verify-token",
	}
	mockUseCase := &MockWhatsAppUseCase{}

	handler := NewWhatsAppHandler(cfg, mockUseCase)

	assert.NotNil(t, handler)
	assert.IsType(t, &WhatsAppHandler{}, handler)
}

func TestWhatsAppHandler_SendMessage_Success(t *testing.T) {
	cfg := config.Config{}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		MessageType:   "text",
	}

	expectedResponse := &domain.SendMessageResponse{
		MessageID: "msg_123",
		Status:    "sent",
	}

	mockUseCase.On("SendMessage", message).Return(expectedResponse, nil)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request body
	jsonBody, _ := json.Marshal(message)
	c.Request = httptest.NewRequest("POST", "/api/v1/whatsapp/send", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendMessage(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.SendMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.MessageID, response.MessageID)
	assert.Equal(t, expectedResponse.Status, response.Status)
	mockUseCase.AssertExpectations(t)
}

func TestWhatsAppHandler_SendMessage_InvalidJSON(t *testing.T) {
	cfg := config.Config{}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Invalid JSON
	c.Request = httptest.NewRequest("POST", "/api/v1/whatsapp/send", bytes.NewReader([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendMessage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_request", errorResponse.Error)
	assert.Equal(t, http.StatusBadRequest, errorResponse.Code)
}

func TestWhatsAppHandler_SendMessage_DefaultMessageType(t *testing.T) {
	cfg := config.Config{}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		// MessageType is empty, should be defaulted to "text"
	}

	expectedMessage := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		MessageType:   "text", // Should be defaulted
	}

	expectedResponse := &domain.SendMessageResponse{
		MessageID: "msg_123",
		Status:    "sent",
	}

	mockUseCase.On("SendMessage", expectedMessage).Return(expectedResponse, nil)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request body
	jsonBody, _ := json.Marshal(message)
	c.Request = httptest.NewRequest("POST", "/api/v1/whatsapp/send", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendMessage(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestWhatsAppHandler_SendMessage_UseCaseError(t *testing.T) {
	cfg := config.Config{}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		MessageType:   "text",
	}

	expectedError := errors.New("use case error")
	mockUseCase.On("SendMessage", message).Return((*domain.SendMessageResponse)(nil), expectedError)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request body
	jsonBody, _ := json.Marshal(message)
	c.Request = httptest.NewRequest("POST", "/api/v1/whatsapp/send", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendMessage(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "send_failed", errorResponse.Error)
	assert.Equal(t, http.StatusInternalServerError, errorResponse.Code)
	mockUseCase.AssertExpectations(t)
}

func TestWhatsAppHandler_ReceiveWebhook_Success(t *testing.T) {
	cfg := config.Config{}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	webhook := domain.WebhookRequest{
		Object: "whatsapp_business_account",
		Entry: []domain.WebhookEntry{
			{
				ID: "entry_123",
				Changes: []domain.WebhookChange{
					{
						Value: domain.WebhookValue{
							MessagingProduct: "whatsapp",
							Metadata: domain.WebhookMetadata{
								DisplayPhoneNumber: "5491112345678",
								PhoneNumberID:      "123456789",
							},
						},
						Field: "messages",
					},
				},
			},
		},
	}

	mockUseCase.On("ProcessIncomingWebhook", &webhook).Return(nil)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request body
	jsonBody, _ := json.Marshal(webhook)
	c.Request = httptest.NewRequest("POST", "/api/v1/whatsapp/webhook", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ReceiveWebhook(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	mockUseCase.AssertExpectations(t)
}

func TestWhatsAppHandler_ReceiveWebhook_InvalidJSON(t *testing.T) {
	cfg := config.Config{}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Invalid JSON
	c.Request = httptest.NewRequest("POST", "/api/v1/whatsapp/webhook", bytes.NewReader([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ReceiveWebhook(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "invalid_webhook", errorResponse.Error)
	assert.Equal(t, http.StatusBadRequest, errorResponse.Code)
}

func TestWhatsAppHandler_ReceiveWebhook_ProcessingError(t *testing.T) {
	cfg := config.Config{}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	webhook := domain.WebhookRequest{
		Object: "whatsapp_business_account",
		Entry:  []domain.WebhookEntry{},
	}

	expectedError := errors.New("processing error")
	mockUseCase.On("ProcessIncomingWebhook", &webhook).Return(expectedError)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request body
	jsonBody, _ := json.Marshal(webhook)
	c.Request = httptest.NewRequest("POST", "/api/v1/whatsapp/webhook", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ReceiveWebhook(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "webhook_processing_failed", errorResponse.Error)
	assert.Equal(t, http.StatusInternalServerError, errorResponse.Code)
	mockUseCase.AssertExpectations(t)
}

func TestWhatsAppHandler_VerifyWebhook_Success(t *testing.T) {
	cfg := config.Config{
		WebhookVerifyToken: "test-verify-token",
	}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Setup query parameters
	c.Request = httptest.NewRequest("GET", "/api/v1/whatsapp/webhook?hub.mode=subscribe&hub.verify_token=test-verify-token&hub.challenge=test-challenge", nil)

	handler.VerifyWebhook(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test-challenge", w.Body.String())
}

func TestWhatsAppHandler_VerifyWebhook_WrongToken(t *testing.T) {
	cfg := config.Config{
		WebhookVerifyToken: "correct-token",
	}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Setup query parameters with wrong token
	c.Request = httptest.NewRequest("GET", "/api/v1/whatsapp/webhook?hub.mode=subscribe&hub.verify_token=wrong-token&hub.challenge=test-challenge", nil)

	handler.VerifyWebhook(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var errorResponse domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "verification_failed", errorResponse.Error)
	assert.Equal(t, http.StatusForbidden, errorResponse.Code)
}

func TestWhatsAppHandler_VerifyWebhook_WrongMode(t *testing.T) {
	cfg := config.Config{
		WebhookVerifyToken: "test-verify-token",
	}
	mockUseCase := &MockWhatsAppUseCase{}
	handler := &WhatsAppHandler{
		whatsappUseCase: mockUseCase,
		config:          cfg,
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Setup query parameters with wrong mode
	c.Request = httptest.NewRequest("GET", "/api/v1/whatsapp/webhook?hub.mode=unsubscribe&hub.verify_token=test-verify-token&hub.challenge=test-challenge", nil)

	handler.VerifyWebhook(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var errorResponse domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "verification_failed", errorResponse.Error)
	assert.Equal(t, http.StatusForbidden, errorResponse.Code)
}
