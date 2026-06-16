package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// joinRequestQueryID parses a chat_join_request_query_id into the MTProto query
// id.
func joinRequestQueryID(id string) (int64, error) {
	v, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, &Error{Code: 400, Description: "Bad Request: invalid chat_join_request_query_id"}
	}

	return v, nil
}

// AnswerChatJoinRequestQuery responds to a chat join request the bot is asked to
// resolve. result must be one of "approve", "decline" or "queue" (queue defers
// the decision to the chat's administrators).
func (b *Bot) AnswerChatJoinRequestQuery(ctx context.Context, chatJoinRequestQueryID, result string) error {
	queryID, err := joinRequestQueryID(chatJoinRequestQueryID)
	if err != nil {
		return err
	}

	var res tg.JoinChatBotResultClass

	switch result {
	case "approve":
		res = &tg.JoinChatBotResultApproved{}
	case "decline":
		res = &tg.JoinChatBotResultDeclined{}
	case "queue":
		res = &tg.JoinChatBotResultQueued{}
	default:
		return &Error{Code: 400, Description: "Bad Request: result must be approve, decline or queue"}
	}

	if _, err := b.raw.BotsSetJoinChatResults(ctx, &tg.BotsSetJoinChatResultsRequest{
		QueryID: queryID,
		Result:  res,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// SendChatJoinRequestWebApp opens a Web App in response to a chat join request
// the bot is asked to resolve, so the user can complete a verification flow.
func (b *Bot) SendChatJoinRequestWebApp(ctx context.Context, chatJoinRequestQueryID, webAppURL string) error {
	queryID, err := joinRequestQueryID(chatJoinRequestQueryID)
	if err != nil {
		return err
	}

	if _, err := b.raw.BotsSetJoinChatResults(ctx, &tg.BotsSetJoinChatResultsRequest{
		QueryID: queryID,
		Result:  &tg.JoinChatBotResultWebView{URL: webAppURL},
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}
