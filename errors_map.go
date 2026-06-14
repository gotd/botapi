package botapi

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gotd/td/tgerr"
)

// Shared Bot API error descriptions reused across the mapping and the methods.
const (
	descChatNotFound = "Bad Request: chat not found"
	descWrongFileID  = "Bad Request: wrong file identifier/HTTP URL specified"
	descUnauthorized = "Unauthorized"
	descInvalidFile  = "Bad Request: invalid file"

	descNotModified = "Bad Request: message is not modified: specified new message content " +
		"and reply markup are exactly the same as a current content and reply markup of the message"
	descPhotoExt = "Bad Request: Photo has unsupported extension. " +
		"Use one of .jpg, .jpeg, .gif, .png, .tif or .bmp"
)

// rpcMapping is the final Bot API code and description for a known MTProto RPC
// error type. Descriptions are the verbatim strings the official Bot API server
// returns (including their "Bad Request:"/"Forbidden:" prefix), so callers see
// behavior identical to the HTTP Bot API.
type rpcMapping struct {
	code int
	desc string
}

// rpcErrorMap translates MTProto RPC error types into Bot API errors. It mirrors
// the official telegram-bot-api server (Client.cpp fail_query_with_error and the
// TDLib send-error rewrites). Where the official description depends on the
// calling method (e.g. "message to edit not found"), a method-independent
// wording is used since asAPIError has no method context.
var rpcErrorMap = map[string]rpcMapping{
	// Chat / peer resolution.
	"PEER_ID_INVALID":       {http.StatusBadRequest, descChatNotFound},
	"CHAT_ID_INVALID":       {http.StatusBadRequest, descChatNotFound},
	"CHAT_NOT_FOUND":        {http.StatusBadRequest, descChatNotFound},
	"CHANNEL_INVALID":       {http.StatusBadRequest, descChatNotFound},
	"USERNAME_INVALID":      {http.StatusBadRequest, descChatNotFound},
	"USERNAME_NOT_OCCUPIED": {http.StatusBadRequest, descChatNotFound},
	"USER_ID_INVALID":       {http.StatusBadRequest, "Bad Request: user not found"},

	// Messages.
	"MESSAGE_ID_INVALID":       {http.StatusBadRequest, "Bad Request: message to be edited not found"},
	"MESSAGE_NOT_MODIFIED":     {http.StatusBadRequest, descNotModified},
	"MESSAGE_EMPTY":            {http.StatusBadRequest, "Bad Request: message text is empty"},
	"MESSAGE_TOO_LONG":         {http.StatusBadRequest, "Bad Request: message is too long"},
	"MEDIA_CAPTION_TOO_LONG":   {http.StatusBadRequest, "Bad Request: message caption is too long"},
	"MESSAGE_DELETE_FORBIDDEN": {http.StatusBadRequest, "Bad Request: message can't be deleted"},
	"MESSAGE_AUTHOR_REQUIRED":  {http.StatusBadRequest, "Bad Request: message can't be edited"},

	// Reply markup / buttons.
	"REPLY_MARKUP_TOO_LONG": {http.StatusBadRequest, "Bad Request: reply markup is too long"},
	"REPLY_MARKUP_INVALID":  {http.StatusBadRequest, "Bad Request: reply markup is invalid"},
	"BUTTON_URL_INVALID":    {http.StatusBadRequest, "Bad Request: button URL is invalid"},
	"BUTTON_DATA_INVALID":   {http.StatusBadRequest, "Bad Request: button data is invalid"},

	// URLs / web pages / media.
	"WC_CONVERT_URL_INVALID":   {http.StatusBadRequest, "Bad Request: Wrong HTTP URL specified"},
	"EXTERNAL_URL_INVALID":     {http.StatusBadRequest, "Bad Request: Wrong HTTP URL specified"},
	"WEBPAGE_CURL_FAILED":      {http.StatusBadRequest, "Bad Request: Failed to get HTTP URL content"},
	"WEBPAGE_MEDIA_EMPTY":      {http.StatusBadRequest, "Bad Request: Wrong type of the web page content"},
	"MEDIA_GROUPED_INVALID":    {http.StatusBadRequest, "Bad Request: Can't use the media of the specified type in the album"},
	"MEDIA_EMPTY":              {http.StatusBadRequest, "Bad Request: Wrong file identifier/HTTP URL specified"},
	"PHOTO_EXT_INVALID":        {http.StatusBadRequest, descPhotoExt},
	"PHOTO_INVALID_DIMENSIONS": {http.StatusBadRequest, "Bad Request: PHOTO_INVALID_DIMENSIONS"},
	"FILE_PARTS_INVALID":       {http.StatusBadRequest, "Bad Request: file is too big"},
	"FILE_REFERENCE_EXPIRED":   {http.StatusBadRequest, descWrongFileID},

	// Inline / callback queries.
	"QUERY_ID_INVALID": {http.StatusBadRequest, "Bad Request: query is too old and response timeout expired or query ID is invalid"},

	// Stickers.
	"PACK_SHORT_NAME_INVALID":  {http.StatusBadRequest, "Bad Request: invalid sticker set name is specified"},
	"PACK_SHORT_NAME_OCCUPIED": {http.StatusBadRequest, "Bad Request: sticker set name is already occupied"},
	"STICKER_EMOJI_INVALID":    {http.StatusBadRequest, "Bad Request: invalid sticker emojis"},
	"STICKERSET_INVALID":       {http.StatusBadRequest, "Bad Request: STICKERSET_INVALID"},

	// Chat description / about.
	"CHAT_ABOUT_NOT_MODIFIED": {http.StatusBadRequest, "Bad Request: chat description is not modified"},

	// Member management.
	"USER_ADMIN_INVALID":      {http.StatusBadRequest, "Bad Request: user is an administrator of the chat"},
	"CHAT_ADMIN_REQUIRED":     {http.StatusBadRequest, "Bad Request: not enough rights"},
	"USER_NOT_PARTICIPANT":    {http.StatusBadRequest, "Bad Request: user is not a member of the chat"},
	"USER_NOT_MUTUAL_CONTACT": {http.StatusBadRequest, "Bad Request: user is not a mutual contact"},
	"ADMINS_TOO_MUCH":         {http.StatusBadRequest, "Bad Request: there are too many administrators in the chat"},
	"USER_CHANNELS_TOO_MUCH":  {http.StatusBadRequest, "Bad Request: the user is a member of too many chats"},

	// Forbidden (403) cases.
	"INPUT_USER_DEACTIVATED":    {http.StatusForbidden, "Forbidden: user is deactivated"},
	"USER_IS_BLOCKED":           {http.StatusForbidden, "Forbidden: bot was blocked by the user"},
	"USER_IS_BOT":               {http.StatusForbidden, "Forbidden: bot can't send messages to bots"},
	"USER_DELETED":              {http.StatusForbidden, "Forbidden: user is deactivated"},
	"CHAT_WRITE_FORBIDDEN":      {http.StatusForbidden, "Forbidden: bot can't send messages to the chat"},
	"CHAT_SEND_MEDIA_FORBIDDEN": {http.StatusForbidden, "Forbidden: not enough rights to send media to the chat"},
	"CHANNEL_PRIVATE":           {http.StatusForbidden, "Forbidden: bot is not a member of the chat"},
	"BOT_GROUPS_BLOCKED":        {http.StatusForbidden, "Forbidden: bot can't be added to groups"},

	// Auth.
	"AUTH_KEY_UNREGISTERED": {http.StatusUnauthorized, descUnauthorized},
	"SESSION_REVOKED":       {http.StatusUnauthorized, descUnauthorized},
	"USER_DEACTIVATED":      {http.StatusUnauthorized, descUnauthorized},
}

// asAPIError normalizes an error from the MTProto layer into a Bot-API-shaped
// *Error. It is a no-op for nil, for errors that are already *Error, and for
// non-RPC errors (e.g. context cancellation or network failures), which are
// returned unchanged so callers can branch on them.
func asAPIError(err error) error {
	if err == nil {
		return nil
	}

	var already *Error

	if errors.As(err, &already) {
		return err
	}

	// Context cancellation/deadline always passes through unchanged (even if it
	// happens to be wrapped by an RPC error) so callers can errors.Is on it.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
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
			Description: "Too Many Requests: retry after " + strconv.Itoa(int(d/time.Second)),
			Parameters:  &ResponseParameters{RetryAfter: int(d / time.Second)},
			err:         err,
		}
	}

	if m, ok := rpcErrorMap[rpcErr.Type]; ok {
		return &Error{Code: m.code, Description: m.desc, err: err}
	}

	// Unknown error: apply the official server's normalization and prefix rules.
	code := normalizeErrorCode(rpcErr.Code, rpcErr.Message)

	return &Error{
		Code:        code,
		Description: prefixDescription(code, rpcErr.Message),
		err:         err,
	}
}

// normalizeErrorCode mirrors the official server's code remapping: codes below
// 400 (and 404) collapse to 400, and a 403 whose message is a bare SCREAMING_CASE
// RPC constant is downgraded to 400.
func normalizeErrorCode(code int, message string) int {
	if code < http.StatusBadRequest || code == http.StatusNotFound {
		return http.StatusBadRequest
	}

	if code == http.StatusForbidden && isScreamingCase(message) {
		return http.StatusBadRequest
	}

	return code
}

// prefixDescription prepends the Bot API status prefix to a raw error message,
// following the official casing rule: a SCREAMING_CASE constant is kept verbatim,
// otherwise the first letter is lowercased.
func prefixDescription(code int, message string) string {
	prefix := statusPrefix(code)
	if message == "" {
		return prefix
	}

	if hasPrefix(message, prefix) {
		return message
	}

	if len(message) >= 2 && (message[1] == '_' || isUpper(message[1])) {
		return prefix + ": " + message
	}

	return prefix + ": " + lowerFirst(message)
}

// statusPrefix returns the Bot API description prefix for a status code. Unknown
// codes use the 400 prefix (matching the official server's default branch).
func statusPrefix(code int) string {
	switch code {
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Bad Request"
	}
}

func isScreamingCase(s string) bool {
	if s == "" {
		return false
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') && c != '_' {
			return false
		}
	}

	return true
}

func isUpper(c byte) bool { return c >= 'A' && c <= 'Z' }

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func lowerFirst(s string) string {
	if s == "" || !isUpper(s[0]) {
		return s
	}

	b := []byte(s)

	b[0] += 'a' - 'A'

	return string(b)
}
