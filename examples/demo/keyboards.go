package main

import (
	"strings"

	"github.com/gotd/botapi"
)

// registerKeyboards wires both keyboard flavors and their feedback channels:
//   - reply keyboards send plain text messages when a button is tapped;
//   - inline keyboards send callback queries the bot answers and acts on.
func registerKeyboards(bot *botapi.Bot) {
	bot.OnCommand("keyboard", "Show an inline keyboard with callbacks", func(c *botapi.Context) error {
		_, err := c.Reply("Pick a color:", botapi.WithReplyMarkup(colorKeyboard()))
		return err
	})

	bot.OnCommand("removekbd", "Remove the reply keyboard", func(c *botapi.Context) error {
		_, err := c.Reply("Reply keyboard removed.", botapi.WithReplyMarkup(botapi.RemoveKeyboard()))
		return err
	})

	// Reply-keyboard taps arrive as ordinary text messages, matched here by exact
	// text with the Or predicate combinator.
	bot.OnMessage(func(c *botapi.Context) error {
		_, err := c.Reply("You tapped: " + c.Message().Text)
		return err
	}, botapi.Or(
		botapi.TextEquals("📊 Poll"),
		botapi.TextEquals("🎲 Dice"),
		botapi.TextEquals("❓ Help"),
	))

	// Inline-keyboard callbacks all share the "demo:" data prefix.
	bot.OnCallbackQuery(func(c *botapi.Context) error {
		cb := c.Update.CallbackQuery
		choice := strings.TrimPrefix(cb.Data, "demo:")

		// Acknowledge the tap with a toast so the client stops its spinner.
		if err := c.AnswerCallback(botapi.WithCallbackText("You chose " + choice)); err != nil {
			return err
		}

		if cb.Message == nil {
			return nil
		}

		// Edit the originating message in place to reflect the choice.
		_, err := c.Bot.EditMessageText(c, botapi.ID(cb.Message.Chat.ID), cb.Message.MessageID,
			"You chose: "+choice, botapi.WithReplyMarkup(colorKeyboard()))

		return err
	}, botapi.CallbackPrefix("demo:"))
}

// mainReplyKeyboard is the persistent reply keyboard offered by /start. It mixes
// plain text buttons with the special contact- and location-request buttons.
func mainReplyKeyboard() *botapi.ReplyKeyboardMarkup {
	return &botapi.ReplyKeyboardMarkup{
		Keyboard: [][]botapi.KeyboardButton{
			botapi.Row(botapi.Button("📊 Poll"), botapi.Button("🎲 Dice"), botapi.Button("❓ Help")),
			botapi.Row(botapi.ButtonContact("📇 Share contact"), botapi.ButtonLocation("📍 Share location")),
		},
		ResizeKeyboard:        true,
		InputFieldPlaceholder: "Tap a button or type a command",
	}
}

// colorKeyboard is an inline keyboard mixing callback-data buttons with a URL
// button.
func colorKeyboard() *botapi.InlineKeyboardMarkup {
	return botapi.InlineKeyboard(
		botapi.InlineRow(
			botapi.InlineButtonData("🔴 Red", "demo:red"),
			botapi.InlineButtonData("🟢 Green", "demo:green"),
			botapi.InlineButtonData("🔵 Blue", "demo:blue"),
		),
		botapi.InlineRow(
			botapi.InlineButtonURL("Built on gotd/td", "https://github.com/gotd/td"),
		),
	)
}
