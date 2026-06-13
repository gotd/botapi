package botapi

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRecoverMiddleware(t *testing.T) {
	b := newTestBot(t)
	h := Recover()(func(c *Context) error { panic("boom") })

	err := h(&Context{Context: context.Background(), Bot: b, Update: &Update{}})
	var apiErr *Error
	if err == nil {
		t.Fatal("Recover must convert a panic into an error")
	}
	if !errors.As(err, &apiErr) || apiErr.Code != 500 {
		t.Fatalf("Recover should yield a 500 *Error, got %v", err)
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	b := newTestBot(t)
	var hasDeadline bool
	h := Timeout(time.Minute)(func(c *Context) error {
		_, hasDeadline = c.Deadline()
		return nil
	})

	if err := h(&Context{Context: context.Background(), Bot: b, Update: &Update{}}); err != nil {
		t.Fatal(err)
	}
	if !hasDeadline {
		t.Fatal("Timeout should give the handler context a deadline")
	}
}
