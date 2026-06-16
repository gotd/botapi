package botapi

import (
	"context"
	"math"

	"github.com/gotd/td/tg"
)

// GetUserPersonalChatMessages returns the most recent messages, up to limit,
// from the personal channel a user has linked to their profile. The messages are
// returned newest-first.
func (b *Bot) GetUserPersonalChatMessages(ctx context.Context, userID int64, limit int) ([]*Message, error) {
	if limit <= 0 {
		return nil, &Error{Code: 400, Description: "Bad Request: limit must be positive"}
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.MessagesGetPersonalChannelHistory(ctx, &tg.MessagesGetPersonalChannelHistoryRequest{
		UserID: user,
		Limit:  limit,
		MaxID:  math.MaxInt32,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	msgs, ok := res.AsModified()
	if !ok {
		return nil, nil
	}

	out := make([]*Message, 0, len(msgs.GetMessages()))

	for _, m := range msgs.GetMessages() {
		msg, ok := m.(*tg.Message)
		if !ok {
			continue
		}

		converted, err := b.convertMessage(ctx, msg)
		if err != nil {
			return nil, err
		}

		out = append(out, converted)
	}

	return out, nil
}
