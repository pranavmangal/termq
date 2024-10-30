package cerebras

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	cfg "github.com/pranavmangal/termq/config"
)

const API_URL = "https://api.cerebras.ai/v1/chat/completions"

var avlModels = []string{"llama3.1-8b", "llama3.1-70b"}

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

func isModelAvailable(model string) bool {
	for _, m := range avlModels {
		if m == model {
			return true
		}
	}

	return false
}

func RunQuery(query string, config cfg.Config) (string, error) {
	cb := config.Cerebras

	if !isModelAvailable(cb.Model) {
		return "", fmt.Errorf(`Model "%s" is not available on Cerebras`, cb.Model)
	}

	messages := []Message{
		{Role: "system", Content: config.SystemPrompt},
		{Role: "user", Content: query},
	}
	body := Request{Model: cb.Model, Messages: messages}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Could not create request body: %v", err)
	}

	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalf("Could not create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cb.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Unsuccessful HTTP request. Status code: %d", resp.StatusCode)
	}

	var jsonResp Response
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		log.Fatalf("Error decoding response JSON: %v", err)
	}

	if len(jsonResp.Choices) == 0 {
		return "", fmt.Errorf("No response from the API")
	}

	return jsonResp.Choices[0].Message.Content, nil
}
