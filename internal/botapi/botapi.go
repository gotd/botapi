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

	me, err := b.peers.Self(ctx)
	if err != nil {
		return errors.Wrap(err, "get self")
	}

	_, isBot := me.ToBot()
	if err := b.gaps.Auth(ctx, b.raw, me.ID(), isBot, false); err != nil {
		return errors.Wrap(err, "init gaps")
	}

	return nil
}

// GetUpdates implements oas.Handler.
func (b *BotAPI) GetUpdates(ctx context.Context, req oas.OptGetUpdates) (oas.ResultArrayOfUpdate, error) {
	return oas.ResultArrayOfUpdate{}, &NotImplementedError{}
}

// SetPassportDataErrors implements oas.Handler.
func (b *BotAPI) SetPassportDataErrors(ctx context.Context, req oas.SetPassportDataErrors) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}
