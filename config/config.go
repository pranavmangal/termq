package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pelletier/go-toml/v2"
)

type Provider struct {
	Model  string `toml:"model"`
	ApiKey string `toml:"api_key"`
}

func (c Provider) IsValid() bool {
	if c.Model == "" || c.ApiKey == "" {
		return false
	}

	return true
}

type Config struct {
	SystemPrompt string   `toml:"system_prompt"`
	Groq         Provider `toml:"groq"`
	Cerebras     Provider `toml:"cerebras"`
	Gemini       Provider `toml:"gemini"`
}

const defaultSystemPrompt = `You are a helpful assistant run from the terminal. Be very concise. Respond in markdown. Code blocks should contain language identifiers.`

func getConfigFilePath() string {
	toolName := "termq"
	configFileName := "config.toml"

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error fetching home directory: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", toolName, configFileName)
	default:
		return filepath.Join(home, ".config", toolName, configFileName)
	}
}

func Exists() bool {
	configPath := getConfigFilePath()

	info, err := os.Stat(configPath)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func Create() string {
	configPath := getConfigFilePath()

	onlyOwnerCanWrite := os.FileMode(0644)
	dir := filepath.Dir(configPath)
	err := os.MkdirAll(dir, onlyOwnerCanWrite)
	if err != nil {
		log.Fatalf("Failed to create config directory at: %v", err)
	}

	defaultConfig := Config{
		SystemPrompt: defaultSystemPrompt,
		Groq:         Provider{},
		Cerebras:     Provider{},
	}

	configToml, err := toml.Marshal(defaultConfig)
	if err != nil {
		log.Fatalf("Failed to create default config: %v", err)
	}

	err = os.WriteFile(configPath, configToml, onlyOwnerCanWrite)
	if err != nil {
		log.Fatalf("Failed to write default config at %s: %v", configPath, err)
	}

	return configPath
}

func Parse() (Config, error) {
	configPath := getConfigFilePath()
	configFileContents, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file at %s: %v", configPath, err)
	}

	var config Config
	err = toml.Unmarshal(configFileContents, &config)

	return config, err
}
