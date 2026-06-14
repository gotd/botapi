// Command business is a bot for verifying the Telegram Business surface over
// MTProto: inbound business updates, acting on behalf of a connected business
// account (text, media, albums, profile edits), and reading the connection.
//
// Setup:
//
//  1. Run the bot with an MTProto app identity (https://my.telegram.org) and a
//     BotFather token.
//
//  2. On a Telegram Premium account, open Settings → Business → Chatbots and add
//     this bot. That delivers the account's 1-to-1 chats to the bot over a
//     business connection.
//
//  3. From another account, write to the business account. Plain text is echoed
//     back (as the account). Messages starting with "!" are commands:
//
//     !help                 list commands
//     !ping                 reply as the account
//     !conn                 dump the business connection + granted rights
//     !balance              the account's Telegram Stars balance
//     !photo                send a generated photo as the account
//     !album                send a 2-photo album as the account
//     !name First [Last]    set the account's name
//     !bio <text>           set the account's bio
//     !username <name>      set the account's username
//     !avatar [public]      set the account's profile photo (public ⇒ fallback)
//     !rmavatar [public]    remove the account's profile photo
//
// Every command replies with "ok" or "<action> failed: <error>", so the online
// behavior of each assumption is visible in the chat.
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/business
package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gotd/log/logzap"
	"go.uber.org/zap"

	"github.com/gotd/botapi"
	"github.com/gotd/botapi/storage"
)

func main() {
	log, _ := zap.NewProduction()
	defer func() { _ = log.Sync() }()

	appID, err := atoi(os.Getenv("APP_ID"))
	if err != nil {
		log.Fatal("APP_ID must be a number (see https://my.telegram.org)", zap.Error(err))
	}

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

	bot.Use(botapi.Recover(), botapi.Timeout(60*time.Second))

	bot.OnCommand("start", "How to use this bot", func(c *botapi.Context) error {
		_, err := c.Reply("Add me to a Telegram Business account " +
			"(Settings → Business → Chatbots), then message that account. Send !help in a business chat.")

		return err
	})

	// Log the connection whenever it is established, disabled or its rights
	// change — lets you eyeball GetBusinessConnection / the rights mapping.
	bot.OnBusinessConnection(func(c *botapi.Context) error {
		conn := c.Update.BusinessConnection
		log.Info("Business connection update",
			zap.String("id", conn.ID),
			zap.Int64("user_id", conn.User.ID),
			zap.Bool("enabled", conn.IsEnabled),
			zap.Any("rights", conn.Rights),
		)

		return nil
	})

	bot.OnBusinessMessage(handleBusiness(log))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting business bot")

	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}

// handleBusiness processes a message delivered over a business connection.
func handleBusiness(log *zap.Logger) botapi.Handler {
	return func(c *botapi.Context) error {
		bm := c.BusinessMessage()
		if bm == nil {
			return nil
		}

		// Messages the account itself sent (including our own replies) come back
		// on the stream with Out set — ignore them to avoid a reply loop.
		if raw := bm.Raw(); raw != nil && raw.Out {
			return nil
		}

		bc, ok := c.Business()
		if !ok {
			return nil
		}

		log.Info("Business message",
			zap.String("connection", bc.ConnectionID()),
			zap.Int64("chat", bm.Chat.ID),
			zap.String("text", bm.Text),
		)

		text := strings.TrimSpace(bm.Text)
		if !strings.HasPrefix(text, "!") {
			_, err := c.Reply("echo (sent as the business account): " + text + "\nSend !help for commands.")
			return err
		}

		return dispatchCommand(c, bc, strings.Fields(text[1:]))
	}
}

// dispatchCommand runs one "!" command from a business chat.
func dispatchCommand(c *botapi.Context, bc *botapi.BusinessContext, fields []string) error {
	if len(fields) == 0 {
		return nil
	}

	cmd, args := fields[0], fields[1:]
	chat := botapi.ID(c.BusinessMessage().Chat.ID)

	switch cmd {
	case "help":
		return reply(c, helpText)
	case "ping":
		return reply(c, "pong — sent on behalf of the business account")
	case "conn":
		return showConnection(c, bc)
	case "balance":
		return showBalance(c, bc)
	case "name":
		return report(c, "set name", bc.SetName(c, first(args), rest(args)))
	case "bio":
		return report(c, "set bio", bc.SetBio(c, strings.Join(args, " ")))
	case "username":
		return report(c, "set username", bc.SetUsername(c, first(args)))
	case "photo":
		_, err := bc.SendPhoto(c, chat, samplePhoto(color.RGBA{R: 0x2a, G: 0x9d, B: 0x8f, A: 0xff}), "sent as the account")
		return report(c, "send photo", err)
	case "album":
		return sendAlbum(c, bc, chat)
	case "avatar":
		photo := botapi.InputProfilePhotoStatic{Photo: samplePhoto(color.RGBA{R: 0xe7, G: 0x6f, B: 0x51, A: 0xff})}
		return report(c, "set avatar", bc.SetProfilePhoto(c, photo, isPublic(args)))
	case "rmavatar":
		return report(c, "remove avatar", bc.RemoveProfilePhoto(c, isPublic(args)))
	default:
		return reply(c, "unknown command; send !help")
	}
}

const helpText = "Commands: !ping, !conn, !balance, !photo, !album, " +
	"!name First [Last], !bio <text>, !username <name>, !avatar [public], !rmavatar [public]"

func showConnection(c *botapi.Context, bc *botapi.BusinessContext) error {
	conn, err := bc.Connection(c)
	if err != nil {
		return report(c, "get connection", err)
	}

	return reply(c, fmt.Sprintf("connection %s\nuser_id %d\nenabled %t\nrights %+v",
		conn.ID, conn.User.ID, conn.IsEnabled, conn.Rights))
}

func showBalance(c *botapi.Context, bc *botapi.BusinessContext) error {
	amount, err := bc.StarBalance(c)
	if err != nil {
		return report(c, "get balance", err)
	}

	return reply(c, fmt.Sprintf("balance: %d stars (%d nanostars)", amount.Amount, amount.NanostarAmount))
}

func sendAlbum(c *botapi.Context, bc *botapi.BusinessContext, chat botapi.ChatID) error {
	media := []botapi.InputMedia{
		&botapi.InputMediaPhoto{Media: upload("one.png", color.RGBA{R: 0x26, G: 0x46, B: 0x53, A: 0xff}), Caption: "one"},
		&botapi.InputMediaPhoto{Media: upload("two.png", color.RGBA{R: 0xf4, G: 0xa2, B: 0x61, A: 0xff}), Caption: "two"},
	}

	_, err := c.Bot.SendMediaGroup(c, chat, media, botapi.WithBusinessConnection(bc.ConnectionID()))

	return report(c, "send album", err)
}

// report replies with the outcome of an action so it is visible in the chat.
func report(c *botapi.Context, action string, err error) error {
	if err != nil {
		return reply(c, action+" failed: "+err.Error())
	}

	return reply(c, action+" ok")
}

func reply(c *botapi.Context, text string) error {
	_, err := c.Reply(text)

	return err
}

// samplePhoto returns a generated solid-color InputFile so the example needs no
// asset files on disk.
func samplePhoto(c color.Color) botapi.InputFile {
	return upload("photo.png", c)
}

func upload(name string, c color.Color) *botapi.InputFileUpload {
	img := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)

	var buf bytes.Buffer

	_ = png.Encode(&buf, img)

	return &botapi.InputFileUpload{Name: name, Bytes: buf.Bytes()}
}

func isPublic(args []string) bool { return len(args) > 0 && args[0] == "public" }

func first(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return args[0]
}

func rest(args []string) string {
	if len(args) < 2 {
		return ""
	}

	return strings.Join(args[1:], " ")
}

func atoi(s string) (int, error) {
	var n int

	_, err := fmt.Sscanf(s, "%d", &n)

	return n, err
}
