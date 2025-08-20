package infrastructure

import (
	"anyzzapp/internal/config"
	"anyzzapp/internal/infrastructure/entity"
	"anyzzapp/pkg/domain"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type LLMRepository struct {
	config config.Config
}

func NewLLMRepository(config config.Config) domain.LLMRepository {
	return &LLMRepository{config: config}
}

func (r *LLMRepository) SendMessage(prompt string) (string, error) {
	payload := entity.Request{
		Prompt: prompt,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", r.config.LLMUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var response entity.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	return response.Response, nil
}
