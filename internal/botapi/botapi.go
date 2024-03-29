// Package botapi contains Telegram Bot API handlers implementation.
package botapi

import (
	"context"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

// BotAPI is Bot API implementation.
type BotAPI struct {
	raw  *tg.Client
	gaps *updates.Manager

	sender *message.Sender
	peers  *peers.Manager

	debug  bool
	logger *zap.Logger

	oas.UnimplementedHandler
}

// NewBotAPI creates new BotAPI instance.
func NewBotAPI(
	raw *tg.Client,
	gaps *updates.Manager,
	peer *peers.Manager,
	opts Options,
) *BotAPI {
	opts.setDefaults()

	return &BotAPI{
		raw:    raw,
		gaps:   gaps,
		sender: message.NewSender(raw),
		peers:  peer,
		debug:  opts.Debug,
		logger: opts.Logger,
	}
}

// Init makes some initialization requests.
func (b *BotAPI) Init(ctx context.Context) error {
	if err := b.peers.Init(ctx); err != nil {
		return errors.Wrap(err, "init peers")
	}

	return nil
}

// GetUpdates implements oas.Handler.
func (b *BotAPI) GetUpdates(ctx context.Context, req oas.OptGetUpdates) (*oas.ResultArrayOfUpdate, error) {
	return nil, &NotImplementedError{}
}

// SetPassportDataErrors implements oas.Handler.
func (b *BotAPI) SetPassportDataErrors(ctx context.Context, req *oas.SetPassportDataErrors) (*oas.Result, error) {
	return nil, &NotImplementedError{}
}
