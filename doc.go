// Package botapi is a Telegram Bot API library implemented over MTProto using
// github.com/gotd/td.
//
// Unlike HTTP Bot API clients, botapi does not talk to api.telegram.org. It
// exposes the familiar Bot API surface (types, methods, updates) but speaks
// MTProto directly. That avoids the public Bot API server's rate limits, keeps
// updates flowing over one persistent connection, and gives access to the raw
// gotd client when the Bot API surface does not cover something (Bot.Raw).
//
// # Construction
//
// New builds an unconnected bot from a BotFather token and an MTProto app
// identity (AppID/AppHash from https://my.telegram.org — these are not the bot
// token). Run connects, authorizes as a bot and serves updates until the
// context is canceled:
//
//	bot, err := botapi.New(token, botapi.Options{AppID: id, AppHash: hash})
//	if err != nil {
//		return err
//	}
//	bot.OnCommand("start", func(c *botapi.Context) error {
//		_, err := c.Reply("hello")
//		return err
//	})
//	return bot.Run(ctx)
//
// # Sending
//
// Outgoing methods hang off *Bot, take a context first and a ChatID target, and
// share functional SendOptions (ReplyTo, Silent, WithReplyMarkup,
// WithParseMode, ...). A ChatID is constructed with ID (numeric) or Username:
//
//	bot.SendMessage(ctx, botapi.ID(chatID), "*hi*", botapi.WithParseMode(botapi.ParseModeMarkdownV2))
//	bot.SendPhoto(ctx, botapi.Username("@channel"), botapi.FileURL("https://.../x.jpg"), "caption")
//
// # Receiving
//
// Updates are dispatched through a small handler framework. Register handlers
// with the On* helpers (OnMessage, OnCommand, OnCallbackQuery, OnInlineQuery,
// ...), narrow them with Predicates (Command, HasText, Regex, CallbackPrefix,
// And/Or/Not, ...), and group shared middleware with Group/Use. Each handler
// receives a *Context with helpers (Message, Sender, Reply, Send,
// AnswerCallback, AnswerInline).
//
// # Errors
//
// Methods return errors shaped like the HTTP Bot API: an *Error with a Code and
// Description (see errors_map.go). Extract it with errors.As, or use the
// helpers Code, AsFloodWait and AsChatMigrated. Context cancellation passes
// through unchanged so callers can errors.Is on ctx.Err().
//
// # Sealed unions
//
// Where the Bot API uses "one of" objects (ChatID, InputFile, ReplyMarkup,
// InputMedia, ChatMember, InlineQueryResult, InputMessageContent, ...) this
// package uses sealed interfaces: an interface with an unexported marker method
// and a fixed set of concrete implementations, so an illegal state is
// unrepresentable and switches over them are checked for exhaustiveness.
//
// See docs/roadmap.md for the implementation status.
package botapi
