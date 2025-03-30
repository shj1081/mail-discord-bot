package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration is a custom YAML unmarshaler for time.Duration
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	parsed, err := time.ParseDuration(value.Value)
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

// types for config
type MailConfig struct {
	Host           string   `yaml:"host"`
	Port           int      `yaml:"port"`
	Username       string   `yaml:"username"`
	Password       string   `yaml:"password"`
	CheckInterval  Duration `yaml:"check_interval"`
	AllowedDomains []string `yaml:"allowed_domains"`
}

type DiscordConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type Config struct {
	Mail    MailConfig    `yaml:"mail"`
	Discord DiscordConfig `yaml:"discord"`
}

var App Config

// Load reads configuration from config.yaml
func Load() error {
	// Read config file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	// Parse YAML directly without environment variable expansion
	if err := yaml.Unmarshal(data, &App); err != nil {
		return err
	}

	return nil
}
