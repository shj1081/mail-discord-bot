package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"scg-mail-discord-bot/internal/config"
	"scg-mail-discord-bot/internal/scheduler"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Create and start scheduler
	sched := scheduler.NewScheduler()
	if err := sched.Start(); err != nil {
		log.Fatal("Error starting scheduler:", err)
	}
	defer sched.Stop()

	log.Println("Bot started. Press Ctrl+C to exit.")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
