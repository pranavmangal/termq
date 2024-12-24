package groq

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pranavmangal/termq/common"
	cfg "github.com/pranavmangal/termq/config"
)

const API_URL = "https://api.groq.com/openai/v1/chat/completions"

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

	avlModels, err := common.GetModels("groq")
	if err == nil && !slices.Contains(avlModels, groq.Model) {
		return "", fmt.Errorf(
			"The model '%s' is not available on Groq. Available models:\n%s",
			groq.Model,
			strings.Join(avlModels, "\n"),
		)
	}

	messages := []Message{
		{Role: "system", Content: config.SystemPrompt},
		{Role: "user", Content: query},
	}
	body := Request{Model: groq.Model, Messages: messages}

	var jsonResp Response
	err = common.MakeRequest(API_URL, body, groq.ApiKey, &jsonResp)
	if err != nil {
		return "", err
	}

	if len(jsonResp.Choices) == 0 {
		return "", fmt.Errorf("No response from the API")
	}

	return jsonResp.Choices[0].Message.Content, nil

}
