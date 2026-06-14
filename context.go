package botapi

import "context"

// Message returns the message the update carries (new/edited message or channel
// post), or nil.
func (c *Context) Message() *Message { return c.Update.EffectiveMessage() }

// Sender returns the user who produced the update (message sender or callback/
// inline-query sender), or nil when there is none.
func (c *Context) Sender() *User {
	switch {
	case c.Update.CallbackQuery != nil:
		return &c.Update.CallbackQuery.From
	case c.Update.InlineQuery != nil:
		return &c.Update.InlineQuery.From
	}

	if m := c.Message(); m != nil {
		return m.From
	}

	return nil
}

// Chat returns the target chat id of the update and whether one is present.
func (c *Context) Chat() (ChatID, bool) {
	if m := c.Message(); m != nil {
		return ID(m.Chat.ID), true
	}

	if cq := c.Update.CallbackQuery; cq != nil && cq.Message != nil {
		return ID(cq.Message.Chat.ID), true
	}

	return nil, false
}

// Send sends a text message to the update's chat.
func (c *Context) Send(text string, opts ...SendOption) (*Message, error) {
	chat, ok := c.Chat()
	if !ok {
		return nil, &Error{Code: 400, Description: "Bad Request: update has no chat to send to"}
	}

	return c.Bot.SendMessage(c, chat, text, opts...)
}

// Reply sends a text message to the update's chat as a reply to the incoming
// message.
func (c *Context) Reply(text string, opts ...SendOption) (*Message, error) {
	m := c.Message()
	if m == nil {
		return nil, &Error{Code: 400, Description: "Bad Request: update has no message to reply to"}
	}

	opts = append([]SendOption{ReplyTo(m.MessageID)}, opts...)

	return c.Bot.SendMessage(c, ID(m.Chat.ID), text, opts...)
}

// AnswerCallback answers the update's callback query. It is an error to call it
// when the update is not a callback query.
func (c *Context) AnswerCallback(opts ...AnswerCallbackQueryOption) error {
	cq := c.Update.CallbackQuery
	if cq == nil {
		return &Error{Code: 400, Description: "Bad Request: update has no callback query to answer"}
	}

	return c.Bot.AnswerCallbackQuery(c, cq.ID, opts...)
}

// Background returns a context tied to the bot's run lifetime, for sends that
// must outlive this handler (a timer, queue or goroutine). The handler's own
// context is per-update and may be canceled (e.g. by Timeout middleware) as soon
// as the handler returns, so do not capture it for background work.
//
// Send messages to any chat in the background with Bot.SendMessage and this
// context:
//
//	ctx := c.Background()
//	go func() { c.Bot.SendMessage(ctx, botapi.ID(other), "hi") }()
func (c *Context) Background() context.Context { return c.Bot.Background() }

// AnswerInline answers the update's inline query with the given results. It is
// an error to call it when the update is not an inline query.
func (c *Context) AnswerInline(results []InlineQueryResult, opts ...AnswerInlineQueryOption) error {
	iq := c.Update.InlineQuery
	if iq == nil {
		return &Error{Code: 400, Description: "Bad Request: update has no inline query to answer"}
	}

	return c.Bot.AnswerInlineQuery(c, iq.ID, results, opts...)
}
