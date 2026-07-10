package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey   string
	MistralAPIKey  string
	DiscordWebhook string
}

func Load() (*Config, error) {
	// loads .env file variables into process memory
	_ = godotenv.Load()
	cfg := &Config{
		GeminiAPIKey:   os.Getenv("GEMINI_API_KEY"),
		MistralAPIKey:  os.Getenv("MISTRAL_API_KEY"),
		DiscordWebhook: os.Getenv("DISCORD_WEBHOOK"),
	}
	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY missing")
	}
	if cfg.MistralAPIKey == "" {
		return nil, fmt.Errorf("MISTRAL_API_KEY missing")
	}
	if cfg.DiscordWebhook == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK missing")
	}
	return cfg, nil
}
