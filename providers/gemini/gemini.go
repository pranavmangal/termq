package gemini

import (
	"fmt"

	"github.com/pranavmangal/termq/common"
	cfg "github.com/pranavmangal/termq/config"
)

const API_URL = "https://generativelanguage.googleapis.com/v1beta"

var avlModels = []string{
	"gemini-1.5-flash",
	"gemini-1.5-flash-8b",
	"gemini-1.5-pro",
	"gemini-1.0-pro",
}

type Parts struct {
	Text string `json:"text"`
}
type Message struct {
	Parts Parts `json:"parts"`
}

type Request struct {
	SystemInstruction Message `json:"system_instruction"`
	Contents          Message `json:"contents"`
}

type Response struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func getFullUrl(model string, apiKey string) string {
	fullUrl := fmt.Sprintf(API_URL+"/models/%s:generateContent?key=%s", model, apiKey)

	return fullUrl
}

func RunQuery(query string, config cfg.Config) (string, error) {
	gc := config.Gemini

	if !common.IsModelAvailable(gc.Model, avlModels) {
		return "", fmt.Errorf(`Model "%s" is not available on Google Gemini`, gc.Model)
	}

	body := Request{
		SystemInstruction: Message{Parts{config.SystemPrompt}},
		Contents:          Message{Parts{query}},
	}

	fullUrl := getFullUrl(gc.Model, gc.ApiKey)

	var jsonResp Response
	err := common.MakeRequest(fullUrl, body, "", &jsonResp)
	if err != nil {
		return "", err
	}

	if len(jsonResp.Candidates) == 0 || len(jsonResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("No response from the API")
	}

	return jsonResp.Candidates[0].Content.Parts[0].Text, nil
}