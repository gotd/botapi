package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// SentWebAppMessage describes an inline message sent by a Web App on behalf of a
// user.
type SentWebAppMessage struct {
	// InlineMessageID is the identifier of the sent inline message, present only
	// if the message has a reply markup with at least one callback or inline
	// query button.
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// AnswerWebAppQuery sets the result of an interaction with a Web App and sends a
// corresponding message on behalf of the user to the chat from which the query
// originated. webAppQueryID is the query id from the Web App.
func (b *Bot) AnswerWebAppQuery(ctx context.Context, webAppQueryID string, result InlineQueryResult) (*SentWebAppMessage, error) {
	if result == nil {
		return nil, &Error{Code: 400, Description: "Bad Request: inline query result is nil"}
	}

	converted, err := result.toTg(ctx, b)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.MessagesSendWebViewResultMessage(ctx, &tg.MessagesSendWebViewResultMessageRequest{
		BotQueryID: webAppQueryID,
		Result:     converted,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	out := &SentWebAppMessage{}

	if id, ok := res.GetMsgID(); ok {
		enc, err := encodeInlineMessageID(id)
		if err != nil {
			return nil, err
		}

		out.InlineMessageID = enc
	}

	return out, nil
}
