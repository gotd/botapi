package botapi

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
	"go.uber.org/zap"
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

	// Logger is the zap logger. Defaults to a no-op logger.
	Logger *zap.Logger

	// Device describes the client to Telegram. Optional.
	Device telegram.DeviceConfig

	// Storage persists session, peers and update state. Optional; in-memory if
	// nil (nothing survives a restart).
	Storage Storage

	// OnStart is called once, after the bot is authorized and update gap
	// recovery is live. Optional.
	OnStart func(ctx context.Context)
}

func (o *Options) setDefaults() {
	if o.Logger == nil {
		o.Logger = zap.NewNop()
	}
}
