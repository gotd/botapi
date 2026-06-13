package botapi

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/gotd/log"
	"github.com/gotd/log/logzap"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// Bot is a Telegram Bot API client implemented over MTProto.
//
// Construct it with New, register update handlers on its Dispatcher, then call
// Run to connect and serve. Bot is safe for concurrent use once running.
type Bot struct {
	token string
	log   *zap.Logger

	client *telegram.Client
	raw    *tg.Client
	sender *message.Sender
	peers  *peers.Manager
	gaps   *updates.Manager
	disp   tg.UpdateDispatcher

	onStart func(ctx context.Context)

	self *tg.User
}

// New constructs an unconnected Bot from a BotFather token. It performs no
// network I/O; call Run to connect, authorize and serve updates.
func New(token string, opt Options) (*Bot, error) {
	if token == "" {
		return nil, errors.New("token is required")
	}
	if opt.AppID == 0 || opt.AppHash == "" {
		return nil, errors.New("AppID and AppHash are required (see https://my.telegram.org)")
	}
	opt.setDefaults()

	// gotd/td logs through the github.com/gotd/log port; bridge the zap logger.
	lg := logzap.New(opt.Logger)

	disp := tg.NewUpdateDispatcher()

	// peers.Manager needs a *tg.Client, but the real one only exists after
	// telegram.NewClient, which in turn needs the update handler built from the
	// manager. Break the cycle with a placeholder we backfill below (same trick
	// the gotd examples use).
	rawPlaceholder := new(tg.Client)
	pm := peers.Options{
		Storage: opt.Storage,
		Cache:   opt.Storage,
		Logger:  log.Named(lg, "peers"),
	}.Build(rawPlaceholder)

	gaps := updates.New(updates.Config{
		Handler:      disp,
		Storage:      opt.Storage,
		AccessHasher: pm,
		Logger:       log.Named(lg, "gaps"),
		OnChannelTooLong: func(channelID int64) {
			opt.Logger.Warn("Channel too long", zap.Int64("channel_id", channelID))
		},
	})

	// Harvest access hashes from every update, then feed gap recovery.
	h := pm.UpdateHook(gaps)

	client := telegram.NewClient(opt.AppID, opt.AppHash, telegram.Options{
		Logger:         log.Named(lg, "client"),
		Device:         opt.Device,
		UpdateHandler:  h,
		SessionStorage: opt.Storage,
		Middlewares: []telegram.Middleware{
			updhook.UpdateHook(h.Handle),
		},
	})
	*rawPlaceholder = *client.API()

	return &Bot{
		token:   token,
		log:     opt.Logger,
		client:  client,
		raw:     client.API(),
		sender:  message.NewSender(client.API()),
		peers:   pm,
		gaps:    gaps,
		disp:    disp,
		onStart: opt.OnStart,
	}, nil
}

// Run connects, authorizes as a bot, and blocks serving updates until ctx is
// canceled or a fatal error occurs. Register handlers before calling Run.
func (b *Bot) Run(ctx context.Context) error {
	return b.client.Run(ctx, func(ctx context.Context) error {
		status, err := b.client.Auth().Status(ctx)
		if err != nil {
			return errors.Wrap(err, "auth status")
		}
		if !status.Authorized {
			if _, err := b.client.Auth().Bot(ctx, b.token); err != nil {
				return errors.Wrap(err, "bot login")
			}
		}

		me, err := b.client.Self(ctx)
		if err != nil {
			return errors.Wrap(err, "get self")
		}
		if !me.Bot {
			return errors.New("authorized account is not a bot")
		}
		b.self = me

		if err := b.peers.Init(ctx); err != nil {
			return errors.Wrap(err, "init peers")
		}

		return b.gaps.Run(ctx, b.raw, me.ID, updates.AuthOptions{
			IsBot: true,
			OnStart: func(ctx context.Context) {
				b.log.Info("Bot started",
					zap.Int64("id", me.ID),
					zap.String("username", me.Username),
				)
				if b.onStart != nil {
					b.onStart(ctx)
				}
			},
		})
	})
}

// Raw returns the underlying gotd/td API client for direct MTProto calls. It is
// the escape hatch for anything the Bot API surface does not (yet) cover.
func (b *Bot) Raw() *tg.Client { return b.raw }

// Dispatcher returns the update dispatcher for registering raw MTProto update
// handlers. The typed Bot API handler framework is built on top of this and
// will be added in a later phase.
func (b *Bot) Dispatcher() *tg.UpdateDispatcher { return &b.disp }

// Sender returns the message sender used for outgoing messages.
func (b *Bot) Sender() *message.Sender { return b.sender }

// Peers returns the peer manager (resolution and access-hash storage).
func (b *Bot) Peers() *peers.Manager { return b.peers }

// Self returns the bot's own user. It is nil until Run has authorized.
func (b *Bot) Self() *tg.User { return b.self }
