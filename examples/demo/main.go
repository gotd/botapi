// Command demo is a single bot that exercises every major feature of
// github.com/gotd/botapi, organized across several files so each subsystem is
// easy to read in isolation:
//
//   - main.go        — wiring: options, storage, logger, middleware, lifecycle
//   - commands.go    — slash commands: formatting, polls, dice, places, editing
//   - keyboards.go   — inline + reply keyboards and callback queries
//   - media.go       — sending media and reacting to incoming media
//   - inline.go      — inline mode and chosen-inline-result handling
//   - admin.go       — a group-scoped command set (chat info, pin, reactions)
//   - text.go        — free-text predicates and edited-message handling
//   - helpers.go     — small shared helpers
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token. Enable inline mode in @BotFather to test inline queries:
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/demo
//
// Logs are JSONL on stderr; pipe through github.com/go-faster/pl to read them
// (see ../README.md).
package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gotd/log/logzap"

	"github.com/gotd/botapi"
	"github.com/gotd/botapi/examples"
	"github.com/gotd/botapi/storage"
)

func main() {
	if err := run(); err != nil {
		// run already logged the detail; surface a non-zero exit.
		os.Exit(1)
	}
}

func run() error {
	log, err := examples.NewLogger()
	if err != nil {
		return err
	}

	defer func() { _ = log.Sync() }()

	appID, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		log.Error("APP_ID must be a number (see https://my.telegram.org)")

		return err
	}

	// Storage persists the MTProto session, peer access hashes and update gap
	// state. Without it every restart re-authorizes the bot and Telegram answers
	// repeated logins with a growing FLOOD_WAIT — so it is effectively mandatory
	// for any long-lived bot.
	store, err := storage.Open("demo-session.bbolt")
	if err != nil {
		log.Error("open storage: " + err.Error())

		return err
	}

	defer func() { _ = store.Close() }()

	bot, err := botapi.New(os.Getenv("BOT_TOKEN"), botapi.Options{
		AppID:   appID,
		AppHash: os.Getenv("APP_HASH"),
		Logger:  logzap.New(log),
		Storage: store,

		// FloodWait transparently waits out 429 FLOOD_WAIT limits instead of
		// failing; RequestsPerSecond is a coarse proactive guard against them.
		FloodWait:         true,
		RequestsPerSecond: 25,
		RequestBurst:      5,

		// OnStart fires once after the bot is authorized and gap recovery is
		// live; a good place for proactive startup work.
		OnStart: func(ctx context.Context) {
			log.Info("Bot authorized and serving updates")
		},
	})
	if err != nil {
		log.Error("create bot: " + err.Error())

		return err
	}

	// Global middleware applies to every handler, outermost first: recover from
	// panics, bound each handler with a timeout, then log the outcome. metrics is
	// a custom middleware defined in middleware.go.
	bot.Use(
		botapi.Recover(),
		botapi.Timeout(30*time.Second),
		botapi.Logging(),
		metrics(),
	)

	// Each module registers its own handlers and commands. Commands registered
	// via OnCommand are published to Telegram's command menu on Run.
	registerCommands(bot)
	registerKeyboards(bot)
	registerMedia(bot)
	registerStreaming(bot)
	registerInline(bot)
	registerAdmin(bot)
	registerText(bot)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting demo bot")

	// Run connects, authorizes and blocks serving updates until ctx is canceled.
	if err := bot.Run(ctx); err != nil && ctx.Err() == nil {
		log.Error("run: " + err.Error())

		return err
	}

	return nil
}
