package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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
	Host           string   `yaml:"host" env:"MAIL_HOST"`
	Port           int      `yaml:"port" env:"MAIL_PORT"`
	Username       string   `yaml:"username" env:"MAIL_USERNAME"`
	Password       string   `yaml:"password" env:"MAIL_PASSWORD"`
	CheckInterval  Duration `yaml:"check_interval" env:"MAIL_CHECK_INTERVAL"`
	AllowedDomains []string `yaml:"allowed_domains"`
}

type DiscordConfig struct {
	WebhookURL string `yaml:"webhook_url" env:"DISCORD_WEBHOOK_URL"`
}

type Config struct {
	Mail    MailConfig    `yaml:"mail"`
	Discord DiscordConfig `yaml:"discord"`
}

var App Config

// Load reads configuration from config.yaml and .env file
func Load() error {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Read config file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, &App); err != nil {
		return err
	}

	// Override with environment variables if they exist
	if envHost := os.Getenv("MAIL_HOST"); envHost != "" {
		App.Mail.Host = envHost
	}
	if envPort := os.Getenv("MAIL_PORT"); envPort != "" {
		App.Mail.Port = parseInt(envPort, App.Mail.Port)
	}
	if envUsername := os.Getenv("MAIL_USERNAME"); envUsername != "" {
		App.Mail.Username = envUsername
	}
	if envPassword := os.Getenv("MAIL_PASSWORD"); envPassword != "" {
		App.Mail.Password = envPassword
	}
	if envWebhookURL := os.Getenv("DISCORD_WEBHOOK_URL"); envWebhookURL != "" {
		App.Discord.WebhookURL = envWebhookURL
	}

	return nil
}

// parseInt helper function
func parseInt(s string, defaultValue int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return v
}
