package main

import (
	"fmt"
	"time"

	glog "github.com/gotd/log"
	"github.com/gotd/td/telegram/message/rich"

	"github.com/gotd/botapi"
)

// registerCommands wires the slash commands. Each (name, description) pair is
// published to Telegram's command menu via SetMyCommands when the bot starts.
func registerCommands(bot *botapi.Bot) {
	bot.OnCommand("start", "Welcome message and reply keyboard", func(c *botapi.Context) error {
		_, err := c.Reply(
			"👋 Welcome to the gotd/botapi demo!\nType /help to see everything I can do.",
			botapi.WithReplyMarkup(mainReplyKeyboard()),
		)

		return err
	})

	bot.OnCommand("help", "List the available demos", func(c *botapi.Context) error {
		_, err := c.Reply(helpText,
			botapi.WithParseMode(botapi.ParseModeHTML),
			botapi.DisableWebPagePreview(),
		)

		return err
	})

	registerFormatting(bot)
	registerContent(bot)
	registerEditing(bot)
}

// registerFormatting shows the three text-formatting routes: parse modes and
// the structured rich-message API.
func registerFormatting(bot *botapi.Bot) {
	bot.OnCommand("html", "Send an HTML-formatted message", func(c *botapi.Context) error {
		_, err := c.Reply(
			`<b>bold</b>, <i>italic</i>, <code>code</code>, <a href="https://telegram.org">link</a>, <tg-spoiler>spoiler</tg-spoiler>`,
			botapi.WithParseMode(botapi.ParseModeHTML),
		)

		return err
	})

	bot.OnCommand("md", "Send a MarkdownV2-formatted message", func(c *botapi.Context) error {
		_, err := c.Reply(
			"*bold*, _italic_, `code`, ||spoiler||",
			botapi.WithParseMode(botapi.ParseModeMarkdownV2),
		)

		return err
	})

	// A rich message (Bot API 10.1): a tree of page blocks rather than a flat
	// string, built with the gotd/td rich package and sent via SendRichMessage.
	bot.OnCommand("rich", "Send a structured rich message", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		msg := rich.New(
			rich.Heading1(rich.Plain("Rich messages")),
			rich.Paragraph(rich.Concat(
				rich.Plain("They carry "),
				rich.Bold(rich.Plain("headings")),
				rich.Plain(", "),
				rich.Italic(rich.Plain("paragraphs")),
				rich.Plain(" and more."),
			)),
			rich.Preformatted(rich.Plain("bot.SendRichMessage(ctx, chat, msg)"), "go"),
		).Input()

		_, err := c.Bot.SendRichMessage(c, chat, msg)

		return err
	})
}

// registerContent covers the non-text content send methods.
func registerContent(bot *botapi.Bot) {
	bot.OnCommand("poll", "Send a poll", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.SendPoll(c, chat, "Best way to build a Telegram bot?",
			[]string{"MTProto (gotd)", "HTTP Bot API", "Both"})

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

	// Send options: a message with no notification, and one whose content can't
	// be forwarded or saved.
	bot.OnCommand("silent", "Send a message without a notification", func(c *botapi.Context) error {
		_, err := c.Reply("🔕 Sent silently.", botapi.Silent())
		return err
	})

	bot.OnCommand("protect", "Send a message that can't be forwarded", func(c *botapi.Context) error {
		_, err := c.Reply("🔒 This message has protected content.", botapi.ProtectContent())
		return err
	})

	// Background send: reply now, then send again from a goroutine that outlives
	// the handler. The handler's own context is canceled on return, so proactive
	// work must use the bot's run-lifetime Background context.
	bot.OnCommand("remind", "Send a reminder in 5 seconds (background send)", func(c *botapi.Context) error {
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
}

// registerEditing demonstrates the edit, forward and delete methods as a single
// self-narrating sequence.
func registerEditing(bot *botapi.Bot) {
	bot.OnCommand("edit", "Send a message, then edit its text and markup", func(c *botapi.Context) error {
		chat, _ := c.Chat()

		m, err := c.Bot.SendMessage(c, chat, "Step 0: this message will change…")
		if err != nil {
			return err
		}

		time.Sleep(time.Second)

		if _, err := c.Bot.EditMessageText(c, chat, m.MessageID, "Step 1: ✏️ text edited"); err != nil {
			return err
		}

		time.Sleep(time.Second)

		kb := botapi.InlineKeyboard(botapi.InlineRow(
			botapi.InlineButtonData("✅ markup edited", "demo:ok"),
		))

		_, err = c.Bot.EditMessageText(c, chat, m.MessageID,
			"Step 2: ✅ all edits succeeded", botapi.WithReplyMarkup(kb))

		return err
	})

	bot.OnCommand("forward", "Forward your command back to you", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.ForwardMessage(c, chat, chat, c.Message().MessageID)

		return err
	})

	// Send a throwaway message and delete it a moment later.
	bot.OnCommand("selfdestruct", "Send a message that deletes itself", func(c *botapi.Context) error {
		chat, _ := c.Chat()

		m, err := c.Bot.SendMessage(c, chat, "💣 This message self-destructs in 3 seconds…")
		if err != nil {
			return err
		}

		ctx := c.Background()

		go func() {
			time.Sleep(3 * time.Second)

			if err := c.Bot.DeleteMessage(ctx, chat, m.MessageID); err != nil {
				glog.For(c.Bot.Logger()).Warn(ctx, "self-destruct delete failed", glog.Error(err))
			}
		}()

		return nil
	})

	// PeerRef turns the current chat into a small, JSON-serializable reference you
	// can persist and later send to with botapi.Peer(ref) — no re-resolution.
	bot.OnCommand("ref", "Show this chat's serializable PeerRef", func(c *botapi.Context) error {
		chat, _ := c.Chat()

		ref, err := c.Bot.PeerRef(c, chat)
		if err != nil {
			return err
		}

		_, err = c.Reply(fmt.Sprintf("PeerRef{Kind: %q, ID: %d, AccessHash: %d}",
			ref.Kind, ref.ID, ref.AccessHash))

		return err
	})
}
