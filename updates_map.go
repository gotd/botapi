package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// messageFromTg converts a dispatched message into a Bot API Message. Service
// and empty messages yield (nil, nil) so callers can skip them.
func (b *Bot) messageFromTg(ctx context.Context, msg tg.MessageClass) (*Message, error) {
	m, ok := msg.(*tg.Message)
	if !ok {
		return nil, nil
	}
	return b.convertMessage(ctx, m)
}

// callbackQueryFromTg builds a CallbackQuery from a bot callback update. The
// sender is taken from the update entities (already harvested by the peer hook),
// so no extra RPC is needed.
func callbackQueryFromTg(e tg.Entities, u *tg.UpdateBotCallbackQuery) *CallbackQuery {
	cq := &CallbackQuery{
		ID:            strconv.FormatInt(u.QueryID, 10),
		ChatInstance:  strconv.FormatInt(u.ChatInstance, 10),
		Data:          string(u.Data),
		GameShortName: u.GameShortName,
	}
	if user, ok := e.Users[u.UserID]; ok {
		cq.From = userFromTgUser(user)
	} else {
		cq.From = User{ID: u.UserID}
	}
	return cq
}

// inlineQueryFromTg builds an InlineQuery from a bot inline-query update.
func inlineQueryFromTg(e tg.Entities, u *tg.UpdateBotInlineQuery) *InlineQuery {
	iq := &InlineQuery{
		ID:     strconv.FormatInt(u.QueryID, 10),
		Query:  u.Query,
		Offset: u.Offset,
	}
	if user, ok := e.Users[u.UserID]; ok {
		iq.From = userFromTgUser(user)
	} else {
		iq.From = User{ID: u.UserID}
	}
	return iq
}
