package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gotd/botapi"
	"github.com/gotd/botapi/storage"
	"github.com/gotd/log/logzap"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log, _ := zap.NewProduction()
	defer func() { _ = log.Sync() }()

	store, err := storage.Open(os.Getenv("STORAGE_PATH"))
	if err != nil {
		log.Fatal("Open storage", zap.Error(err))
	}
	defer func() { _ = store.Close() }()

	appID, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		log.Fatal("App ID", zap.Error(err))
	}

	bot, err := botapi.New(os.Getenv("BOT_TOKEN"), botapi.Options{
		AppID:     appID,
		AppHash:   os.Getenv("APP_HASH"),
		Logger:    logzap.New(log),
		Storage:   store,
		FloodWait: true,
	})
	if err != nil {
		log.Fatal("Create bot", zap.Error(err))
	}

	bot.UseOuter(func(next botapi.Handler) botapi.Handler {
		return func(c *botapi.Context) error {
			return next(c)
		}
	})

	bot.Use(botapi.Recover(), botapi.Timeout(time.Minute), botapi.Logging())

	log.Info("Starting bot")
	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}
