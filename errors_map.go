package botapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/gotd/td/tgerr"
)

// asAPIError normalizes an error from the MTProto layer into a Bot-API-shaped
// *Error. It is a no-op for nil, for errors that are already *Error, and for
// non-RPC errors (e.g. context cancellation or network failures), which are
// returned unchanged so callers can branch on them.
//
// The mapping here covers the common cases; the full tgerr -> Bot API table is
// completed in Phase 6.
func asAPIError(err error) error {
	if err == nil {
		return nil
	}
	var already *Error
	if errors.As(err, &already) {
		return err
	}
	rpcErr, ok := tgerr.As(err)
	if !ok {
		return err
	}

	// Flood wait maps to 429 with a retry hint.
	if d, ok := tgerr.AsFloodWait(err); ok {
		return &Error{
			Code:        http.StatusTooManyRequests,
			Description: "Too Many Requests: retry later",
			Parameters:  &ResponseParameters{RetryAfter: int(d / time.Second)},
			err:         err,
		}
	}

	out := &Error{Code: rpcErr.Code, Description: rpcErr.Message, err: err}
	switch rpcErr.Type {
	case "PEER_ID_INVALID", "CHAT_ID_INVALID", "CHAT_NOT_FOUND":
		out.Code, out.Description = http.StatusBadRequest, "Bad Request: chat not found"
	case "MESSAGE_ID_INVALID", "MESSAGE_NOT_FOUND":
		out.Code, out.Description = http.StatusBadRequest, "Bad Request: message not found"
	case "MESSAGE_NOT_MODIFIED":
		out.Code = http.StatusBadRequest
		out.Description = "Bad Request: message is not modified: specified new message content " +
			"and reply markup are exactly the same as a current content and reply markup of the message"
	case "MESSAGE_EMPTY":
		out.Code, out.Description = http.StatusBadRequest, "Bad Request: message text is empty"
	case "REPLY_MARKUP_TOO_LONG":
		out.Code, out.Description = http.StatusBadRequest, "Bad Request: reply markup is too long"
	case "USER_IS_BLOCKED":
		out.Code, out.Description = http.StatusForbidden, "Forbidden: bot was blocked by the user"
	case "INPUT_USER_DEACTIVATED":
		out.Code, out.Description = http.StatusForbidden, "Forbidden: user is deactivated"
	case "CHANNEL_PRIVATE":
		out.Code, out.Description = http.StatusForbidden, "Forbidden: bot is not a member of the chat"
	case "USER_ADMIN_INVALID":
		out.Code, out.Description = http.StatusBadRequest, "Bad Request: user is an administrator of the chat"
	default:
		if out.Description == "" {
			out.Description = http.StatusText(out.Code)
		}
	}
	return out
}
