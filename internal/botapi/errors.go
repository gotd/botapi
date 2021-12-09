package botapi

import (
	"context"
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
	)
	if errors.As(err, &notImplemented) {
		return errorOf(http.StatusNotImplemented)
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
