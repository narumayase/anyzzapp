package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

type HttpClientImpl struct {
	client      *http.Client
	bearerToken string
}

type HttpClient interface {
	Post(payload interface{}, url string) (*http.Response, error)
}

func NewHttpClient(client *http.Client, bearerToken string) HttpClient {
	return &HttpClientImpl{
		client:      client,
		bearerToken: bearerToken,
	}
}

func (c *HttpClientImpl) Post(payload interface{}, url string) (*http.Response, error) {
	// TODO add context

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	log.Debug().Msgf("payload to send: %s", string(jsonPayload))
	log.Debug().Msgf("url %s", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	// TODO add content type?

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	return resp, nil
}
