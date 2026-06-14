package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// reactionToTg converts a Bot API reaction type to the MTProto representation.
//
// The switch over the sealed ReactionType union is exhaustive (gochecksumtype).
func reactionToTg(r ReactionType) (tg.ReactionClass, error) {
	switch v := r.(type) {
	case ReactionTypeEmoji:
		return &tg.ReactionEmoji{Emoticon: v.Emoji}, nil
	case ReactionTypeCustomEmoji:
		id, err := strconv.ParseInt(v.CustomEmojiID, 10, 64)
		if err != nil {
			return nil, errInvalidCustomEmojiID()
		}

		return &tg.ReactionCustomEmoji{DocumentID: id}, nil
	case ReactionTypePaid:
		return &tg.ReactionPaid{}, nil
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: invalid reaction type"}
	}
}

// ReactionOption configures a SetMessageReaction call.
type ReactionOption func(*reactionConfig)

type reactionConfig struct {
	big bool
}

// WithBigReaction shows the reaction with a big, animated effect.
func WithBigReaction() ReactionOption {
	return func(c *reactionConfig) { c.big = true }
}

// SetMessageReaction sets the bot's reactions on a message. Passing an empty
// reactions slice removes the bot's reaction. As a non-premium user a bot can
// set at most one reaction per message.
func (b *Bot) SetMessageReaction(ctx context.Context, chat ChatID, messageID int, reactions []ReactionType, opts ...ReactionOption) error {
	var cfg reactionConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	tgReactions := make([]tg.ReactionClass, 0, len(reactions))

	for _, r := range reactions {
		converted, err := reactionToTg(r)
		if err != nil {
			return err
		}

		tgReactions = append(tgReactions, converted)
	}

	if _, err := b.raw.MessagesSendReaction(ctx, &tg.MessagesSendReactionRequest{
		Big:      cfg.big,
		Peer:     peer,
		MsgID:    messageID,
		Reaction: tgReactions,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}
