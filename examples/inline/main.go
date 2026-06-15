// Command inline is an inline bot built on github.com/gotd/botapi. Type the
// bot's @username followed by a query in any chat and it offers article results
// echoing the query; picking one sends the text.
//
// It also demonstrates inline_message_id: each result carries an inline
// keyboard, so when a user picks a result Telegram assigns it an
// inline_message_id. That id is surfaced both on the chosen-inline-result update
// and on callback queries from the inline message's button — a bot can persist
// it to edit the inline message later.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token. Both inline mode and inline feedback (for chosen-inline-result updates)
// must be enabled for the bot via @BotFather:
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

	// A callback keyboard on the result makes Telegram mint an inline_message_id
	// for the sent message and lets the user interact with it afterwards.
	keyboard := botapi.InlineKeyboard(
		botapi.InlineRow(botapi.InlineKeyboardButton{Text: "🔁 Shout again", CallbackData: "shout"}),
	)

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
				ReplyMarkup: keyboard,
			},
			&botapi.InlineQueryResultArticle{
				ID:          "echo",
				Title:       "Echo",
				Description: q,
				InputMessageContent: &botapi.InputTextMessageContent{
					MessageText: q,
				},
				ReplyMarkup: keyboard,
			},
		}

		return c.AnswerInline(results, botapi.WithInlineCacheTime(1))
	})

	// Fires when a user picks one of the offered results. The inline_message_id
	// identifies the message that was sent into the user's chat; persist it to
	// edit the message later. Requires inline feedback enabled in @BotFather.
	bot.OnChosenInlineResult(func(c *botapi.Context) error {
		r := c.Update.ChosenInlineResult

		log.Info("Inline result chosen",
			zap.String("result_id", r.ResultID),
			zap.String("query", r.Query),
			zap.String("inline_message_id", r.InlineMessageID),
		)

		return nil
	})

	// Callback queries from the inline message's button carry an
	// inline_message_id instead of a Message (the bot never sees the host chat).
	bot.OnCallbackQuery(func(c *botapi.Context) error {
		cq := c.Update.CallbackQuery
		if cq.InlineMessageID == "" {
			// Callback from a normal message (not an inline one).
			return c.AnswerCallback()
		}

		log.Info("Callback on inline message",
			zap.String("data", cq.Data),
			zap.String("inline_message_id", cq.InlineMessageID),
		)

		return c.AnswerCallback(botapi.WithCallbackText("Got it!"), botapi.WithCallbackAlert())
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting inline bot")

	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}
