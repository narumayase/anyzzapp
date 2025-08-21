package application

import (
	"anyzzapp/pkg/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWhatsAppRepository is a mock implementation of WhatsAppRepository
type MockWhatsAppRepository struct {
	mock.Mock
}

func (m *MockWhatsAppRepository) SendMessage(message domain.Message) (*domain.SendMessageResponse, error) {
	args := m.Called(message)
	return args.Get(0).(*domain.SendMessageResponse), args.Error(1)
}

func (m *MockWhatsAppRepository) MarkAsRead(phoneNumberID, messageID string) error {
	args := m.Called(phoneNumberID, messageID)
	return args.Error(0)
}

// MockLLMRepository is a mock implementation of LLMRepository
type MockLLMRepository struct {
	mock.Mock
}

func (m *MockLLMRepository) SendMessage(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

func TestNewWhatsAppUseCase(t *testing.T) {
	mockWhatsAppRepo := &MockWhatsAppRepository{}
	mockLLMRepo := &MockLLMRepository{}
	useCase := NewWhatsAppUseCase(mockWhatsAppRepo, mockLLMRepo)

	assert.NotNil(t, useCase)
	assert.IsType(t, &WhatsAppUseCase{}, useCase)
}

func TestWhatsAppUseCase_SendMessage_Success(t *testing.T) {
	mockWhatsAppRepo := &MockWhatsAppRepository{}
	mockLLMRepo := &MockLLMRepository{}
	useCase := &WhatsAppUseCase{
		whatsappRepo: mockWhatsAppRepo,
		llmRepo:      mockLLMRepo,
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

	mockWhatsAppRepo.On("SendMessage", message).Return(expectedResponse, nil)

	result, err := useCase.SendMessage(message)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockWhatsAppRepo.AssertExpectations(t)
}

func TestWhatsAppUseCase_SendMessage_EmptyPhoneNumberID(t *testing.T) {
	mockWhatsAppRepo := &MockWhatsAppRepository{}
	mockLLMRepo := &MockLLMRepository{}
	useCase := &WhatsAppUseCase{
		whatsappRepo: mockWhatsAppRepo,
		llmRepo:      mockLLMRepo,
	}

	message := domain.Message{
		PhoneNumberID: "",
		To:            "5491112345678",
		Content:       "Hello, World!",
		MessageType:   "text",
	}

	result, err := useCase.SendMessage(message)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "phone number ID is required")
}

func TestWhatsAppUseCase_SendMessage_EmptyTo(t *testing.T) {
	mockWhatsAppRepo := &MockWhatsAppRepository{}
	mockLLMRepo := &MockLLMRepository{}
	useCase := &WhatsAppUseCase{
		whatsappRepo: mockWhatsAppRepo,
		llmRepo:      mockLLMRepo,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "",
		Content:       "Hello, World!",
		MessageType:   "text",
	}

	result, err := useCase.SendMessage(message)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "recipient phone number is required")
}

func TestWhatsAppUseCase_SendMessage_EmptyContent(t *testing.T) {
	mockWhatsAppRepo := &MockWhatsAppRepository{}
	mockLLMRepo := &MockLLMRepository{}
	useCase := &WhatsAppUseCase{
		whatsappRepo: mockWhatsAppRepo,
		llmRepo:      mockLLMRepo,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "",
		MessageType:   "text",
	}

	result, err := useCase.SendMessage(message)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "message content is required")
}

func TestWhatsAppUseCase_SendMessage_RepositoryError(t *testing.T) {
	mockWhatsAppRepo := &MockWhatsAppRepository{}
	mockLLMRepo := &MockLLMRepository{}
	useCase := &WhatsAppUseCase{
		whatsappRepo: mockWhatsAppRepo,
		llmRepo:      mockLLMRepo,
	}

	message := domain.Message{
		PhoneNumberID: "123456789",
		To:            "5491112345678",
		Content:       "Hello, World!",
		MessageType:   "text",
	}

	expectedError := errors.New("API connection failed")
	mockWhatsAppRepo.On("SendMessage", message).Return((*domain.SendMessageResponse)(nil), expectedError)

	result, err := useCase.SendMessage(message)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to send message")
	mockWhatsAppRepo.AssertExpectations(t)
}

func TestWhatsAppUseCase_ProcessIncomingWebhook_Success(t *testing.T) {
	mockWhatsAppRepo := &MockWhatsAppRepository{}
	mockLLMRepo := &MockLLMRepository{}
	useCase := &WhatsAppUseCase{
		whatsappRepo: mockWhatsAppRepo,
		llmRepo:      mockLLMRepo,
	}

	webhook := &domain.WebhookRequest{
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
							Messages: []domain.WebhookMessage{
								{
									From:      "5491112345678",
									ID:        "msg_123",
									Timestamp: "1234567890",
									Text: &domain.WebhookText{
										Body: "Hello, how are you?",
									},
									Type: "text",
								},
							},
						},
						Field: "messages",
					},
				},
			},
		},
	}

	expectedLLMResponse := "I'm doing well, thank you!"
	expectedReplyMessage := domain.Message{
		PhoneNumberID: "123456789",
		To:            "541112345678", // removeNine applied
		Content:       expectedLLMResponse,
		MessageType:   "text",
	}
	expectedSendResponse := &domain.SendMessageResponse{
		MessageID: "reply_msg_123",
		Status:    "sent",
	}

	mockWhatsAppRepo.On("MarkAsRead", "123456789", "msg_123").Return(nil)
	mockLLMRepo.On("SendMessage", "Hello, how are you?").Return(expectedLLMResponse, nil)
	mockWhatsAppRepo.On("SendMessage", expectedReplyMessage).Return(expectedSendResponse, nil)

	err := useCase.ProcessIncomingWebhook(webhook)

	assert.NoError(t, err)
	mockWhatsAppRepo.AssertExpectations(t)
	mockLLMRepo.AssertExpectations(t)
}

func TestWhatsAppUseCase_ProcessIncomingWebhook_NilWebhook(t *testing.T) {
	mockWhatsAppRepo := &MockWhatsAppRepository{}
	mockLLMRepo := &MockLLMRepository{}
	useCase := &WhatsAppUseCase{
		whatsappRepo: mockWhatsAppRepo,
		llmRepo:      mockLLMRepo,
	}

	err := useCase.ProcessIncomingWebhook(nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "webhook data cannot be nil")
}

func TestWhatsAppUseCase_ProcessIncomingWebhook_LLMError(t *testing.T) {
	mockWhatsAppRepo := &MockWhatsAppRepository{}
	mockLLMRepo := &MockLLMRepository{}
	useCase := &WhatsAppUseCase{
		whatsappRepo: mockWhatsAppRepo,
		llmRepo:      mockLLMRepo,
	}

	webhook := &domain.WebhookRequest{
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
							Messages: []domain.WebhookMessage{
								{
									From:      "5491112345678",
									ID:        "msg_123",
									Timestamp: "1234567890",
									Text: &domain.WebhookText{
										Body: "Hello, how are you?",
									},
									Type: "text",
								},
							},
						},
						Field: "messages",
					},
				},
			},
		},
	}

	expectedError := errors.New("LLM service unavailable")

	mockWhatsAppRepo.On("MarkAsRead", "123456789", "msg_123").Return(nil)
	mockLLMRepo.On("SendMessage", "Hello, how are you?").Return("", expectedError)

	err := useCase.ProcessIncomingWebhook(webhook)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to process messages")
	mockWhatsAppRepo.AssertExpectations(t)
	mockLLMRepo.AssertExpectations(t)
}

func TestRemoveNine(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		expected    string
	}{
		{
			name:        "Argentine number with 9",
			phoneNumber: "5491112345678",
			expected:    "541112345678",
		},
		{
			name:        "Argentine number without 9",
			phoneNumber: "541112345678",
			expected:    "541112345678",
		},
		{
			name:        "Non-Argentine number",
			phoneNumber: "1234567890123",
			expected:    "1234567890123",
		},
		{
			name:        "Short number",
			phoneNumber: "123456789",
			expected:    "123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeNine(tt.phoneNumber)
			assert.Equal(t, tt.expected, result)
		})
	}
}
