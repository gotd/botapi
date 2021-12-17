package botapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"go.uber.org/zap"

	"github.com/gotd/botapi/internal/oas"
)

// NotImplementedError is stub error for not implemented methods.
type NotImplementedError struct {
	Message string
}

// Error implements error.
func (n *NotImplementedError) Error() string {
	if n.Message == "" {
		return "method not implemented yet"
	}
	return n.Message
}

// PeerNotFoundError reports that BotAPI cannot find this peer.
type PeerNotFoundError struct {
	ID oas.ID
}

// Error implements error.
func (p *PeerNotFoundError) Error() string {
	if p.ID.IsString() {
		return fmt.Sprintf("peer %q not found", p.ID.String)
	}
	return fmt.Sprintf("peer %d not found", p.ID.Int64)
}

// BadRequestError reports bad request.
type BadRequestError struct {
	Message string
}

// Error implements error.
func (p *BadRequestError) Error() string {
	return p.Message
}

func errorOf(code int) oas.ErrorStatusCode {
	return oas.ErrorStatusCode{
		StatusCode: code,
		Response: oas.Error{
			ErrorCode:   code,
			Description: http.StatusText(code),
		},
	}
}

// NewError maps error to status code.
func (b *BotAPI) NewError(ctx context.Context, err error) oas.ErrorStatusCode {
	// TODO(tdakkota): pass request context info.
	b.logger.Warn("Request error", zap.Error(err))

	var (
		notImplemented *NotImplementedError
		peerNotFound   *PeerNotFoundError
		badRequest     *BadRequestError
	)
	// TODO(tdakkota): better error mapping.
	switch {
	case errors.As(err, &notImplemented):
		return errorOf(http.StatusNotImplemented)
	case errors.As(err, &peerNotFound):
		return errorOf(http.StatusNotFound)
	case errors.As(err, &badRequest):
		return errorOf(http.StatusBadRequest)
	}

	resp := errorOf(http.StatusInternalServerError)
	if b.debug && err != nil {
		resp.Response.Description = err.Error()
	}
	return resp
}

// NotFound is default not found handler.
func NotFound(w http.ResponseWriter, _ *http.Request) {
	apiError := errorOf(http.StatusNotFound)

	e := jx.GetEncoder()
	defer jx.PutEncoder(e)

	apiError.Encode(e)
	_, _ = e.WriteTo(w)
}
