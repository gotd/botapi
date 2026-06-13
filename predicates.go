package botapi

import (
	"regexp"
	"strings"
)

// EffectiveMessage returns the message carried by the update, regardless of
// whether it is a new/edited message or channel post. It is nil for updates
// that carry no message (e.g. callback or inline queries).
func (u *Update) EffectiveMessage() *Message {
	switch {
	case u.Message != nil:
		return u.Message
	case u.EditedMessage != nil:
		return u.EditedMessage
	case u.ChannelPost != nil:
		return u.ChannelPost
	case u.EditedChannelPost != nil:
		return u.EditedChannelPost
	default:
		return nil
	}
}

// Text returns the text of the effective message, or the callback query data,
// or empty when neither is present.
func (u *Update) Text() string {
	if m := u.EffectiveMessage(); m != nil {
		return m.Text
	}
	if u.CallbackQuery != nil {
		return u.CallbackQuery.Data
	}
	return ""
}

// commandName extracts the bot command name from message text: "/start@bot foo"
// yields ("start", true). Pure.
func commandName(text string) (string, bool) {
	if !strings.HasPrefix(text, "/") {
		return "", false
	}
	field := text
	if i := strings.IndexAny(text, " \t\n"); i >= 0 {
		field = text[:i]
	}
	field = field[1:] // drop leading slash
	if at := strings.IndexByte(field, '@'); at >= 0 {
		field = field[:at]
	}
	return field, field != ""
}

// Command matches a message whose first token is the given bot command (with or
// without a leading slash, and ignoring a trailing @botusername).
func Command(name string) Predicate {
	name = strings.TrimPrefix(name, "/")
	return func(u *Update) bool {
		m := u.EffectiveMessage()
		if m == nil {
			return false
		}
		got, ok := commandName(m.Text)
		return ok && got == name
	}
}

// HasPrefix matches a message whose text starts with prefix.
func HasPrefix(prefix string) Predicate {
	return func(u *Update) bool {
		m := u.EffectiveMessage()
		return m != nil && strings.HasPrefix(m.Text, prefix)
	}
}

// HasText matches any message that carries non-empty text.
func HasText() Predicate {
	return func(u *Update) bool {
		m := u.EffectiveMessage()
		return m != nil && m.Text != ""
	}
}

// TextEquals matches a message whose text equals s exactly.
func TextEquals(s string) Predicate {
	return func(u *Update) bool {
		m := u.EffectiveMessage()
		return m != nil && m.Text == s
	}
}

// Regex matches a message whose text matches the pattern. It panics if the
// pattern does not compile (a programming error caught at registration).
func Regex(pattern string) Predicate {
	re := regexp.MustCompile(pattern)
	return func(u *Update) bool {
		m := u.EffectiveMessage()
		return m != nil && re.MatchString(m.Text)
	}
}

// ChatTypeIs matches a message sent in a chat of the given type.
func ChatTypeIs(t ChatType) Predicate {
	return func(u *Update) bool {
		m := u.EffectiveMessage()
		return m != nil && m.Chat.Type == t
	}
}

// CallbackData matches a callback query whose data equals s.
func CallbackData(s string) Predicate {
	return func(u *Update) bool {
		return u.CallbackQuery != nil && u.CallbackQuery.Data == s
	}
}

// CallbackPrefix matches a callback query whose data starts with prefix.
func CallbackPrefix(prefix string) Predicate {
	return func(u *Update) bool {
		return u.CallbackQuery != nil && strings.HasPrefix(u.CallbackQuery.Data, prefix)
	}
}

// Not inverts a predicate.
func Not(p Predicate) Predicate {
	return func(u *Update) bool { return !p(u) }
}

// Or matches when any of the given predicates matches.
func Or(predicates ...Predicate) Predicate {
	return func(u *Update) bool {
		for _, p := range predicates {
			if p(u) {
				return true
			}
		}
		return false
	}
}

// OnCommand registers a handler for the given bot command (e.g. "start").
func (b *Bot) OnCommand(name string, h Handler, predicates ...Predicate) {
	b.OnMessage(h, prepend(Command(name), predicates)...)
}
