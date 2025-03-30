package scheduler

import (
	"fmt"
	"log"

	"scg-mail-discord-bot/internal/config"
	"scg-mail-discord-bot/internal/discord"
	"scg-mail-discord-bot/internal/mail"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) Start() error {
	_, err := s.cron.AddFunc(fmt.Sprintf("@every %s", config.App.Mail.CheckInterval.Duration), s.checkEmails)
	if err != nil {
		return err
	}

	s.cron.Start()
	log.Printf("Scheduler started with interval: %s", config.App.Mail.CheckInterval.Duration)
	return nil
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) checkEmails() {
	messages, err := mail.CheckEmails()
	if err != nil {
		log.Printf("Error checking emails: %v", err)
		return
	}

	if len(messages) > 0 {
		if err := discord.SendEmailNotification(messages); err != nil {
			log.Printf("Error sending Discord notifications: %v", err)
		}
	}
}
