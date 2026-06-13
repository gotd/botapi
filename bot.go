package botapi

import (
	"context"
	"strings"
	"sync"

	"github.com/go-faster/errors"
	"github.com/gotd/log"
	"github.com/gotd/log/logzap"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
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

	router router

	onStart func(ctx context.Context)

	// commands collects the (command, description) pairs registered through
	// OnCommand so Run can publish them to Telegram. registerCommands gates the
	// publish.
	commandsMu       sync.Mutex
	commands         []BotCommand
	registerCommands bool

	// runMu guards runCtx, the bot's run-lifetime context, used for background
	// (proactive) sends that must outlive a single update handler.
	runMu  sync.Mutex
	runCtx context.Context

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
		Middlewares:    buildMiddlewares(opt, h),
	})
	*rawPlaceholder = *client.API()

	b := &Bot{
		token:            token,
		log:              opt.Logger,
		client:           client,
		raw:              client.API(),
		sender:           message.NewSender(client.API()),
		peers:            pm,
		gaps:             gaps,
		disp:             disp,
		onStart:          opt.OnStart,
		registerCommands: !opt.DisableCommandRegistration,
	}
	b.installHandlers()
	return b, nil
}

// registerCommand records a command registered through OnCommand so Run can
// publish the set to Telegram. Duplicate names keep their first description.
func (b *Bot) registerCommand(name, description string) {
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		return
	}
	b.commandsMu.Lock()
	defer b.commandsMu.Unlock()
	for _, c := range b.commands {
		if c.Command == name {
			return
		}
	}
	b.commands = append(b.commands, BotCommand{Command: name, Description: description})
}

// publishCommands reports the collected OnCommand set to Telegram (default
// scope). Failures are logged, not fatal, since a bad description should not
// stop the bot from serving.
func (b *Bot) publishCommands(ctx context.Context) {
	if !b.registerCommands {
		return
	}
	b.commandsMu.Lock()
	cmds := append([]BotCommand(nil), b.commands...)
	b.commandsMu.Unlock()
	if len(cmds) == 0 {
		return
	}
	if err := b.SetMyCommands(ctx, cmds); err != nil {
		b.log.Warn("Register bot commands", zap.Error(err))
		return
	}
	b.log.Debug("Registered bot commands", zap.Int("count", len(cmds)))
}

// Run connects, authorizes as a bot, and blocks serving updates until ctx is
// canceled or a fatal error occurs. Register handlers before calling Run.
func (b *Bot) Run(ctx context.Context) error {
	return b.client.Run(ctx, func(ctx context.Context) error {
		// Expose the connection-lifetime context for background sends, and clear
		// it on shutdown so background work stops with the bot.
		b.setRunCtx(ctx)
		defer b.setRunCtx(nil)

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
				b.publishCommands(ctx)
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

// Logger returns the bot's zap logger.
func (b *Bot) Logger() *zap.Logger { return b.log }

func (b *Bot) setRunCtx(ctx context.Context) {
	b.runMu.Lock()
	b.runCtx = ctx
	b.runMu.Unlock()
}

// Background returns a context tied to the bot's run lifetime, for proactive or
// background sends that are not in response to an update — e.g. from a timer,
// queue or goroutine. It is canceled when the bot stops.
//
// A handler's own context is per-update (and may carry a Timeout deadline), so
// it must not be used for work that outlives the handler. Use Background
// instead:
//
//	bot.OnCommand("remind", "Remind in a minute", func(c *botapi.Context) error {
//		chat, _ := c.Chat()
//		ctx := c.Bot.Background()
//		go func() {
//			time.Sleep(time.Minute)
//			c.Bot.SendMessage(ctx, chat, "⏰ reminder")
//		}()
//		return nil
//	})
//
// Before Run has connected (or after it has stopped) Background returns an
// already-canceled context, so background sends fail fast rather than block.
func (b *Bot) Background() context.Context {
	b.runMu.Lock()
	ctx := b.runCtx
	b.runMu.Unlock()
	if ctx == nil {
		return canceledContext
	}
	return ctx
}

// canceledContext is returned by Background before the bot is running.
var canceledContext = func() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}()
