package common

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const (
	modelsFile = "available_models.json"
	modelsURL  = "https://cdn.jsdelivr.net/gh/pranavmangal/termq@master/models/available_models.json"
)

type AvailableModels map[string][]string

func getModelsFilePath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("Could not get cache directory: %v", err)
	}

	toolName := "termq"
	return filepath.Join(cacheDir, toolName, modelsFile)
}

func ModelCacheExists() bool {
	_, err := os.Stat(getModelsFilePath())
	return err == nil
}

func GetModels(provider string) ([]string, error) {
	file, err := os.Open(getModelsFilePath())
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	var data AvailableModels
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		log.Fatalf("Could not decode JSON: %v", err)
	}

	modelsList, ok := data[provider]
	if !ok {
		return []string{}, err
	}

	return modelsList, nil
}

func fetchLatestAvailableModels() ([]byte, error) {
	resp, err := http.Get(modelsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	models, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func CreateModelCache() {
	modelData, err := fetchLatestAvailableModels()
	if err != nil {
		log.Fatalf("Could not fetch latest available models: %v", err)
	}

	modelsFilePath := getModelsFilePath()
	if err := os.MkdirAll(path.Dir(modelsFilePath), 0755); err != nil {
		log.Fatalf("Could not create cache directory: %v", err)
	}

	err = os.WriteFile(modelsFilePath, modelData, 0644)
	if err != nil {
		log.Fatalf("Could not write to cache file: %v", err)
	}
}

func UpdateModelCache() {
	modelData, err := fetchLatestAvailableModels()
	if err != nil {
		return
	}

	os.WriteFile(getModelsFilePath(), modelData, 0644)
}
