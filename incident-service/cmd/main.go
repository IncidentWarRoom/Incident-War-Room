package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/bot"
	"github.com/cQu1x/Incident-War-Room/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	tgBot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("bot: %v", err)
	}

	bot.Register(tgBot)

	fmt.Println("Bot started")
	tgBot.Start()
}
