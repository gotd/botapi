// Command advanced is a full-featured demo bot that exercises most of the
// github.com/gotd/botapi surface: commands (auto-registered), formatting,
// inline and reply keyboards, callback queries, inline mode, the media send
// methods, incoming-media handling, editing, forwarding, chat actions, polls,
// dice, location/venue/contact and the predicate/middleware framework.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token (enable inline mode in @BotFather to test inline queries):
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/advanced
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	glog "github.com/gotd/log"
	"github.com/gotd/log/logzap"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram/message/rich"

	"github.com/gotd/botapi"
	"github.com/gotd/botapi/storage"
)

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
		AppID:     appID,
		AppHash:   os.Getenv("APP_HASH"),
		Logger:    logzap.New(log),
		Storage:   store,
		FloodWait: true, // transparently wait out flood limits
	})
	if err != nil {
		log.Fatal("Create bot", zap.Error(err))
	}

	// Global middleware: recover from panics, bound each handler, log outcomes.
	bot.Use(botapi.Recover(), botapi.Timeout(time.Minute), botapi.Logging())

	registerCommands(bot)
	registerKeyboards(bot)
	registerMedia(bot)
	registerMisc(bot)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting advanced bot")

	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}

// registerCommands wires the slash commands. Each description is published to
// the client command menu via SetMyCommands on start.
func registerCommands(bot *botapi.Bot) {
	bot.OnCommand("start", "Welcome message and reply keyboard", func(c *botapi.Context) error {
		_, err := c.Reply("👋 Welcome! Type /help to see everything I can do.",
			botapi.WithReplyMarkup(mainReplyKeyboard()))

		return err
	})

	bot.OnCommand("help", "List the available demos", func(c *botapi.Context) error {
		_, err := c.Reply(helpText, botapi.WithParseMode(botapi.ParseModeHTML), botapi.DisableWebPagePreview())
		return err
	})

	bot.OnCommand("html", "Send an HTML-formatted message", func(c *botapi.Context) error {
		_, err := c.Reply("<b>bold</b>, <i>italic</i>, <code>code</code>, <a href=\"https://telegram.org\">link</a>",
			botapi.WithParseMode(botapi.ParseModeHTML))

		return err
	})

	bot.OnCommand("md", "Send a MarkdownV2-formatted message", func(c *botapi.Context) error {
		_, err := c.Reply("*bold*, _italic_, `code`, ||spoiler||",
			botapi.WithParseMode(botapi.ParseModeMarkdownV2))

		return err
	})

	// A rich message (Bot API 10.1): structured page-block content built with
	// the gotd/td rich package.
	bot.OnCommand("rich", "Send a structured rich message", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		msg := rich.New(
			rich.Heading1(rich.Plain("Rich messages")),
			rich.Paragraph(rich.Concat(
				rich.Plain("They carry "),
				rich.Bold(rich.Plain("headings")),
				rich.Plain(", "),
				rich.Italic(rich.Plain("paragraphs")),
				rich.Plain(" and more — a tree of blocks, not a flat string."),
			)),
			rich.Preformatted(rich.Plain("bot.SendRichMessage(ctx, chat, msg)"), "go"),
		).Input()
		_, err := c.Bot.SendRichMessage(c, chat, msg)

		return err
	})

	bot.OnCommand("richhtml", "Send a rich message from HTML", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.SendRichHTML(c, chat, "<h1>From HTML</h1><p>Parsed by Telegram into a <b>rich</b> message.</p>")

		return err
	})

	bot.OnCommand("keyboard", "Show an inline keyboard with callbacks", func(c *botapi.Context) error {
		_, err := c.Reply("Pick one:", botapi.WithReplyMarkup(demoInlineKeyboard()))
		return err
	})

	bot.OnCommand("poll", "Send a poll", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.SendPoll(c, chat, "Best Go web framework?",
			[]string{"net/http", "gin", "echo", "fiber"})

		return err
	})

	bot.OnCommand("dice", "Roll a dice", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.SendDice(c, chat, botapi.DiceDie)

		return err
	})

	bot.OnCommand("location", "Send a location", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.SendLocation(c, chat, 55.7558, 37.6173)

		return err
	})

	bot.OnCommand("venue", "Send a venue", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.SendVenue(c, chat, 55.7520, 37.6175, "Red Square", "Moscow, Russia")

		return err
	})

	bot.OnCommand("contact", "Send a contact", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.SendContact(c, chat, "+1234567890", "Ada", "Lovelace")

		return err
	})

	bot.OnCommand("photo", "Send a photo by URL", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.SendPhoto(c, chat, botapi.FileURL(samplePhotoURL), "A cat 🐈 sent by URL")

		return err
	})

	bot.OnCommand("typing", "Show a chat action, then reply", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		if err := c.Bot.SendChatAction(c, chat, botapi.ChatActionTyping); err != nil {
			return err
		}

		time.Sleep(2 * time.Second)

		_, err := c.Reply("Done typing!")

		return err
	})

	// A self-verifying edit test: send a message, then edit its text and its
	// reply markup, reporting any failure back to the chat (and the log).
	bot.OnCommand("edit", "Send a message, then edit it (text + markup)", func(c *botapi.Context) error {
		chat, _ := c.Chat()

		m, err := c.Bot.SendMessage(c, chat, "Step 0: this message will change…")
		if err != nil {
			return err
		}

		time.Sleep(time.Second)

		if _, err := c.Bot.EditMessageText(c, chat, m.MessageID, "Step 1: ✏️ text edited"); err != nil {
			_, _ = c.Reply("EditMessageText failed: " + err.Error())
			return err
		}

		time.Sleep(time.Second)

		kb := botapi.InlineKeyboard([]botapi.InlineKeyboardButton{
			botapi.InlineButtonData("✅ markup edited", "demo:ok"),
		})
		if _, err := c.Bot.EditMessageReplyMarkup(c, chat, m.MessageID, kb); err != nil {
			_, _ = c.Reply("EditMessageReplyMarkup failed: " + err.Error())
			return err
		}

		time.Sleep(time.Second)

		if _, err := c.Bot.EditMessageText(c, chat, m.MessageID,
			"Step 2: ✅ all edits succeeded", botapi.WithReplyMarkup(kb)); err != nil {
			_, _ = c.Reply("final EditMessageText failed: " + err.Error())
			return err
		}

		return nil
	})

	bot.OnCommand("forward", "Forward your command back to you", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.ForwardMessage(c, chat, chat, c.Message().MessageID)

		return err
	})

	bot.OnCommand("silent", "Send a message without a notification", func(c *botapi.Context) error {
		_, err := c.Reply("🔕 Sent silently.", botapi.Silent())
		return err
	})

	// Background send: reply now, then send again from a goroutine that outlives
	// the handler. The handler's own context would be canceled on return, so use
	// the bot's run-lifetime Background context.
	bot.OnCommand("remind", "Send a reminder in 5 seconds (background)", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		ctx := c.Background()

		go func() {
			time.Sleep(5 * time.Second)

			if _, err := c.Bot.SendMessage(ctx, chat, "⏰ Here's your reminder!"); err != nil {
				glog.For(c.Bot.Logger()).Warn(ctx, "background reminder failed", glog.Error(err))
			}
		}()

		_, err := c.Reply("OK, I'll remind you in 5 seconds.")

		return err
	})

	bot.OnCommand("protect", "Send a message that can't be forwarded", func(c *botapi.Context) error {
		_, err := c.Reply("🔒 This message has protected content.", botapi.ProtectContent())
		return err
	})
}

// registerKeyboards handles reply-keyboard interactions and callback queries.
func registerKeyboards(bot *botapi.Bot) {
	bot.OnCommand("removekbd", "Remove the reply keyboard", func(c *botapi.Context) error {
		_, err := c.Reply("Keyboard removed.", botapi.WithReplyMarkup(&botapi.ReplyKeyboardRemove{RemoveKeyboard: true}))
		return err
	})

	// Reply-keyboard buttons arrive as plain text messages.
	bot.OnMessage(func(c *botapi.Context) error {
		_, err := c.Reply("You pressed: " + c.Message().Text)
		return err
	}, botapi.Or(botapi.TextEquals("📊 Poll"), botapi.TextEquals("🎲 Dice"), botapi.TextEquals("📍 Location")))

	// Inline-keyboard callbacks (prefix "demo:").
	bot.OnCallbackQuery(func(c *botapi.Context) error {
		data := strings.TrimPrefix(c.Update.CallbackQuery.Data, "demo:")
		if err := c.AnswerCallback(botapi.WithCallbackText("You chose " + data)); err != nil {
			return err
		}

		msg := c.Update.CallbackQuery.Message
		if msg == nil {
			return nil
		}

		_, err := c.Bot.EditMessageText(c, botapi.ID(msg.Chat.ID), msg.MessageID,
			"You chose: "+data, botapi.WithReplyMarkup(demoInlineKeyboard()))

		return err
	}, botapi.CallbackPrefix("demo:"))
}

// registerMedia echoes incoming media and reports details.
func registerMedia(bot *botapi.Bot) {
	bot.OnMessage(func(c *botapi.Context) error {
		photos := c.Message().Photo
		largest := photos[len(photos)-1]
		chat, _ := c.Chat()
		_, err := c.Bot.SendPhoto(c, chat, botapi.FileID(largest.FileID),
			fmt.Sprintf("Got your photo: %d×%d", largest.Width, largest.Height))

		return err
	}, hasPhoto)

	bot.OnMessage(func(c *botapi.Context) error {
		doc := c.Message().Document

		file, err := c.Bot.GetFile(c, doc.FileID)
		if err != nil {
			return err
		}

		_, err = c.Reply(fmt.Sprintf("Got document %q (%d bytes)\nfile_unique_id=%s",
			doc.FileName, doc.FileSize, file.FileUniqueID))

		return err
	}, hasDocument)

	bot.OnMessage(func(c *botapi.Context) error {
		s := c.Message().Sticker
		_, err := c.Reply("Nice sticker " + s.Emoji)

		return err
	}, hasSticker)

	bot.OnMessage(func(c *botapi.Context) error {
		loc := c.Message().Location
		_, err := c.Reply(fmt.Sprintf("You are at %.4f, %.4f", loc.Latitude, loc.Longitude))

		return err
	}, hasLocation)

	bot.OnMessage(func(c *botapi.Context) error {
		ct := c.Message().Contact
		_, err := c.Reply("Contact: " + ct.FirstName + " " + ct.PhoneNumber)

		return err
	}, hasContact)
}

// registerMisc wires inline mode, edited-message handling and a text predicate.
func registerMisc(bot *botapi.Bot) {
	// Inline mode: type "@yourbot something" in any chat.
	bot.OnInlineQuery(func(c *botapi.Context) error {
		q := strings.TrimSpace(c.Update.InlineQuery.Query)
		if q == "" {
			return c.AnswerInline(nil)
		}

		results := []botapi.InlineQueryResult{
			article("upper", "UPPERCASE", strings.ToUpper(q)),
			article("lower", "lowercase", strings.ToLower(q)),
			article("reverse", "Reversed", reverse(q)),
		}

		return c.AnswerInline(results, botapi.WithInlineCacheTime(1))
	})

	bot.OnEditedMessage(func(c *botapi.Context) error {
		_, err := c.Reply("👀 I noticed you edited a message.")
		return err
	})

	// Greet on "hi"/"hello"/"hey" (case-insensitive), but not commands.
	bot.OnMessage(func(c *botapi.Context) error {
		_, err := c.Reply("Hello, " + displayName(c.Sender()) + "! 👋")
		return err
	}, botapi.Regex(`(?i)^(hi|hello|hey)\b`))
}

// --- helpers ---

func mainReplyKeyboard() *botapi.ReplyKeyboardMarkup {
	return &botapi.ReplyKeyboardMarkup{
		Keyboard: [][]botapi.KeyboardButton{
			{botapi.Button("📊 Poll"), botapi.Button("🎲 Dice")},
			{botapi.ButtonLocation("📍 Location"), botapi.ButtonContact("📇 Contact")},
		},
		ResizeKeyboard: true,
	}
}

func demoInlineKeyboard() *botapi.InlineKeyboardMarkup {
	return botapi.InlineKeyboard(
		[]botapi.InlineKeyboardButton{
			botapi.InlineButtonData("🔴 Red", "demo:red"),
			botapi.InlineButtonData("🟢 Green", "demo:green"),
			botapi.InlineButtonData("🔵 Blue", "demo:blue"),
		},
		[]botapi.InlineKeyboardButton{
			botapi.InlineButtonURL("gotd/td", "https://github.com/gotd/td"),
		},
	)
}

func article(id, title, text string) botapi.InlineQueryResult {
	return &botapi.InlineQueryResultArticle{
		ID:                  id,
		Title:               title,
		Description:         text,
		InputMessageContent: &botapi.InputTextMessageContent{MessageText: text},
	}
}

func displayName(u *botapi.User) string {
	if u == nil {
		return "stranger"
	}

	if u.Username != "" {
		return "@" + u.Username
	}

	return u.FirstName
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return string(r)
}

func hasPhoto(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && len(m.Photo) > 0
}

func hasDocument(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && m.Document != nil
}

func hasSticker(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && m.Sticker != nil
}

func hasLocation(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && m.Location != nil
}

func hasContact(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && m.Contact != nil
}

const helpText = `<b>Advanced demo bot</b>

Commands:
/html — HTML formatting
/md — MarkdownV2 formatting
/rich — structured rich message (page blocks)
/richhtml — rich message from HTML
/keyboard — inline keyboard + callbacks
/removekbd — remove the reply keyboard
/photo — send a photo by URL
/poll — send a poll
/dice — roll a dice
/location — send a location
/venue — send a venue
/contact — send a contact
/typing — chat action
/edit — send then edit a message
/forward — forward your message back
/silent — silent message
/protect — protected-content message
/remind — background send after a delay

Also try:
• send me a <i>photo</i>, <i>document</i>, <i>sticker</i>, <i>location</i> or <i>contact</i>
• edit one of your messages
• type <code>@thisbot hello</code> in any chat (inline mode)
• say "hi"`
