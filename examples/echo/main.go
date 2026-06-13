// Command echo is a minimal bot built on github.com/gotd/botapi: it greets on
// /start and echoes any other text message back as a reply.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token:
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/echo
package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gotd/log/logzap"
	"go.uber.org/zap"

	"github.com/gotd/botapi"
)

func main() {
	log, _ := zap.NewProduction()
	defer func() { _ = log.Sync() }()

	appID, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		log.Fatal("APP_ID must be a number (see https://my.telegram.org)", zap.Error(err))
	}

	bot, err := botapi.New(os.Getenv("BOT_TOKEN"), botapi.Options{
		AppID:   appID,
		AppHash: os.Getenv("APP_HASH"),
		Logger:  logzap.New(log),
	})
	if err != nil {
		log.Fatal("Create bot", zap.Error(err))
	}

	// Middleware applies to every handler.
	bot.Use(botapi.Recover(), botapi.Timeout(30*time.Second))

	bot.OnCommand("start", "Show the welcome message", func(c *botapi.Context) error {
		_, err := c.Reply("Hi! Send me any text and I'll echo it back.")
		return err
	})

	// Any text message that is not a command.
	bot.OnMessage(func(c *botapi.Context) error {
		_, err := c.Reply(c.Message().Text)
		return err
	}, botapi.HasText(), botapi.Not(botapi.HasPrefix("/")))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting echo bot")
	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}
