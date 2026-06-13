package botapi

import (
	"errors"
	"fmt"
	"time"

	"github.com/gotd/td/tgerr"
)

// ResponseParameters describes why a request failed and how to recover, mirroring
// the Bot API "ResponseParameters" object.
type ResponseParameters struct {
	// MigrateToChatID is the new chat identifier when the group was migrated to a
	// supergroup. Zero when not applicable.
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`
	// RetryAfter is the number of seconds to wait before repeating the request
	// when a flood-wait limit was hit. Zero when not applicable.
	RetryAfter int `json:"retry_after,omitempty"`
}

// Error is a Bot-API-shaped error. Methods return it (wrapped) so callers can
// branch on a stable Code/Description regardless of the underlying MTProto error.
//
// Use errors.As to extract it:
//
//	var apiErr *botapi.Error
//	if errors.As(err, &apiErr) && apiErr.Code == 403 { ... }
type Error struct {
	// Code is the Bot-API-compatible error code (e.g. 400, 403, 429).
	Code int
	// Description is a human-readable message.
	Description string
	// Parameters carries optional recovery hints (retry_after, migrate_to_chat_id).
	Parameters *ResponseParameters
	// err is the wrapped underlying error (typically *tgerr.Error), if any.
	err error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Description == "" {
		return fmt.Sprintf("botapi: error %d", e.Code)
	}
	return fmt.Sprintf("botapi: %d: %s", e.Code, e.Description)
}

// Unwrap returns the underlying error, enabling errors.Is/errors.As to reach it.
func (e *Error) Unwrap() error { return e.err }

// ErrNotImplemented is returned by methods whose translation is not yet wired.
var ErrNotImplemented = errors.New("botapi: not implemented")

// AsFloodWait reports whether err (or anything it wraps) is a flood-wait error
// and, if so, how long the caller should wait before retrying.
//
// It understands both this package's *Error (via Parameters.RetryAfter) and the
// underlying github.com/gotd/td/tgerr flood-wait representation.
func AsFloodWait(err error) (retryAfter time.Duration, ok bool) {
	var apiErr *Error
	if errors.As(err, &apiErr) && apiErr.Parameters != nil && apiErr.Parameters.RetryAfter > 0 {
		return time.Duration(apiErr.Parameters.RetryAfter) * time.Second, true
	}
	return tgerr.AsFloodWait(err)
}

// AsChatMigrated reports whether err indicates that a basic group was upgraded
// to a supergroup and, if so, the new supergroup chat id the caller should use.
func AsChatMigrated(err error) (newChatID int64, ok bool) {
	var apiErr *Error
	if errors.As(err, &apiErr) && apiErr.Parameters != nil && apiErr.Parameters.MigrateToChatID != 0 {
		return apiErr.Parameters.MigrateToChatID, true
	}
	return 0, false
}

// Code returns the Bot API error code carried by err, or 0 if err is not a
// *Error. It is a convenience over errors.As for the common branch-on-code case.
func Code(err error) int {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.Code
	}
	return 0
}
