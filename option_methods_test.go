package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

// TestOptionBuildersThroughMethods covers option constructors by passing them
// through the methods that consume them.
func TestOptionBuildersThroughMethods(t *testing.T) {
	t.Run("shipping", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.MessagesSetBotShippingResultsRequestTypeID, &tg.BoolTrue{})

		b := newMockBot(inv)

		err := b.AnswerShippingQuery(context.Background(), "1", true,
			WithShippingOptions(ShippingOption{ID: "x", Title: "Std", Prices: []LabeledPrice{{Label: "p", Amount: 1}}}))
		if err != nil {
			t.Fatalf("AnswerShippingQuery ok: %v", err)
		}

		inv2 := newMockInvoker()
		inv2.reply(tg.MessagesSetBotShippingResultsRequestTypeID, &tg.BoolTrue{})

		b2 := newMockBot(inv2)
		if err := b2.AnswerShippingQuery(context.Background(), "1", false, WithShippingError("nope")); err != nil {
			t.Fatalf("AnswerShippingQuery err: %v", err)
		}
	})

	t.Run("precheckout", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.MessagesSetBotPrecheckoutResultsRequestTypeID, &tg.BoolTrue{})

		b := newMockBot(inv)
		if err := b.AnswerPreCheckoutQuery(context.Background(), "1", false, WithPreCheckoutError("bad")); err != nil {
			t.Fatalf("AnswerPreCheckoutQuery: %v", err)
		}
	})

	t.Run("livelocation", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))

		b := newMockBot(inv)
		kb := InlineKeyboard([]InlineKeyboardButton{InlineButtonData("o", "d")})

		_, err := b.EditMessageLiveLocation(context.Background(), userRef(10, 20), 5, 1, 2,
			WithHeading(90), WithProximityAlertRadius(100), WithHorizontalAccuracy(10), WithLiveLocationMarkup(kb))
		if err != nil {
			t.Fatalf("EditMessageLiveLocation: %v", err)
		}
	})

	t.Run("invitelink", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.MessagesExportChatInviteRequestTypeID, exportedInvite())

		b := newMockBot(inv)

		_, err := b.CreateChatInviteLink(context.Background(), tdlibChannel(50),
			WithInviteLinkExpire(1700000000), WithInviteLinkJoinRequest())
		if err != nil {
			t.Fatalf("CreateChatInviteLink: %v", err)
		}
	})

	t.Run("inlineswitchpm", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.MessagesSetInlineBotResultsRequestTypeID, &tg.BoolTrue{})

		b := newMockBot(inv)

		err := b.AnswerInlineQuery(context.Background(), "1", nil, WithInlineSwitchPM("Login", "start"))
		if err != nil {
			t.Fatalf("AnswerInlineQuery switchpm: %v", err)
		}
	})
}
