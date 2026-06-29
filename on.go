package botapi

import (
	"context"

	"github.com/gotd/log"
	"github.com/gotd/td/tg"
)

// installHandlers wires the raw tg.UpdateDispatcher to the Bot API router. It is
// called once from New. Update-conversion failures are logged and swallowed so
// a single bad update never tears down the update stream.
func (b *Bot) installHandlers() {
	b.disp.OnNewMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
		b.dispatchMessage(ctx, e, u.Message, false)

		return nil
	})
	b.disp.OnEditMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateEditMessage) error {
		b.dispatchMessage(ctx, e, u.Message, true)

		return nil
	})
	b.disp.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateNewChannelMessage) error {
		b.dispatchMessage(ctx, e, u.Message, false)

		return nil
	})
	b.disp.OnEditChannelMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateEditChannelMessage) error {
		b.dispatchMessage(ctx, e, u.Message, true)

		return nil
	})
	b.disp.OnBotCallbackQuery(func(ctx context.Context, e tg.Entities, u *tg.UpdateBotCallbackQuery) error {
		cq := callbackQueryFromTg(e, u)

		chat, err := b.chatByPeer(ctx, u.Peer)
		if err == nil {
			cq.Message = &Message{
				MessageID: u.MsgID,
				Chat:      chat,
			}
		}

		b.route(ctx, &Update{CallbackQuery: callbackQueryFromTg(e, u)})

		return nil
	})
	b.disp.OnBotInlineQuery(func(ctx context.Context, e tg.Entities, u *tg.UpdateBotInlineQuery) error {
		b.route(ctx, &Update{InlineQuery: inlineQueryFromTg(e, u), Entities: e})

		return nil
	})
	b.disp.OnBotInlineSend(func(ctx context.Context, e tg.Entities, u *tg.UpdateBotInlineSend) error {
		b.route(ctx, &Update{ChosenInlineResult: chosenInlineResultFromTg(e, u), Entities: e})

		return nil
	})
	b.disp.OnInlineBotCallbackQuery(func(ctx context.Context, e tg.Entities, u *tg.UpdateInlineBotCallbackQuery) error {
		b.route(ctx, &Update{CallbackQuery: inlineCallbackQueryFromTg(e, u), Entities: e})

		return nil
	})
	b.disp.OnBotShippingQuery(func(ctx context.Context, e tg.Entities, u *tg.UpdateBotShippingQuery) error {
		b.route(ctx, &Update{ShippingQuery: shippingQueryFromTg(e, u), Entities: e})

		return nil
	})
	b.disp.OnBotPrecheckoutQuery(func(ctx context.Context, e tg.Entities, u *tg.UpdateBotPrecheckoutQuery) error {
		b.route(ctx, &Update{PreCheckoutQuery: preCheckoutQueryFromTg(e, u), Entities: e})

		return nil
	})
	b.installBusinessHandlers()
}

// dispatchMessage converts a message and routes it as the appropriate update
// field. Channel-broadcast messages become channel posts; everything else is a
// regular message. edited selects the edited_* fields.
func (b *Bot) dispatchMessage(ctx context.Context, e tg.Entities, msg tg.MessageClass, edited bool) {
	// Drop the bot's own outgoing messages. MTProto echoes them back on the
	// update stream, but the HTTP Bot API never delivers them — without this a
	// bot that replies to messages would answer its own replies in a loop.
	if tgm, ok := msg.(*tg.Message); ok && tgm.Out {
		return
	}

	m, err := b.messageFromTg(ctx, msg)
	if err != nil {
		b.logger().Error(ctx, "Convert message", log.Error(err))

		return
	}

	if m == nil {
		return
	}

	u := &Update{}

	u.Entities = e

	switch {
	case m.Chat.Type == ChatTypeChannel && edited:
		u.EditedChannelPost = m
	case m.Chat.Type == ChatTypeChannel:
		u.ChannelPost = m
	case edited:
		u.EditedMessage = m
	default:
		u.Message = m
	}

	b.route(ctx, u)
}

// Kind predicates select an update by which field it carries. They are shared
// by Bot.On* and Group.On*.
func hasMessage(c *Context) bool            { return c.Update.Message != nil }
func hasEditedMessage(c *Context) bool      { return c.Update.EditedMessage != nil }
func hasChannelPost(c *Context) bool        { return c.Update.ChannelPost != nil }
func hasCallbackQuery(c *Context) bool      { return c.Update.CallbackQuery != nil }
func hasInlineQuery(c *Context) bool        { return c.Update.InlineQuery != nil }
func hasChosenInlineResult(c *Context) bool { return c.Update.ChosenInlineResult != nil }
func hasShippingQuery(c *Context) bool      { return c.Update.ShippingQuery != nil }
func hasPreCheckoutQuery(c *Context) bool   { return c.Update.PreCheckoutQuery != nil }

// OnMessage registers a handler for new messages matching the predicates.
func (b *Bot) OnMessage(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasMessage, predicates)...)
}

// OnEditedMessage registers a handler for edited messages.
func (b *Bot) OnEditedMessage(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasEditedMessage, predicates)...)
}

// OnChannelPost registers a handler for new channel posts.
func (b *Bot) OnChannelPost(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasChannelPost, predicates)...)
}

// OnCallbackQuery registers a handler for callback queries from inline keyboards.
func (b *Bot) OnCallbackQuery(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasCallbackQuery, predicates)...)
}

// OnInlineQuery registers a handler for inline queries.
func (b *Bot) OnInlineQuery(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasInlineQuery, predicates)...)
}

// OnChosenInlineResult registers a handler for inline results chosen by users.
// Requires inline feedback to be enabled with BotFather.
func (b *Bot) OnChosenInlineResult(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasChosenInlineResult, predicates)...)
}

// OnShippingQuery registers a handler for shipping queries (flexible invoices).
func (b *Bot) OnShippingQuery(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasShippingQuery, predicates)...)
}

// OnPreCheckoutQuery registers a handler for pre-checkout queries.
func (b *Bot) OnPreCheckoutQuery(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasPreCheckoutQuery, predicates)...)
}

// prepend returns p followed by rest, without mutating rest.
func prepend(p Predicate, rest []Predicate) []Predicate {
	out := make([]Predicate, 0, len(rest)+1)

	out = append(out, p)
	out = append(out, rest...)

	return out
}
