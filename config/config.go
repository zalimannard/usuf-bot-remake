package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LoggerValue  Logger  `envconfig:"LOGGER"`
	DiscordValue Discord `envconfig:"DISCORD"`
	YouTubeValue YouTube `envconfig:"YOUTUBE"`
}

func (c Config) Logger() Logger {
	return c.LoggerValue
}

func (c Config) Discord() Discord {
	return c.DiscordValue
}

func (c Config) YouTube() YouTube {
	return c.YouTubeValue
}

type Logger struct {
	LevelValue string `envconfig:"LEVEL"`
}

func (l Logger) Level() string {
	return l.LevelValue
}

type Discord struct {
	PrefixValue string `envconfig:"PREFIX"`
	TokenValue  string `envconfig:"TOKEN"`
}

func (d Discord) Prefix() string {
	return d.PrefixValue
}

func (d Discord) Token() string {
	return d.TokenValue
}

type YouTube struct {
	APIKeyValue string `envconfig:"API_KEY"`
}

func (d YouTube) APIKey() string {
	return d.APIKeyValue
}

func Parse() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
