package botapi

import (
	"context"
	"testing"
)

func TestOnRegistrationRoutesByKind(t *testing.T) {
	b := newTestBot(t)
	var got string
	b.OnMessage(func(c *Context) error { got = "message"; return nil })
	b.OnCallbackQuery(func(c *Context) error { got = "callback"; return nil })
	b.OnInlineQuery(func(c *Context) error { got = "inline"; return nil })

	cases := []struct {
		update *Update
		want   string
	}{
		{&Update{Message: &Message{}}, "message"},
		{&Update{CallbackQuery: &CallbackQuery{}}, "callback"},
		{&Update{InlineQuery: &InlineQuery{}}, "inline"},
	}
	for _, c := range cases {
		got = ""
		b.route(context.Background(), c.update)
		if got != c.want {
			t.Fatalf("update %+v routed to %q, want %q", c.update, got, c.want)
		}
	}
}

func TestOnMessageDoesNotFireForCallback(t *testing.T) {
	b := newTestBot(t)
	fired := false
	b.OnMessage(func(c *Context) error { fired = true; return nil })
	b.route(context.Background(), &Update{CallbackQuery: &CallbackQuery{}})
	if fired {
		t.Fatal("OnMessage handler must not fire for a callback-only update")
	}
}
