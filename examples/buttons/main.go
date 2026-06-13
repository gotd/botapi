// Command buttons is a bot built on github.com/gotd/botapi that demonstrates
// inline keyboards and callback queries: /menu shows a keyboard, and tapping a
// button answers the callback and edits the message.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token:
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/buttons
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

func menu() *botapi.InlineKeyboardMarkup {
	return &botapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]botapi.InlineKeyboardButton{
			{
				{Text: "👍 Like", CallbackData: "vote:up"},
				{Text: "👎 Dislike", CallbackData: "vote:down"},
			},
			{
				{Text: "gotd/td", URL: "https://github.com/gotd/td"},
			},
		},
	}
}

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

	bot.Use(botapi.Recover(), botapi.Timeout(30*time.Second))

	bot.OnCommand("menu", "Show the voting menu", func(c *botapi.Context) error {
		_, err := c.Reply("How do you like this bot?", botapi.WithReplyMarkup(menu()))
		return err
	})

	// Respond to the inline-keyboard taps. The CallbackPrefix predicate filters
	// to our "vote:" buttons.
	bot.OnCallbackQuery(func(c *botapi.Context) error {
		choice := "🤷"
		switch c.Update.CallbackQuery.Data {
		case "vote:up":
			choice = "👍"
		case "vote:down":
			choice = "👎"
		}
		// Acknowledge the tap (removes the client-side loading state).
		if err := c.AnswerCallback(botapi.WithCallbackText("Thanks for voting " + choice)); err != nil {
			return err
		}
		// Edit the original message to reflect the choice.
		msg := c.Update.CallbackQuery.Message
		if msg == nil {
			return nil
		}
		_, err := c.Bot.EditMessageText(c, botapi.ID(msg.Chat.ID), msg.MessageID, "You voted "+choice)
		return err
	}, botapi.CallbackPrefix("vote:"))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting buttons bot")
	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}
