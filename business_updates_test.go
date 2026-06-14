package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestBusinessRouteByKind(t *testing.T) {
	b := newMockBot(newMockInvoker())

	var fired string

	b.OnBusinessMessage(func(c *Context) error { fired = "message"; return nil })
	b.OnEditedBusinessMessage(func(c *Context) error { fired = "edited"; return nil })
	b.OnBusinessConnection(func(c *Context) error { fired = "connection"; return nil })
	b.OnDeletedBusinessMessages(func(c *Context) error { fired = "deleted"; return nil })

	cases := []struct {
		u    *Update
		want string
	}{
		{&Update{BusinessMessage: &Message{}}, "message"},
		{&Update{EditedBusinessMessage: &Message{}}, "edited"},
		{&Update{BusinessConnection: &BusinessConnection{ID: "bc1"}}, "connection"},
		{&Update{DeletedBusinessMessages: &BusinessMessagesDeleted{BusinessConnectionID: "bc1"}}, "deleted"},
	}
	for _, c := range cases {
		fired = ""

		b.route(context.Background(), c.u)

		if fired != c.want {
			t.Fatalf("update %+v fired %q, want %q", c.u, fired, c.want)
		}
	}
}

func TestDispatchBusinessMessageKeepsOutgoing(t *testing.T) {
	b := newMockBot(newMockInvoker())

	var got *Message

	b.OnBusinessMessage(func(c *Context) error { got = c.BusinessMessage(); return nil })

	// Unlike a normal message, an outgoing business message is still delivered.
	msg := &tg.Message{ID: 7, Out: true, Message: "hi", PeerID: &tg.PeerUser{UserID: 10}}
	b.dispatchBusinessMessage(context.Background(), tg.Entities{}, "bc1", msg, false)

	if got == nil {
		t.Fatal("outgoing business message should be dispatched")
	}

	if got.BusinessConnectionID != "bc1" {
		t.Fatalf("connection id = %q", got.BusinessConnectionID)
	}

	if got.MessageID != 7 {
		t.Fatalf("message id = %d", got.MessageID)
	}
}

func TestDispatchBusinessMessageEdited(t *testing.T) {
	b := newMockBot(newMockInvoker())

	fired := false

	b.OnEditedBusinessMessage(func(c *Context) error { fired = true; return nil })
	b.dispatchBusinessMessage(context.Background(), tg.Entities{}, "bc1", &tg.Message{ID: 8, PeerID: &tg.PeerUser{UserID: 10}}, true)

	if !fired {
		t.Fatal("edited business message handler did not fire")
	}
}

func TestDispatchBusinessMessageDedup(t *testing.T) {
	b := newMockBot(newMockInvoker())

	var fires int

	b.OnBusinessMessage(func(c *Context) error { fires++; return nil })

	e := tg.Entities{Users: map[int64]*tg.User{10: {ID: 10, AccessHash: 1}}}
	msg := &tg.Message{ID: 50, PeerID: &tg.PeerUser{UserID: 10}}

	// A redelivered (identical) business message must be handled only once.
	b.dispatchBusinessMessage(context.Background(), e, "bc1", msg, false)
	b.dispatchBusinessMessage(context.Background(), e, "bc1", msg, false)

	if fires != 1 {
		t.Fatalf("handler fired %d times, want 1", fires)
	}

	// The same id on a different connection is a different message.
	b.dispatchBusinessMessage(context.Background(), e, "bc2", msg, false)

	if fires != 2 {
		t.Fatalf("handler fired %d times, want 2", fires)
	}
}

func TestBusinessDedupFresh(t *testing.T) {
	d := newBusinessDedup(2)

	if !d.fresh("a") || d.fresh("a") {
		t.Fatal("first 'a' fresh, second not")
	}

	// Fill past capacity to evict "a", which then reads as fresh again.
	d.fresh("b")
	d.fresh("c")

	if !d.fresh("a") {
		t.Fatal("evicted key should read fresh again")
	}
}

func TestContextBusinessFromMessage(t *testing.T) {
	b := newMockBot(newMockInvoker())
	c := &Context{Context: context.Background(), Bot: b, Update: &Update{
		BusinessMessage: &Message{MessageID: 1, BusinessConnectionID: "bc7"},
	}}

	bc, ok := c.Business()
	if !ok || bc.ConnectionID() != "bc7" {
		t.Fatalf("business = %#v, ok=%v", bc, ok)
	}

	if c.BusinessMessage() == nil {
		t.Fatal("BusinessMessage accessor returned nil")
	}
}

func TestContextBusinessFromConnection(t *testing.T) {
	b := newMockBot(newMockInvoker())
	c := &Context{Context: context.Background(), Bot: b, Update: &Update{
		BusinessConnection: &BusinessConnection{ID: "bc8"},
	}}

	bc, ok := c.Business()
	if !ok || bc.ConnectionID() != "bc8" {
		t.Fatalf("business = %#v, ok=%v", bc, ok)
	}
}

func TestContextBusinessAbsent(t *testing.T) {
	b := newMockBot(newMockInvoker())
	c := &Context{Context: context.Background(), Bot: b, Update: &Update{
		Message: &Message{MessageID: 1},
	}}

	if bc, ok := c.Business(); ok || bc != nil {
		t.Fatalf("expected no business context, got %#v", bc)
	}

	if c.BusinessMessage() != nil {
		t.Fatal("non-business update should have no business message")
	}
}
