package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// AnswerGuestQuery sets the result of a guest interaction with the bot and sends
// a corresponding message to the chat from which the query originated.
// guestQueryID is the query id received in the update. It returns the sent
// inline message id.
func (b *Bot) AnswerGuestQuery(ctx context.Context, guestQueryID string, result InlineQueryResult) (*SentWebAppMessage, error) {
	if result == nil {
		return nil, errNilInlineResult()
	}

	queryID, err := strconv.ParseInt(guestQueryID, 10, 64)
	if err != nil {
		return nil, &Error{Code: 400, Description: "Bad Request: invalid guest_query_id"}
	}

	converted, err := result.toTg(ctx, b)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.MessagesSetBotGuestChatResult(ctx, &tg.MessagesSetBotGuestChatResultRequest{
		QueryID: queryID,
		Result:  converted,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	enc, err := encodeInlineMessageID(res)
	if err != nil {
		return nil, err
	}

	return &SentWebAppMessage{InlineMessageID: enc}, nil
}
