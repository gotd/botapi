package botapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"

	"github.com/gotd/botapi/internal/oas"
)

// NotImplementedError is stub error for not implemented methods.
type NotImplementedError struct{}

// Error implements error.
func (n *NotImplementedError) Error() string {
	return "method not implemented yet"
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
func (b BotAPI) NewError(ctx context.Context, err error) oas.ErrorStatusCode {
	var (
		notImplemented *NotImplementedError
		peerNotFound   *PeerNotFoundError
	)
	switch {
	case errors.As(err, &notImplemented):
		return errorOf(http.StatusNotImplemented)
	case errors.As(err, &peerNotFound):
		return errorOf(http.StatusNotFound)
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
