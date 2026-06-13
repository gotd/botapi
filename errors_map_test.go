package botapi

import (
	"context"
	"errors"
	"testing"

	"github.com/gotd/td/tgerr"
)

func TestAsAPIError(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		if asAPIError(nil) != nil {
			t.Fatal("nil should stay nil")
		}
	})
	t.Run("PassThroughNonRPC", func(t *testing.T) {
		if got := asAPIError(context.Canceled); !errors.Is(got, context.Canceled) {
			t.Fatalf("non-RPC error should pass through, got %v", got)
		}
	})
	t.Run("AlreadyAPIError", func(t *testing.T) {
		in := &Error{Code: 400, Description: "x"}
		if got := asAPIError(in); got != error(in) {
			t.Fatalf("existing *Error should be returned unchanged, got %v", got)
		}
	})
	t.Run("FloodWait", func(t *testing.T) {
		got := asAPIError(&tgerr.Error{Code: 420, Type: "FLOOD_WAIT", Argument: 12})
		var apiErr *Error
		if !errors.As(got, &apiErr) || apiErr.Code != 429 {
			t.Fatalf("want 429, got %#v", got)
		}
		if apiErr.Parameters == nil || apiErr.Parameters.RetryAfter != 12 {
			t.Fatalf("want retry_after=12, got %#v", apiErr.Parameters)
		}
	})
	t.Run("Forbidden", func(t *testing.T) {
		got := asAPIError(&tgerr.Error{Code: 400, Type: "USER_IS_BLOCKED"})
		var apiErr *Error
		if !errors.As(got, &apiErr) || apiErr.Code != 403 {
			t.Fatalf("USER_IS_BLOCKED should map to 403, got %#v", got)
		}
	})
	t.Run("ChatNotFound", func(t *testing.T) {
		got := asAPIError(&tgerr.Error{Code: 400, Type: "PEER_ID_INVALID"})
		var apiErr *Error
		if !errors.As(got, &apiErr) || apiErr.Code != 400 || apiErr.Description != "Bad Request: chat not found" {
			t.Fatalf("unexpected: %#v", got)
		}
	})
}
