// Package botapi contains Telegram Bot API handlers implementation.
package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
	"github.com/gotd/botapi/internal/pool"
)

// BotAPI is Bot API implementation.
type BotAPI struct {
	pool  *pool.Pool
	debug bool
}

// NewBotAPI creates new BotAPI.
func NewBotAPI(pool *pool.Pool, debug bool) *BotAPI {
	return &BotAPI{
		pool:  pool,
		debug: debug,
	}
}

// GetUpdates implements oas.Handler.
func (b *BotAPI) GetUpdates(ctx context.Context, req oas.GetUpdates) (oas.ResultArrayOfUpdate, error) {
	return oas.ResultArrayOfUpdate{}, &NotImplementedError{}
}

// SetPassportDataErrors implements oas.Handler.
func (b *BotAPI) SetPassportDataErrors(ctx context.Context, req oas.SetPassportDataErrors) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}
