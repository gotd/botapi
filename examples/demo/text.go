package main

import (
	"github.com/gotd/botapi"
)

// registerText wires free-text handlers that aren't commands or button taps:
// a greeting matched by a case-insensitive regex, an echo for any other text,
// and an edited-message notice. Handler registration order matters — the first
// matching handler wins — so the specific predicates come before the catch-all.
func registerText(bot *botapi.Bot) {
	// Greet on "hi" / "hello" / "hey".
	bot.OnMessage(func(c *botapi.Context) error {
		_, err := c.Reply("Hello, " + displayName(c.Sender()) + "! 👋")
		return err
	}, botapi.Regex(`(?i)^(hi|hello|hey)\b`))

	// Echo any other non-command text. Not(HasPrefix("/")) skips commands so they
	// fall through to their OnCommand handlers.
	bot.OnMessage(func(c *botapi.Context) error {
		_, err := c.Reply("You said: " + c.Message().Text)
		return err
	}, botapi.HasText(), botapi.Not(botapi.HasPrefix("/")))

	// Notice when a user edits one of their messages.
	bot.OnEditedMessage(func(c *botapi.Context) error {
		_, err := c.Reply("👀 I noticed you edited a message.")
		return err
	})
}
