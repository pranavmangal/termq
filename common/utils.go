package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func IsModelAvailable(model string, avlModels []string) bool {
	for _, m := range avlModels {
		if m == model {
			return true
		}
	}

	return false
}

func MakeRequest[Body, Res any](url string, body Body, apiKey string, res *Res) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Could not create request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalf("Could not create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Unsuccessful HTTP request. Status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Fatalf("Error decoding response JSON: %v", err)
	}

	return nil
}
