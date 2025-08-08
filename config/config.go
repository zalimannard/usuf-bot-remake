package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LoggerValue  Logger  `envconfig:"LOGGER"`
	DiscordValue Discord `envconfig:"DISCORD"`
}

type Logger struct {
	LevelValue string `envconfig:"LEVEL"`
}

func (l Logger) Level() string {
	return l.LevelValue
}

func (c Config) Logger() Logger {
	return c.LoggerValue
}

func (c Config) Discord() Discord {
	return c.DiscordValue
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

func Parse() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
