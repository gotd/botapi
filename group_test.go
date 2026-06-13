package botapi

import (
	"context"
	"testing"
)

func TestGroupPredicatesGuardHandlers(t *testing.T) {
	b := newTestBot(t)
	var fired bool
	g := b.Group(ChatTypeIs(ChatTypeSupergroup))
	g.OnMessage(func(c *Context) error { fired = true; return nil })

	// Wrong chat type: group predicate fails, handler must not fire.
	b.route(context.Background(), &Update{Message: &Message{Chat: Chat{Type: ChatTypePrivate}}})
	if fired {
		t.Fatal("group predicate should have blocked the handler")
	}

	// Right chat type: handler fires.
	b.route(context.Background(), &Update{Message: &Message{Chat: Chat{Type: ChatTypeSupergroup}}})
	if !fired {
		t.Fatal("handler should fire when group predicate matches")
	}
}

func TestGroupMiddlewareScopedAndOrdered(t *testing.T) {
	b := newTestBot(t)
	var order []string
	b.Use(func(next Handler) Handler {
		return func(c *Context) error { order = append(order, "global"); return next(c) }
	})
	g := b.Group()
	g.Use(func(next Handler) Handler {
		return func(c *Context) error { order = append(order, "group"); return next(c) }
	})
	g.OnMessage(func(c *Context) error { order = append(order, "handler"); return nil })

	b.route(context.Background(), &Update{Message: &Message{}})

	want := []string{"global", "group", "handler"}
	for i := range want {
		if i >= len(order) || order[i] != want[i] {
			t.Fatalf("order = %v, want %v (global outermost, group inside)", order, want)
		}
	}
}
