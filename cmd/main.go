package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("bot: %v", err)
	}

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Incident War Room is running.")
	})

	fmt.Println("Bot started")
	bot.Start()
}
