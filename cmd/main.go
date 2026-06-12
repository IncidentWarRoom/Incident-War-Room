package main

import (
	"fmt"
	"log"
	"strings"
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

	bot.Handle("/incident", func(c telebot.Context) error {
		args := c.Args()
		if len(args) == 0 {
			return c.Send("Usage:\n/incident create — open a new incident\n/incident <message> — add an update to the timeline")
		}

		switch args[0] {
		case "create":
			return c.Send("[stub] Incident created. (not implemented yet)")
		default:
			message := strings.Join(args, " ")
			return c.Send("[stub] Update added to timeline: " + message + " (not implemented yet)")
		}
	})

	bot.Handle("/timeline", func(c telebot.Context) error {
		return c.Send("[stub] Incident timeline is empty. (not implemented yet)")
	})

	fmt.Println("Bot started")
	bot.Start()
}
