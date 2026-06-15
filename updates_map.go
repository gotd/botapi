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

// inlineCallbackQueryFromTg builds a CallbackQuery from a callback on an inline
// message (one sent via an inline query result). Unlike a regular callback, it
// carries an inline_message_id instead of a Message.
func inlineCallbackQueryFromTg(e tg.Entities, u *tg.UpdateInlineBotCallbackQuery) *CallbackQuery {
	cq := &CallbackQuery{
		ID:            strconv.FormatInt(u.QueryID, 10),
		ChatInstance:  strconv.FormatInt(u.ChatInstance, 10),
		Data:          string(u.Data),
		GameShortName: u.GameShortName,
	}

	if id, err := encodeInlineMessageID(u.MsgID); err == nil {
		cq.InlineMessageID = id
	}

	if user, ok := e.Users[u.UserID]; ok {
		cq.From = userFromTgUser(user)
	} else {
		cq.From = User{ID: u.UserID}
	}

	return cq
}

// chosenInlineResultFromTg builds a ChosenInlineResult from a bot inline-send
// update (the user picked one of the bot's inline results).
func chosenInlineResultFromTg(e tg.Entities, u *tg.UpdateBotInlineSend) *ChosenInlineResult {
	r := &ChosenInlineResult{
		ResultID: u.ID,
		Query:    u.Query,
	}

	if id, ok := u.GetMsgID(); ok {
		if enc, err := encodeInlineMessageID(id); err == nil {
			r.InlineMessageID = enc
		}
	}

	if geo, ok := u.GetGeo(); ok {
		if g, ok := geo.(*tg.GeoPoint); ok {
			r.Location = &Location{
				Latitude:           g.Lat,
				Longitude:          g.Long,
				HorizontalAccuracy: float64(g.AccuracyRadius),
			}
		}
	}

	if user, ok := e.Users[u.UserID]; ok {
		r.From = userFromTgUser(user)
	} else {
		r.From = User{ID: u.UserID}
	}

	return r
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

// shippingQueryFromTg builds a ShippingQuery from a bot shipping-query update.
func shippingQueryFromTg(e tg.Entities, u *tg.UpdateBotShippingQuery) *ShippingQuery {
	q := &ShippingQuery{
		ID:              strconv.FormatInt(u.QueryID, 10),
		InvoicePayload:  string(u.Payload),
		ShippingAddress: shippingAddressFromTg(u.ShippingAddress),
	}
	if user, ok := e.Users[u.UserID]; ok {
		q.From = userFromTgUser(user)
	} else {
		q.From = User{ID: u.UserID}
	}

	return q
}

// preCheckoutQueryFromTg builds a PreCheckoutQuery from a bot pre-checkout
// update.
func preCheckoutQueryFromTg(e tg.Entities, u *tg.UpdateBotPrecheckoutQuery) *PreCheckoutQuery {
	q := &PreCheckoutQuery{
		ID:               strconv.FormatInt(u.QueryID, 10),
		Currency:         u.Currency,
		TotalAmount:      int(u.TotalAmount),
		InvoicePayload:   string(u.Payload),
		ShippingOptionID: u.ShippingOptionID,
	}
	if user, ok := e.Users[u.UserID]; ok {
		q.From = userFromTgUser(user)
	} else {
		q.From = User{ID: u.UserID}
	}

	if info, ok := u.GetInfo(); ok {
		q.OrderInfo = orderInfoFromTg(info)
	}

	return q
}

// shippingAddressFromTg converts an MTProto post address.
func shippingAddressFromTg(a tg.PostAddress) ShippingAddress {
	return ShippingAddress{
		CountryCode: a.CountryISO2,
		State:       a.State,
		City:        a.City,
		StreetLine1: a.StreetLine1,
		StreetLine2: a.StreetLine2,
		PostCode:    a.PostCode,
	}
}

// orderInfoFromTg converts MTProto requested payment info.
func orderInfoFromTg(info tg.PaymentRequestedInfo) *OrderInfo {
	out := &OrderInfo{
		Name:        info.Name,
		PhoneNumber: info.Phone,
		Email:       info.Email,
	}
	if addr, ok := info.GetShippingAddress(); ok {
		a := shippingAddressFromTg(addr)

		out.ShippingAddress = &a
	}

	return out
}
