package botapi

import (
	"context"
	"errors"
	"fmt"
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
	t.Run("PassThroughDeadline", func(t *testing.T) {
		if got := asAPIError(context.DeadlineExceeded); !errors.Is(got, context.DeadlineExceeded) {
			t.Fatalf("deadline error should pass through, got %v", got)
		}
	})
	t.Run("ContextWrappedByRPCStillPasses", func(t *testing.T) {
		// An RPC-typed error that wraps a cancellation must still surface the
		// cancellation so callers can branch on ctx.Err().
		wrapped := fmt.Errorf("rpc: %w", context.Canceled)
		if got := asAPIError(wrapped); !errors.Is(got, context.Canceled) {
			t.Fatalf("wrapped cancellation should pass through, got %v", got)
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
		if apiErr.Description != "Too Many Requests: retry after 12" {
			t.Fatalf("description: %q", apiErr.Description)
		}
	})
	t.Run("Forbidden", func(t *testing.T) {
		got := asAPIError(&tgerr.Error{Code: 400, Type: "USER_IS_BLOCKED"})
		var apiErr *Error
		if !errors.As(got, &apiErr) || apiErr.Code != 403 {
			t.Fatalf("USER_IS_BLOCKED should map to 403, got %#v", got)
		}
		if apiErr.Description != "Forbidden: bot was blocked by the user" {
			t.Fatalf("description: %q", apiErr.Description)
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

func TestAsAPIErrorFallback(t *testing.T) {
	t.Run("ScreamingCaseKeptVerbatim", func(t *testing.T) {
		got := asAPIError(&tgerr.Error{Code: 400, Type: "SOME_NEW_ERROR", Message: "SOME_NEW_ERROR"})
		var apiErr *Error
		if !errors.As(got, &apiErr) {
			t.Fatal("want *Error")
		}
		if apiErr.Description != "Bad Request: SOME_NEW_ERROR" {
			t.Fatalf("description: %q", apiErr.Description)
		}
	})
	t.Run("Forbidden403AllCapsDowngradesTo400", func(t *testing.T) {
		got := asAPIError(&tgerr.Error{Code: 403, Type: "MYSTERY_FORBIDDEN", Message: "MYSTERY_FORBIDDEN"})
		if Code(got) != 400 {
			t.Fatalf("all-caps 403 should downgrade to 400, got %d", Code(got))
		}
	})
	t.Run("SubFourHundredCollapsesTo400", func(t *testing.T) {
		got := asAPIError(&tgerr.Error{Code: 303, Type: "SEE_OTHER", Message: "Something happened"})
		var apiErr *Error
		if !errors.As(got, &apiErr) || apiErr.Code != 400 {
			t.Fatalf("want 400, got %#v", got)
		}
		// Non-screaming message has its first letter lowercased.
		if apiErr.Description != "Bad Request: something happened" {
			t.Fatalf("description: %q", apiErr.Description)
		}
	})
	t.Run("Genuine403WithProseKept", func(t *testing.T) {
		got := asAPIError(&tgerr.Error{Code: 403, Type: "X", Message: "Have no rights"})
		var apiErr *Error
		if !errors.As(got, &apiErr) || apiErr.Code != 403 {
			t.Fatalf("prose 403 should stay 403, got %#v", got)
		}
		if apiErr.Description != "Forbidden: have no rights" {
			t.Fatalf("description: %q", apiErr.Description)
		}
	})
}

func TestErrorPredicates(t *testing.T) {
	t.Run("AsFloodWait", func(t *testing.T) {
		err := asAPIError(&tgerr.Error{Code: 420, Type: "FLOOD_WAIT", Argument: 7})
		d, ok := AsFloodWait(err)
		if !ok || d.Seconds() != 7 {
			t.Fatalf("AsFloodWait: %v %v", d, ok)
		}
	})
	t.Run("AsChatMigrated", func(t *testing.T) {
		err := &Error{Code: 400, Parameters: &ResponseParameters{MigrateToChatID: -100123}}
		id, ok := AsChatMigrated(err)
		if !ok || id != -100123 {
			t.Fatalf("AsChatMigrated: %v %v", id, ok)
		}
		if _, ok := AsChatMigrated(context.Canceled); ok {
			t.Fatal("non-migrate error should report false")
		}
	})
	t.Run("Code", func(t *testing.T) {
		if Code(&Error{Code: 403}) != 403 {
			t.Fatal("Code should extract 403")
		}
		if Code(context.Canceled) != 0 {
			t.Fatal("Code of non-*Error should be 0")
		}
	})
}
