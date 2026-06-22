package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestAccessors(t *testing.T) {
	b := newMockBot(newMockInvoker())
	if b.Raw() == nil || b.Sender() == nil || b.Peers() == nil {
		t.Fatal("nil accessor")
	}

	if b.Dispatcher() == nil {
		t.Fatal("nil dispatcher")
	}

	if b.Self() == nil || b.Self().ID != 1 {
		t.Fatalf("self = %#v", b.Self())
	}

	if b.Logger() == nil {
		t.Fatal("nil logger")
	}
}

func TestContextSend(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMessageRequestTypeID, messageUpdates(&tg.Message{
		ID: 1, Message: "hi", PeerID: &tg.PeerUser{UserID: 10},
	}))

	b := newMockBot(inv)
	c := &Context{Context: context.Background(), Bot: b, Update: &Update{
		Message: &Message{MessageID: 5, Chat: Chat{ID: 10, Type: ChatTypePrivate}},
	}}

	if _, err := c.Send("hi"); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if _, err := c.Reply("re"); err != nil {
		t.Fatalf("Reply: %v", err)
	}

	if !inv.called(tg.MessagesSendMessageRequestTypeID) {
		t.Fatal("expected messages.sendMessage")
	}
}

func TestContextSendNoChat(t *testing.T) {
	b := newMockBot(newMockInvoker())
	c := &Context{Context: context.Background(), Bot: b, Update: &Update{}}

	if _, err := c.Send("x"); err == nil {
		t.Fatal("expected error with no chat")
	}

	if _, err := c.Reply("x"); err == nil {
		t.Fatal("expected error with no message")
	}
}

func TestContextAnswerCallback(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetBotCallbackAnswerRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)
	c := &Context{Context: context.Background(), Bot: b, Update: &Update{
		CallbackQuery: &CallbackQuery{ID: "55"},
	}}

	if err := c.AnswerCallback(WithCallbackText("ok")); err != nil {
		t.Fatalf("AnswerCallback: %v", err)
	}

	// Wrong update kind is rejected.
	noCb := &Context{Context: context.Background(), Bot: b, Update: &Update{}}
	if err := noCb.AnswerCallback(); err == nil {
		t.Fatal("expected error without callback query")
	}
}

func TestContextAnswerInline(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetInlineBotResultsRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)
	c := &Context{Context: context.Background(), Bot: b, Update: &Update{
		InlineQuery: &InlineQuery{ID: "77"},
	}}
	results := []InlineQueryResult{
		&InlineQueryResultArticle{ID: "1", Title: "A", InputMessageContent: &InputTextMessageContent{MessageText: "x"}},
	}

	if err := c.AnswerInline(results); err != nil {
		t.Fatalf("AnswerInline: %v", err)
	}

	noIq := &Context{Context: context.Background(), Bot: b, Update: &Update{}}
	if err := noIq.AnswerInline(nil); err == nil {
		t.Fatal("expected error without inline query")
	}
}

func TestOnRegistrarsRouteByKind(t *testing.T) {
	b := newMockBot(newMockInvoker())

	var fired string

	b.OnEditedMessage(func(c *Context) error { fired = "edited"; return nil })
	b.OnChannelPost(func(c *Context) error { fired = "channel"; return nil })
	b.OnShippingQuery(func(c *Context) error { fired = "shipping"; return nil })
	b.OnPreCheckoutQuery(func(c *Context) error { fired = "precheckout"; return nil })

	cases := []struct {
		u    *Update
		want string
	}{
		{&Update{EditedMessage: &Message{}}, "edited"},
		{&Update{ChannelPost: &Message{}}, "channel"},
		{&Update{ShippingQuery: &ShippingQuery{}}, "shipping"},
		{&Update{PreCheckoutQuery: &PreCheckoutQuery{}}, "precheckout"},
	}
	for _, c := range cases {
		fired = ""

		b.route(context.Background(), c.u)

		if fired != c.want {
			t.Fatalf("update %+v fired %q, want %q", c.u, fired, c.want)
		}
	}
}

func TestGroupOnCallbackQuery(t *testing.T) {
	b := newMockBot(newMockInvoker())
	g := b.Group()
	fired := false

	g.OnCallbackQuery(func(c *Context) error { fired = true; return nil })
	b.route(context.Background(), &Update{CallbackQuery: &CallbackQuery{}})

	if !fired {
		t.Fatal("group callback handler did not fire")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	b := newMockBot(newMockInvoker())
	fired := false

	b.Use(Logging())
	b.OnMessage(func(c *Context) error { fired = true; return nil })
	b.route(context.Background(), &Update{Message: &Message{}})

	if !fired {
		t.Fatal("handler under Logging middleware did not fire")
	}
}

func TestDispatchMessageDropsOutgoing(t *testing.T) {
	b := newMockBot(newMockInvoker())
	fired := false

	b.OnMessage(func(c *Context) error { fired = true; return nil })
	// An outgoing (bot's own) message must be dropped, not dispatched.
	b.dispatchMessage(context.Background(), tg.Entities{}, &tg.Message{Out: true, Message: "self"}, false)

	if fired {
		t.Fatal("outgoing message must not be dispatched")
	}
}
