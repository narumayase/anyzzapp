package infrastructure

import (
	"anyzzapp/internal/config"
	"anyzzapp/internal/infrastructure/entity"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewLLMRepository(t *testing.T) {
	cfg := config.Config{
		LLMUrl:         "https://api.llm.example.com/chat",
		LLMBearerToken: "test-bearer-token",
	}
	mockClient := &MockHttpClient{}

	repo := NewLLMRepository(cfg, mockClient)

	assert.NotNil(t, repo)
	assert.IsType(t, &LLMRepository{}, repo)
}

func TestLLMRepository_SendMessage_Success(t *testing.T) {
	cfg := config.Config{
		LLMUrl:         "https://api.llm.example.com/chat",
		LLMBearerToken: "test-bearer-token",
	}
	mockClient := &MockHttpClient{}
	repo := &LLMRepository{
		config: cfg,
		client: mockClient,
	}

	prompt := "Hello, how are you?"
	expectedPayload := entity.Request{
		Prompt: prompt,
	}

	responseBody := entity.Response{
		Response: "I'm doing well, thank you!",
	}
	responseJSON, _ := json.Marshal(responseBody)
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseJSON)),
	}

	mockClient.On("Post", expectedPayload, cfg.LLMUrl).Return(mockResponse, nil)

	result, err := repo.SendMessage(prompt)

	assert.NoError(t, err)
	assert.Equal(t, "I'm doing well, thank you!", result)
	mockClient.AssertExpectations(t)
}

func TestLLMRepository_SendMessage_HttpClientError(t *testing.T) {
	cfg := config.Config{
		LLMUrl:         "https://api.llm.example.com/chat",
		LLMBearerToken: "test-bearer-token",
	}
	mockClient := &MockHttpClient{}
	repo := &LLMRepository{
		config: cfg,
		client: mockClient,
	}

	prompt := "Hello, how are you?"
	expectedError := errors.New("network error")

	mockClient.On("Post", mock.Anything, mock.Anything).Return((*http.Response)(nil), expectedError)

	result, err := repo.SendMessage(prompt)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, expectedError, err)
	mockClient.AssertExpectations(t)
}

func TestLLMRepository_SendMessage_InvalidJSON(t *testing.T) {
	cfg := config.Config{
		LLMUrl:         "https://api.llm.example.com/chat",
		LLMBearerToken: "test-bearer-token",
	}
	mockClient := &MockHttpClient{}
	repo := &LLMRepository{
		config: cfg,
		client: mockClient,
	}

	prompt := "Hello, how are you?"

	// Invalid JSON response
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
	}

	mockClient.On("Post", mock.Anything, mock.Anything).Return(mockResponse, nil)

	result, err := repo.SendMessage(prompt)

	assert.Error(t, err)
	assert.Empty(t, result)
	mockClient.AssertExpectations(t)
}

func TestLLMRepository_SendMessage_EmptyPrompt(t *testing.T) {
	cfg := config.Config{
		LLMUrl:         "https://api.llm.example.com/chat",
		LLMBearerToken: "test-bearer-token",
	}
	mockClient := &MockHttpClient{}
	repo := &LLMRepository{
		config: cfg,
		client: mockClient,
	}

	prompt := ""
	expectedPayload := entity.Request{
		Prompt: prompt,
	}

	responseBody := entity.Response{
		Response: "Please provide a valid prompt",
	}
	responseJSON, _ := json.Marshal(responseBody)
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseJSON)),
	}

	mockClient.On("Post", expectedPayload, cfg.LLMUrl).Return(mockResponse, nil)

	result, err := repo.SendMessage(prompt)

	assert.NoError(t, err)
	assert.Equal(t, "Please provide a valid prompt", result)
	mockClient.AssertExpectations(t)
}
