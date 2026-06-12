package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/bot"
	"github.com/cQu1x/Incident-War-Room/internal/config"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
	"github.com/cQu1x/Incident-War-Room/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	pool, err := repository.NewPool(ctx, cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer pool.Close()

	tgBot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("%v", errs.Wrapf(errs.KindUnavailable, "main", err, "connect to Telegram Bot API"))
	}

	bot.Register(tgBot)

	fmt.Println("Bot started")
	tgBot.Start()
}
