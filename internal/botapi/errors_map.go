package botapi

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tgerr"

	"github.com/gotd/botapi/internal/oas"
)

func errorOf(code int) oas.ErrorStatusCode {
	return errorStatusCode(code, "")
}

func errorStatusCode(code int, description string) oas.ErrorStatusCode {
	if description == "" {
		description = http.StatusText(code)
	}
	return oas.ErrorStatusCode{
		StatusCode: code,
		Response: oas.Error{
			ErrorCode:   code,
			Description: description,
		},
	}
}

// See https://github.com/tdlib/telegram-bot-api/blob/90f52477814a2d8a08c9ffb1d780fd179815d715/telegram-bot-api/Client.cpp#L86.
func tryMapRPCError(err error) (r oas.ErrorStatusCode, _ bool) {
	rpcErr, ok := tgerr.As(err)
	if !ok || rpcErr.Code != 400 {
		return r, false
	}

	var (
		errorCode    = rpcErr.Code
		errorMessage = rpcErr.Message
	)
	switch rpcErr.Type {
	case "MESSAGE_NOT_MODIFIED":
		errorMessage = "message is not modified: specified new message content " +
			"and reply markup are exactly the same as a current content " +
			"and reply markup of the message"
	case "WC_CONVERT_URL_INVALID", "EXTERNAL_URL_INVALID":
		errorMessage = "Wrong HTTP URL specified"
	case "WEBPAGE_CURL_FAILED":
		errorMessage = "Failed to get HTTP URL content"
	case "WEBPAGE_MEDIA_EMPTY":
		errorMessage = "Wrong type of the web page content"
	case "MEDIA_GROUPED_INVALID":
		errorMessage = "Can't use the media of the specified type in the album"
	case "REPLY_MARKUP_TOO_LONG":
		errorMessage = "reply markup is too long"
	case "INPUT_USER_DEACTIVATED":
		errorCode = 403
		errorMessage = "Forbidden: user is deactivated"
	case "USER_IS_BLOCKED":
		errorCode = 403
		errorMessage = "bot was blocked by the user"
	case "USER_ADMIN_INVALID":
		errorCode = 400
		errorMessage = "user is an administrator of the chat"
	case "File generation failed":
		errorCode = 400
		errorMessage = "can't upload file by URL"
	case "CHAT_ABOUT_NOT_MODIFIED":
		errorCode = 400
		errorMessage = "chat description is not modified"
	case "PACK_SHORT_NAME_INVALID":
		errorCode = 400
		errorMessage = "invalid sticker set name is specified"
	case "PACK_SHORT_NAME_OCCUPIED":
		errorCode = 400
		errorMessage = "sticker set name is already occupied"
	case "STICKER_EMOJI_INVALID":
		errorCode = 400
		errorMessage = "invalid sticker emojis"
	case "QUERY_ID_INVALID":
		errorCode = 400
		errorMessage = "query is too old and response timeout expired or query ID is invalid"
	case "MESSAGE_DELETE_FORBIDDEN":
		errorCode = 400
		errorMessage = "message can't be deleted"
	}

	return errorStatusCode(errorCode, errorMessage), true
}

func mapGotdError(err error) error {
	var (
		peerNotFound *peers.PeerNotFoundError
	)
	switch {
	case errors.As(err, &peerNotFound):
		return &BadRequestError{Message: "Bad Request: chat not found"}
	}
	return err
}

// NewError maps error to status code.
func (b *BotAPI) NewError(ctx context.Context, err error) (r oas.ErrorStatusCode) {
	// TODO(tdakkota): pass request context info.
	defer func() {
		level := zap.DebugLevel
		if r.StatusCode >= 500 {
			level = zap.WarnLevel
		}
		if e := b.logger.Check(level, "Request error"); e != nil {
			e.Write(zap.Error(err))
		}
	}()

	err = mapGotdError(err)
	var (
		notImplemented *NotImplementedError
		badRequest     *BadRequestError
	)
	// TODO(tdakkota): better error mapping.
	switch {
	case errors.As(err, &notImplemented):
		return errorOf(http.StatusNotImplemented)
	case errors.As(err, &badRequest):
		return errorStatusCode(http.StatusBadRequest, badRequest.Message)
	}

	if r, ok := tryMapRPCError(err); ok {
		return r
	}

	resp := errorOf(http.StatusInternalServerError)
	if b.debug && err != nil {
		resp.Response.Description = err.Error()
	}
	return resp
}

var encodedNotFoundError = func() (r []byte) {
	e := jx.GetWriter()
	defer jx.PutWriter(e)

	errorOf(http.StatusNotFound).Encode(e)
	return append(r, e.Buf...)
}()

// NotFound is default not found handler.
func NotFound(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write(encodedNotFoundError)
}
