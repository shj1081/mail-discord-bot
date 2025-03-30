package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

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

func Load() error {
	// Try to load from config.yaml first
	if err := loadFromYaml(); err == nil {
		return nil
	}

	// If config.yaml doesn't exist or has error, try environment variables
	return loadFromEnv()
}

func loadFromYaml() error {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &App); err != nil {
		return err
	}

	return validate()
}

func loadFromEnv() error {
	App.Mail.Host = os.Getenv("MAIL_HOST")
	App.Mail.Port = parseEnvInt("MAIL_PORT", 993)
	App.Mail.Username = os.Getenv("MAIL_USERNAME")
	App.Mail.Password = os.Getenv("MAIL_PASSWORD")

	if interval := os.Getenv("MAIL_CHECK_INTERVAL"); interval != "" {
		duration, err := time.ParseDuration(interval)
		if err == nil {
			App.Mail.CheckInterval.Duration = duration
		}
	}

	if App.Mail.CheckInterval.Duration == 0 {
		App.Mail.CheckInterval.Duration = 1 * time.Hour
	}

	App.Discord.WebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")

	return validate()
}

func validate() error {
	if App.Mail.Host == "" {
		return fmt.Errorf("mail host is required")
	}
	if App.Mail.Username == "" {
		return fmt.Errorf("mail username is required")
	}
	if App.Mail.Password == "" {
		return fmt.Errorf("mail password is required")
	}
	if App.Discord.WebhookURL == "" {
		return fmt.Errorf("discord webhook URL is required")
	}
	return nil
}

func parseEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return int(parsed)
		}
	}
	return defaultValue
}
