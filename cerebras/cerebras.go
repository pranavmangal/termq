package cerebras

import (
	"fmt"

	"github.com/pranavmangal/termq/common"
	cfg "github.com/pranavmangal/termq/config"
)

const API_URL = "https://api.cerebras.ai/v1/chat/completions"

var avlModels = []string{
	"llama3.1-8b",
	"llama3.1-70b",
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func RunQuery(query string, config cfg.Config) (string, error) {
	cb := config.Cerebras

	if !common.IsModelAvailable(cb.Model, avlModels) {
		return "", fmt.Errorf(`Model "%s" is not available on Cerebras`, cb.Model)
	}

	messages := []Message{
		{Role: "system", Content: config.SystemPrompt},
		{Role: "user", Content: query},
	}
	body := Request{Model: cb.Model, Messages: messages}

	var jsonResp Response
	err := common.MakeRequest(API_URL, body, cb.ApiKey, &jsonResp)
	if err != nil {
		return "", err
	}

	if len(jsonResp.Choices) == 0 {
		return "", fmt.Errorf("No response from the API")
	}

	return jsonResp.Choices[0].Message.Content, nil
}