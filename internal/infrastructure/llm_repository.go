package infrastructure

import (
	"anyzzapp/internal/config"
	client2 "anyzzapp/internal/infrastructure/client"
	"anyzzapp/internal/infrastructure/entity"
	"anyzzapp/pkg/domain"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

type LLMRepository struct {
	config config.Config
	client client2.HttpClient
}

func NewLLMRepository(config config.Config, client client2.HttpClient) domain.LLMRepository {
	return &LLMRepository{
		config: config,
		client: client,
	}
}

func (r *LLMRepository) SendMessage(prompt string) (string, error) {
	payload := entity.Request{
		Prompt: prompt,
	}
	// Execute POST
	resp, err := r.client.Post(payload, r.config.LLMUrl)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("body: %s\n", string(body))

	var response entity.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	defer resp.Body.Close()
	log.Debug().Msgf("llm response: %v", response)

	return response.Response, nil
}
