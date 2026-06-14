// Command pp is a bot that pretty-prints the raw MTProto message behind any
// Bot API message. Send it anything and it replies with the tdp.Format dump of
// the underlying tg.Message, reachable through Message.Raw — useful for
// inspecting fields the typed Bot API surface does not expose.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token:
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/pp
package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gotd/log/logzap"
	"github.com/gotd/td/tdp"
	"go.uber.org/zap"

	"github.com/gotd/botapi"
	"github.com/gotd/botapi/storage"
)

func main() {
	log, _ := zap.NewProduction()
	defer func() { _ = log.Sync() }()

	appID, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		log.Fatal("APP_ID must be a number (see https://my.telegram.org)", zap.Error(err))
	}

	// Persist session, peers and update state so the bot resumes across restarts.
	store, err := storage.Open("session.bbolt")
	if err != nil {
		log.Fatal("Open storage", zap.Error(err))
	}

	defer func() { _ = store.Close() }()

	bot, err := botapi.New(os.Getenv("BOT_TOKEN"), botapi.Options{
		AppID:   appID,
		AppHash: os.Getenv("APP_HASH"),
		Logger:  logzap.New(log),
		Storage: store,
	})
	if err != nil {
		log.Fatal("Create bot", zap.Error(err))
	}

	bot.Use(botapi.Recover(), botapi.Timeout(30*time.Second))

	bot.OnCommand("start", "Show the welcome message", func(c *botapi.Context) error {
		_, err := c.Reply("Reply to any message with /pp to dump its raw MTProto message.")
		return err
	})

	// /pp pretty-prints the MTProto message the command replies to. The replied
	// message is not carried by the update, so it is fetched with Bot.GetMessage,
	// whose raw tg.Message (via Message.Raw) is dumped with tdp.Format.
	bot.OnCommand("pp", "Pretty-print the replied-to MTProto message", func(c *botapi.Context) error {
		msg := c.Message()

		reply := msg.ReplyToMessage
		if reply == nil {
			_, err := c.Reply("Reply to a message with /pp to dump it.")
			return err
		}

		target, err := c.Bot.GetMessage(c, botapi.ID(msg.Chat.ID), reply.MessageID)
		if err != nil {
			_, replyErr := c.Reply("Could not fetch the replied message: " + err.Error())
			if replyErr != nil {
				return replyErr
			}

			return err
		}

		_, err = c.Reply(tdp.Format(target.Raw()))

		return err
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting pp bot")

	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}
