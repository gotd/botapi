package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// PreparedInlineMessage describes an inline message to be sent by a user of a Mini
// App. It is the result of SavePreparedInlineMessage.
type PreparedInlineMessage struct {
	// ID is the unique identifier of the prepared message.
	ID string `json:"id"`
	// ExpirationDate is the Unix time when the prepared message can no longer be
	// used.
	ExpirationDate int `json:"expiration_date"`
}

// PreparedInlineMessageOption configures the set of chat types a prepared inline
// message may be sent to.
type PreparedInlineMessageOption func(*preparedInlineConfig)

type preparedInlineConfig struct {
	allowUsers    bool
	allowBots     bool
	allowGroups   bool
	allowChannels bool
}

// WithAllowUserChats permits the message to be sent to private chats with users.
func WithAllowUserChats() PreparedInlineMessageOption {
	return func(c *preparedInlineConfig) { c.allowUsers = true }
}

// WithAllowBotChats permits the message to be sent to private chats with bots.
func WithAllowBotChats() PreparedInlineMessageOption {
	return func(c *preparedInlineConfig) { c.allowBots = true }
}

// WithAllowGroupChats permits the message to be sent to group and supergroup chats.
func WithAllowGroupChats() PreparedInlineMessageOption {
	return func(c *preparedInlineConfig) { c.allowGroups = true }
}

// WithAllowChannelChats permits the message to be sent to channel chats.
func WithAllowChannelChats() PreparedInlineMessageOption {
	return func(c *preparedInlineConfig) { c.allowChannels = true }
}

// peerTypes turns the allow-* flags into MTProto inline query peer types.
func (c preparedInlineConfig) peerTypes() []tg.InlineQueryPeerTypeClass {
	var out []tg.InlineQueryPeerTypeClass

	if c.allowUsers {
		out = append(out, &tg.InlineQueryPeerTypePM{})
	}

	if c.allowBots {
		out = append(out, &tg.InlineQueryPeerTypeBotPM{})
	}

	if c.allowGroups {
		out = append(out, &tg.InlineQueryPeerTypeChat{}, &tg.InlineQueryPeerTypeMegagroup{})
	}

	if c.allowChannels {
		out = append(out, &tg.InlineQueryPeerTypeBroadcast{})
	}

	return out
}

// SavePreparedInlineMessage stores a message that can later be sent by a user of a
// Mini App. By default the message may not be sent anywhere; use the WithAllow*
// options to permit specific chat types.
func (b *Bot) SavePreparedInlineMessage(
	ctx context.Context, userID int64, result InlineQueryResult, opts ...PreparedInlineMessageOption,
) (*PreparedInlineMessage, error) {
	if result == nil {
		return nil, errNilInlineResult()
	}

	var cfg preparedInlineConfig

	for _, o := range opts {
		o(&cfg)
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	converted, err := result.toTg(ctx, b)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.MessagesSavePreparedInlineMessage(ctx, &tg.MessagesSavePreparedInlineMessageRequest{
		Result:    converted,
		UserID:    user,
		PeerTypes: cfg.peerTypes(),
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	return &PreparedInlineMessage{ID: res.ID, ExpirationDate: res.ExpireDate}, nil
}
