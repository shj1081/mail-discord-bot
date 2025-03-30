package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"scg-mail-discord-bot/internal/config"

	"github.com/emersion/go-imap"
)

// types for discord webhook message
type DiscordWebhookMessage struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

// types for discord embed
type DiscordEmbed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Fields      []EmbedField `json:"fields"`
}

// types for discord embed field
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// send email notification to discord
func SendEmailNotification(messages []*imap.Message) error {
	if len(messages) == 0 {
		return nil
	}

	var fields []EmbedField
	for _, msg := range messages {
		fields = append(fields, EmbedField{
			Name: msg.Envelope.Subject,
			Value: fmt.Sprintf("**From:** %s\n**Date:** %s",
				msg.Envelope.From[0].Address(),
				msg.Envelope.Date.Format("2006-01-02 15:04:05")),
			Inline: false,
		})
	}

	webhookMessage := DiscordWebhookMessage{
		Embeds: []DiscordEmbed{
			{
				Title:       fmt.Sprintf("📫 New emails: %d", len(messages)),
				Description: "Check your SCG email inbox for unread emails",
				Color:       0x00ff00,
				Fields:      fields,
			},
		},
	}

	jsonData, err := json.Marshal(webhookMessage)
	if err != nil {
		log.Printf("Error marshaling webhook message: %v", err)
		return err
	}

	resp, err := http.Post(config.App.Discord.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending Discord webhook: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Discord webhook returned status: %d, body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}

	return nil
}
