package botapi

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gotd/td/tgerr"
)

func TestAsFloodWait(t *testing.T) {
	t.Run("FromError", func(t *testing.T) {
		err := &Error{Code: 429, Description: "Too Many Requests", Parameters: &ResponseParameters{RetryAfter: 5}}
		d, ok := AsFloodWait(err)

		if !ok || d != 5*time.Second {
			t.Fatalf("got (%v, %v), want (5s, true)", d, ok)
		}
	})
	t.Run("FromWrappedError", func(t *testing.T) {
		err := fmt.Errorf("send: %w", &Error{Code: 429, Parameters: &ResponseParameters{RetryAfter: 3}})
		if d, ok := AsFloodWait(err); !ok || d != 3*time.Second {
			t.Fatalf("got (%v, %v), want (3s, true)", d, ok)
		}
	})
	t.Run("FromTGErr", func(t *testing.T) {
		err := &tgerr.Error{Code: 420, Type: "FLOOD_WAIT", Argument: 7}
		if d, ok := AsFloodWait(err); !ok || d != 7*time.Second {
			t.Fatalf("got (%v, %v), want (7s, true)", d, ok)
		}
	})
	t.Run("NotFloodWait", func(t *testing.T) {
		if _, ok := AsFloodWait(errors.New("boom")); ok {
			t.Fatal("plain error should not be a flood wait")
		}
	})
}

func TestErrorUnwrap(t *testing.T) {
	sentinel := errors.New("underlying")
	err := &Error{Code: 400, Description: "Bad Request", err: sentinel}

	if !errors.Is(err, sentinel) {
		t.Fatal("errors.Is should reach the wrapped error")
	}

	var apiErr *Error

	if !errors.As(fmt.Errorf("ctx: %w", err), &apiErr) || apiErr.Code != 400 {
		t.Fatal("errors.As should extract *Error through a wrap")
	}
}
