package botapi

import (
	"context"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

// GetMessage fetches a single message from chat by its id and converts it to a
// Bot API Message. The returned Message exposes the original MTProto message
// through Raw, so callers can reach fields the typed surface does not cover.
//
// It issues the channel- or user-appropriate getMessages call, so chat must be
// a peer the bot can address (seen before or backed by stored peer data); a
// send-only PeerRef is rejected.
func (b *Bot) GetMessage(ctx context.Context, chat ChatID, messageID int) (*Message, error) {
	p, err := b.resolvePeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	ids := []tg.InputMessageClass{&tg.InputMessageID{ID: messageID}}

	var (
		res    tg.MessagesMessagesClass
		rpcErr error
	)

	if ch, ok := p.(peers.Channel); ok {
		res, rpcErr = b.raw.ChannelsGetMessages(ctx, &tg.ChannelsGetMessagesRequest{
			Channel: ch.InputChannel(),
			ID:      ids,
		})
	} else {
		res, rpcErr = b.raw.MessagesGetMessages(ctx, ids)
	}

	if rpcErr != nil {
		return nil, asAPIError(rpcErr)
	}

	msgs, ok := res.AsModified()
	if !ok {
		return nil, &Error{Code: 400, Description: "Bad Request: message not found"}
	}

	for _, m := range msgs.GetMessages() {
		msg, ok := m.(*tg.Message)
		if !ok || msg.ID != messageID {
			continue
		}

		return b.convertMessage(ctx, msg)
	}

	return nil, &Error{Code: 400, Description: "Bad Request: message not found"}
}
