package http

import (
	"anyzzapp/internal/config"
	"anyzzapp/pkg/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWhatsAppUseCase for router tests
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

func TestNewRouter(t *testing.T) {
	cfg := config.Config{
		WebhookVerifyToken: "test-token",
	}
	mockUseCase := &MockWhatsAppUseCase{}

	router := NewRouter(cfg, mockUseCase)

	assert.NotNil(t, router)
	assert.IsType(t, &gin.Engine{}, router)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	cfg := config.Config{}
	mockUseCase := &MockWhatsAppUseCase{}

	router := NewRouter(cfg, mockUseCase)

	// Test health endpoint
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "anyzzapp API is running", response["message"])
}

func TestRouter_WhatsAppEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		WebhookVerifyToken: "test-verify-token",
	}
	mockUseCase := &MockWhatsAppUseCase{}

	router := NewRouter(cfg, mockUseCase)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "POST /api/v1/whatsapp/send - should reach handler",
			method:         "POST",
			path:           "/api/v1/whatsapp/send",
			expectedStatus: http.StatusBadRequest, // Will fail due to invalid JSON, but route exists
		},
		{
			name:           "POST /api/v1/whatsapp/webhook - should reach handler",
			method:         "POST",
			path:           "/api/v1/whatsapp/webhook",
			expectedStatus: http.StatusBadRequest, // Will fail due to invalid JSON, but route exists
		},
		{
			name:           "GET /api/v1/whatsapp/webhook - should reach handler",
			method:         "GET",
			path:           "/api/v1/whatsapp/webhook",
			expectedStatus: http.StatusForbidden, // Will fail verification, but route exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, req)

			// We're just testing that the routes exist and are reachable
			// The actual handler logic is tested in handler_test.go
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRouter_NonExistentRoute(t *testing.T) {
	cfg := config.Config{}
	mockUseCase := &MockWhatsAppUseCase{}

	router := NewRouter(cfg, mockUseCase)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/non-existent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
