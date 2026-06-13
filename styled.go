package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message/entity"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

// styledMessage turns text + parse mode into a plain message string and its
// MTProto entities. It is the entity-producing counterpart of styledText (which
// targets the message builder) and is used where a raw (message, entities) pair
// is needed, such as inline-result message content.
func (b *Bot) styledMessage(ctx context.Context, text string, mode ParseMode) (string, []tg.MessageEntityClass, error) {
	if mode == ParseModeNone {
		return text, nil, nil
	}
	opts, err := styledText(text, mode, b.peers.UserResolveHook(ctx))
	if err != nil {
		return "", nil, err
	}
	var builder entity.Builder
	if err := styling.Perform(&builder, opts...); err != nil {
		return "", nil, &Error{Code: 400, Description: "Bad Request: can't parse entities: " + err.Error()}
	}
	msg, entities := builder.Complete()
	return msg, entities, nil
}
