// Package botapi contains Telegram Bot API handlers implementation.
package botapi

import (
	"context"
	"sync"

	"github.com/go-faster/errors"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/peer"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
	"github.com/gotd/botapi/internal/peers"
)

// BotAPI is Bot API implementation.
type BotAPI struct {
	client *telegram.Client
	raw    *tg.Client
	gaps   *updates.Manager

	sender   *message.Sender
	resolver peer.Resolver
	peers    peers.Storage

	self    *tg.User
	selfID  atomic.Int64
	selfMux sync.Mutex

	debug  bool
	logger *zap.Logger
}

// NewBotAPI creates new BotAPI instance.
func NewBotAPI(
	client *telegram.Client,
	gaps *updates.Manager,
	store peers.Storage,
	opts Options,
) *BotAPI {
	opts.setDefaults()

	raw := client.API()
	resolver := peer.SingleflightResolver(peer.Plain(raw))
	return &BotAPI{
		client:   client,
		raw:      raw,
		gaps:     gaps,
		sender:   message.NewSender(raw).WithResolver(resolver),
		resolver: resolver,
		peers:    store,
		debug:    opts.Debug,
		logger:   opts.Logger,
	}
}

// Init makes some initialization requests.
func (b *BotAPI) Init(ctx context.Context) error {
	me, err := b.client.Self(ctx)
	if err != nil {
		return errors.Wrap(err, "self")
	}

	if err := b.gaps.Auth(ctx, b.raw, me.ID, true, false); err != nil {
		return errors.Wrap(err, "init gaps")
	}

	b.updateSelf(me)
	return nil
}

func (b *BotAPI) updateSelf(user *tg.User) {
	b.selfMux.Lock()
	b.self = user
	b.selfID.Store(user.ID)
	b.selfMux.Unlock()
}

func (b *BotAPI) getSelf() *tg.User {
	b.selfMux.Lock()
	self := b.self
	b.selfMux.Unlock()
	return self
}

// Client returns *telegram.Client used by this instance of BotAPI.
func (b *BotAPI) Client() *telegram.Client {
	return b.client
}

// GetUpdates implements oas.Handler.
func (b *BotAPI) GetUpdates(ctx context.Context, req oas.OptGetUpdates) (oas.ResultArrayOfUpdate, error) {
	return oas.ResultArrayOfUpdate{}, &NotImplementedError{}
}

// SetPassportDataErrors implements oas.Handler.
func (b *BotAPI) SetPassportDataErrors(ctx context.Context, req oas.SetPassportDataErrors) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}
