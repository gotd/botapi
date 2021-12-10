// Package botapi contains Telegram Bot API handlers implementation.
package botapi

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message/peer"

	"github.com/gotd/botapi/internal/oas"
	"github.com/gotd/botapi/internal/peers"
)

// BotAPI is Bot API implementation.
type BotAPI struct {
	client   *telegram.Client
	resolver peer.Resolver
	peers    peers.Storage
	debug    bool
}

// NewBotAPI creates new BotAPI instance.
func NewBotAPI(
	client *telegram.Client,
	peers peers.Storage,
	debug bool,
) *BotAPI {
	return &BotAPI{
		client:   client,
		resolver: peer.SingleflightResolver(peer.Plain(client.API())),
		peers:    peers,
		debug:    debug,
	}
}

// Client returns *telegram.Client used by this instance of BotAPI.
func (b *BotAPI) Client() *telegram.Client {
	return b.client
}

// GetUpdates implements oas.Handler.
func (b *BotAPI) GetUpdates(ctx context.Context, req oas.GetUpdates) (oas.ResultArrayOfUpdate, error) {
	return oas.ResultArrayOfUpdate{}, &NotImplementedError{}
}

// SetPassportDataErrors implements oas.Handler.
func (b *BotAPI) SetPassportDataErrors(ctx context.Context, req oas.SetPassportDataErrors) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}
