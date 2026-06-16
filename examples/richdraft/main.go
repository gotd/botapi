// Command richdraft demonstrates streaming a rich message draft (Bot API 10.1)
// with github.com/gotd/botapi: the pattern an AI bot uses to show a live preview
// while it generates output.
//
// On /stream the bot builds a rich message incrementally and, after each new
// block, calls SendRichMessageDraft with the SAME draft id. Each call replaces an
// ephemeral ~30-second preview in the chat (it does not persist). When the
// "generation" finishes, the bot calls SendRichMessage once to persist the
// finished message.
//
// The draft content is a tg.InputRichMessageClass, built exactly as for
// SendRichMessage: from native page blocks with rich.New(...).Input(), or — for a
// whole document parsed by Telegram's servers — rich.HTML / rich.Markdown.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token:
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/richdraft
package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gotd/log/logzap"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram/message/rich"
	"github.com/gotd/td/tg"

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

	bot.Use(botapi.Recover(), botapi.Logging())

	bot.OnCommand("start", "Explain the demo", func(c *botapi.Context) error {
		_, err := c.Reply("Send /stream to watch a rich message stream in, then persist.")
		return err
	})

	bot.OnCommand("stream", "Stream a rich message draft, then finalize it", stream)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting richdraft bot")

	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}

// stream simulates a bot generating a rich message block by block, previewing
// each step as a draft and persisting the result at the end.
func stream(c *botapi.Context) error {
	chat, ok := c.Chat()
	if !ok {
		return nil
	}

	// A non-zero id that identifies this draft; reuse it across the whole stream
	// so each call updates the same preview instead of starting a new one.
	draftID := time.Now().UnixNano()

	// The body grows one block at a time, as if generated token by token.
	body := []tg.PageBlockClass{rich.Heading2(rich.Plain("Streaming a rich message"))}

	for _, part := range []string{
		"First, the bot drafts an opening paragraph…",
		"…then it adds detail: every SendRichMessageDraft call updates the same ephemeral preview.",
		"Finally it wraps up. Nothing is persisted until the bot sends the finished message.",
	} {
		body = append(body, rich.Paragraph(rich.Plain(part)))

		// Preview the growing partial message. rich.New(...).Input() yields the
		// same tg.InputRichMessageClass SendRichMessage takes; rich.Markdown /
		// rich.HTML would work here too for server-parsed documents.
		if err := c.Bot.SendRichMessageDraft(c, chat, draftID, rich.New(body...).Input()); err != nil {
			return err
		}

		// Stand in for generation latency. The preview is ephemeral (~30s).
		time.Sleep(800 * time.Millisecond)
	}

	// Generation finished: persist the complete message with SendRichMessage.
	body = append(body, rich.Footer(rich.Plain("sent with SendRichMessage")))

	_, err := c.Bot.SendRichMessage(c, chat, rich.New(body...).Input())

	return err
}
