package botapi

import (
	"context"

	"github.com/gotd/log"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
)

// Storage persists everything a bot needs across restarts: the MTProto session,
// peer access hashes and peer cache, and the update gap state.
//
// A single implementation must satisfy all of these; storage.BBoltStorage does.
// When Options.Storage is nil the bot keeps all of this in memory.
type Storage interface {
	telegram.SessionStorage
	peers.Storage
	peers.Cache
	updates.StateStorage
}

// Options configures a Bot.
type Options struct {
	// AppID and AppHash identify the MTProto application. Required: even bots
	// need an app identity to connect to MTProto (obtain one at
	// https://my.telegram.org). They are NOT the bot token.
	AppID   int
	AppHash string

	// Logger is the structured logger the bot writes to, via the
	// github.com/gotd/log port. Defaults to a no-op logger. Wrap a *zap.Logger
	// with github.com/gotd/log/logzap.New, or a *slog.Logger with logslog.New.
	Logger log.Logger

	// Device describes the client to Telegram. Optional.
	Device telegram.DeviceConfig

	// Storage persists session, peers and update state. Optional; in-memory if
	// nil (nothing survives a restart).
	Storage Storage

	// OnStart is called once, after the bot is authorized and update gap
	// recovery is live. Optional.
	OnStart func(ctx context.Context)

	// FloodWait enables transparent flood-wait handling: a request that hits a
	// FLOOD_WAIT limit is retried after sleeping for the indicated duration,
	// instead of failing with a 429 error. Off by default.
	FloodWait bool

	// MaxFloodWaitRetries bounds how many times a flood-waited request is retried
	// when FloodWait is enabled. Zero uses the underlying library default.
	MaxFloodWaitRetries int

	// RequestsPerSecond, when greater than zero, proactively rate-limits outgoing
	// MTProto requests to this many per second via a global token bucket. It is a
	// coarse guard against hitting Telegram's limits; off by default.
	RequestsPerSecond float64

	// RequestBurst is the token-bucket burst size for RequestsPerSecond. Defaults
	// to 1 when RequestsPerSecond is set.
	RequestBurst int

	// DisableCommandRegistration stops Run from publishing the commands
	// registered via OnCommand to Telegram (SetMyCommands, default scope). By
	// default the bot's command menu is kept in sync with its OnCommand handlers.
	DisableCommandRegistration bool

	// resolver, publicKeys and dcList override the MTProto endpoints the client
	// connects to. They are unexported test seams used to point a Bot at an
	// in-process tgtest server; production code reaches Telegram's real DCs.
	resolver   dcs.Resolver
	publicKeys []telegram.PublicKey
	dcList     dcs.List
}

func (o *Options) setDefaults() {
	o.Logger = log.OrNop(o.Logger)
}
