// Command inline is an inline bot built on github.com/gotd/botapi. Type the
// bot's @username followed by a query in any chat and it offers article results
// echoing the query; picking one sends the text.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token (inline mode must be enabled for the bot via @BotFather):
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/inline
package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/gotd/log/logzap"
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
	store, err := storage.Open("inline-session.bbolt")
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

	bot.OnInlineQuery(func(c *botapi.Context) error {
		q := strings.TrimSpace(c.Update.InlineQuery.Query)
		if q == "" {
			return c.AnswerInline(nil)
		}
		results := []botapi.InlineQueryResult{
			&botapi.InlineQueryResultArticle{
				ID:          "upper",
				Title:       "UPPERCASE",
				Description: strings.ToUpper(q),
				InputMessageContent: &botapi.InputTextMessageContent{
					MessageText: strings.ToUpper(q),
				},
			},
			&botapi.InlineQueryResultArticle{
				ID:          "echo",
				Title:       "Echo",
				Description: q,
				InputMessageContent: &botapi.InputTextMessageContent{
					MessageText: q,
				},
			},
		}
		return c.AnswerInline(results, botapi.WithInlineCacheTime(1))
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting inline bot")
	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}
