package botapi

import (
	"context"
	"errors"
	"testing"
)

func TestRouterFirstMatchWins(t *testing.T) {
	b := newTestBot(t)

	var calls []string

	b.on(func(c *Context) error { calls = append(calls, "skipped"); return nil }, func(c *Context) bool { return false })
	b.on(func(c *Context) error { calls = append(calls, "matched"); return nil })
	b.on(func(c *Context) error { calls = append(calls, "second-match"); return nil })

	b.route(context.Background(), &Update{UpdateID: 1})

	if len(calls) != 1 || calls[0] != "matched" {
		t.Fatalf("expected only the first matching handler, got %v", calls)
	}
}

func TestMiddlewareOrder(t *testing.T) {
	b := newTestBot(t)

	var order []string

	b.Use(func(next Handler) Handler {
		return func(c *Context) error {
			order = append(order, "outer-in")
			defer func() { order = append(order, "outer-out") }()

			return next(c)
		}
	})
	b.Use(func(next Handler) Handler {
		return func(c *Context) error {
			order = append(order, "inner-in")
			defer func() { order = append(order, "inner-out") }()

			return next(c)
		}
	})
	b.on(func(c *Context) error { order = append(order, "handler"); return nil })

	b.route(context.Background(), &Update{})

	want := []string{"outer-in", "inner-in", "handler", "inner-out", "outer-out"}
	if len(order) != len(want) {
		t.Fatalf("order = %v, want %v", order, want)
	}

	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order = %v, want %v", order, want)
		}
	}
}

func TestRouterHandlerErrorIsContained(t *testing.T) {
	b := newTestBot(t)
	b.on(func(c *Context) error { return errors.New("boom") })
	// Must not panic or propagate; the error is logged.
	b.route(context.Background(), &Update{})
}
