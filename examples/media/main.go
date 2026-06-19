// Command media is a bot built on github.com/gotd/botapi that demonstrates
// sending and receiving media:
//
//   - /photo sends a photo by URL.
//   - sending the bot a photo echoes it back by file_id (no re-upload).
//   - sending the bot a document echoes it back and reports its size via
//     GetFile.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token:
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/media
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gotd/log/logzap"
	"go.uber.org/zap"

	"github.com/gotd/botapi"
	"github.com/gotd/botapi/storage"
)

// A small public domain image to send for /photo.
const samplePhotoURL = "https://upload.wikimedia.org/wikipedia/commons/3/3a/Cat03.jpg"

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

	bot.Use(botapi.Recover(), botapi.Timeout(time.Minute), botapi.Logging())

	bot.OnCommand("start", "Show the welcome message", func(c *botapi.Context) error {
		_, err := c.Reply("Send me a photo or a file and I'll send it back. Or try /photo.")
		return err
	})

	// Send a photo by URL — Telegram fetches it server-side.
	bot.OnCommand("photo", "Send a sample photo", func(c *botapi.Context) error {
		chat, ok := c.Chat()
		if !ok {
			return nil
		}

		_, err := c.Bot.SendPhoto(c, chat, botapi.FileURL(samplePhotoURL), "Here's a cat 🐈")

		return err
	})

	// Echo incoming photos back by file_id (no download/re-upload round trip).
	bot.OnMessage(func(c *botapi.Context) error {
		photos := c.Message().Photo
		largest := photos[len(photos)-1] // last size is the highest resolution
		chat, _ := c.Chat()
		_, err := c.Bot.SendPhoto(c, chat, botapi.FileID(largest.FileID),
			fmt.Sprintf("%d×%d, file_unique_id=%s", largest.Width, largest.Height, largest.FileUniqueID))

		return err
	}, hasPhoto)

	// Echo incoming documents back, reporting metadata resolved from the file_id.
	bot.OnMessage(func(c *botapi.Context) error {
		doc := c.Message().Document
		chat, _ := c.Chat()

		file, err := c.Bot.GetFile(c, doc.FileID)
		if err != nil {
			return err
		}

		caption := fmt.Sprintf("%s (%d bytes)\nfile_unique_id=%s",
			doc.FileName, doc.FileSize, file.FileUniqueID)

		_, err = c.Bot.SendDocument(c, chat, botapi.FileID(doc.FileID), caption)

		return err
	}, hasDocument)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting media bot")

	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}

// hasPhoto matches messages that carry a photo.
func hasPhoto(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && len(m.Photo) > 0
}

// hasDocument matches messages that carry a document.
func hasDocument(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && m.Document != nil
}
