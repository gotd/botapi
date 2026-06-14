package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSetChatTitle(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsEditTitleRequestTypeID, okUpdates())

	b := newMockBot(inv)

	if err := b.SetChatTitle(context.Background(), tdlibChannel(50), "New Title"); err != nil {
		t.Fatalf("SetChatTitle: %v", err)
	}

	var req tg.ChannelsEditTitleRequest

	inv.decode(t, tg.ChannelsEditTitleRequestTypeID, &req)

	if req.Title != "New Title" {
		t.Fatalf("title = %q", req.Title)
	}
}

func TestSetChatDescription(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditChatAboutRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	if err := b.SetChatDescription(context.Background(), tdlibChannel(50), "About"); err != nil {
		t.Fatalf("SetChatDescription: %v", err)
	}
}

func TestSetChatPermissions(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditChatDefaultBannedRightsRequestTypeID, okUpdates())

	b := newMockBot(inv)

	perms := ChatPermissions{CanSendMessages: true}
	if err := b.SetChatPermissions(context.Background(), tdlibChannel(50), perms); err != nil {
		t.Fatalf("SetChatPermissions: %v", err)
	}
}

func TestPinChatMessage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUpdatePinnedMessageRequestTypeID, okUpdates())

	b := newMockBot(inv)

	if err := b.PinChatMessage(context.Background(), userRef(10, 20), 7); err != nil {
		t.Fatalf("PinChatMessage: %v", err)
	}

	var req tg.MessagesUpdatePinnedMessageRequest

	inv.decode(t, tg.MessagesUpdatePinnedMessageRequestTypeID, &req)

	if req.ID != 7 {
		t.Fatalf("id = %d", req.ID)
	}
}

func TestUnpinChatMessage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUpdatePinnedMessageRequestTypeID, okUpdates())

	b := newMockBot(inv)

	if err := b.UnpinChatMessage(context.Background(), userRef(10, 20), 7); err != nil {
		t.Fatalf("UnpinChatMessage: %v", err)
	}

	var req tg.MessagesUpdatePinnedMessageRequest

	inv.decode(t, tg.MessagesUpdatePinnedMessageRequestTypeID, &req)

	if !req.Unpin {
		t.Fatal("expected Unpin")
	}
}

func TestUnpinAllChatMessages(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUnpinAllMessagesRequestTypeID, &tg.MessagesAffectedHistory{})

	b := newMockBot(inv)

	if err := b.UnpinAllChatMessages(context.Background(), userRef(10, 20)); err != nil {
		t.Fatalf("UnpinAllChatMessages: %v", err)
	}
}

func TestSetChatStickerSet(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsSetStickersRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	if err := b.SetChatStickerSet(context.Background(), tdlibChannel(50), "mypack"); err != nil {
		t.Fatalf("SetChatStickerSet: %v", err)
	}

	var req tg.ChannelsSetStickersRequest

	inv.decode(t, tg.ChannelsSetStickersRequestTypeID, &req)

	if s, ok := req.Stickerset.(*tg.InputStickerSetShortName); !ok || s.ShortName != "mypack" {
		t.Fatalf("stickerset = %#v", req.Stickerset)
	}
}

func TestDeleteChatStickerSet(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsSetStickersRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	if err := b.DeleteChatStickerSet(context.Background(), tdlibChannel(50)); err != nil {
		t.Fatalf("DeleteChatStickerSet: %v", err)
	}
}

func TestDeleteChatPhoto(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsEditPhotoRequestTypeID, okUpdates())

	b := newMockBot(inv)

	if err := b.DeleteChatPhoto(context.Background(), tdlibChannel(50)); err != nil {
		t.Fatalf("DeleteChatPhoto: %v", err)
	}
}

func TestGetStickerSet(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesGetStickerSetRequestTypeID, &tg.MessagesStickerSet{
		Set: tg.StickerSet{ShortName: "mypack", Title: "My Pack"},
	})

	b := newMockBot(inv)

	set, err := b.GetStickerSet(context.Background(), "mypack")
	if err != nil {
		t.Fatalf("GetStickerSet: %v", err)
	}

	if set.Name != "mypack" || set.Title != "My Pack" {
		t.Fatalf("set = %#v", set)
	}
}

func TestSendGame(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, messageUpdates(&tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}}))

	b := newMockBot(inv)

	if _, err := b.SendGame(context.Background(), userRef(10, 20), "mygame"); err != nil {
		t.Fatalf("SendGame: %v", err)
	}
}

func TestAnswerShippingQuery(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetBotShippingResultsRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	if err := b.AnswerShippingQuery(context.Background(), "12345", true); err != nil {
		t.Fatalf("AnswerShippingQuery: %v", err)
	}
}

func TestAnswerPreCheckoutQuery(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetBotPrecheckoutResultsRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	if err := b.AnswerPreCheckoutQuery(context.Background(), "12345", true); err != nil {
		t.Fatalf("AnswerPreCheckoutQuery: %v", err)
	}
}

func TestSetPassportDataErrors(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.UsersSetSecureValueErrorsRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	errs := []PassportElementError{
		&PassportElementErrorDataField{Type: "personal_details", FieldName: "first_name", DataHash: "aGFzaA==", Message: "bad"},
	}
	if err := b.SetPassportDataErrors(context.Background(), 99, errs); err != nil {
		t.Fatalf("SetPassportDataErrors: %v", err)
	}
}

func TestStopPoll(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsGetMessagesRequestTypeID, &tg.MessagesChannelMessages{
		Messages: []tg.MessageClass{&tg.Message{
			ID:     7,
			PeerID: &tg.PeerChannel{ChannelID: 50},
			Media: &tg.MessageMediaPoll{
				Poll:    tg.Poll{ID: 1, Question: tg.TextWithEntities{Text: "Q?"}, Answers: []tg.PollAnswerClass{&tg.PollAnswer{Text: tg.TextWithEntities{Text: "a"}, Option: []byte{0}}}},
				Results: tg.PollResults{},
			},
		}},
	})
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 7, PeerID: &tg.PeerChannel{ChannelID: 50}}))

	b := newMockBot(inv)

	poll, err := b.StopPoll(context.Background(), tdlibChannel(50), 7, nil)
	if err != nil {
		t.Fatalf("StopPoll: %v", err)
	}

	if !poll.IsClosed {
		t.Fatal("poll should be closed")
	}
}
