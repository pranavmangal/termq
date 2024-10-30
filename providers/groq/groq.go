package groq

import (
	"fmt"

	"github.com/pranavmangal/termq/common"
	cfg "github.com/pranavmangal/termq/config"
)

const API_URL = "https://api.groq.com/openai/v1/chat/completions"

var avlModels = []string{
	"gemma-7b-it",
	"gemma2-9b-it",
	"llama-3.2-1b-preview",
	"llama-3.2-3b-preview",
	"llama-3.2-11b-vision-preview",
	"llama-3.2-90b-vision-preview",
	"llama-3.1-8b-instant",
	"llama-3.1-70b-versatile",
	"llama3-8b-8192",
	"llama3-70b-8192",
	"mixtral-8x7b-32768",
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
	groq := config.Groq

	if !common.IsModelAvailable(groq.Model, avlModels) {
		return "", fmt.Errorf(`Model "%s" is not available on Groq`, groq.Model)
	}

	messages := []Message{
		{Role: "system", Content: config.SystemPrompt},
		{Role: "user", Content: query},
	}
	body := Request{Model: groq.Model, Messages: messages}

	var jsonResp Response
	err := common.MakeRequest(API_URL, body, groq.ApiKey, &jsonResp)
	if err != nil {
		return "", err
	}

	if len(jsonResp.Choices) == 0 {
		return "", fmt.Errorf("No response from the API")
	}

	return jsonResp.Choices[0].Message.Content, nil

}
