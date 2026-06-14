package botapi

import (
	"bytes"
	"context"
	"testing"

	"github.com/gotd/td/telegram/message/rich"
	"github.com/gotd/td/tg"
)

func TestSendRichMessage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMessageRequestTypeID, sendMediaOK())
	b := newMockBot(inv)

	msg := rich.New(rich.Paragraph(rich.Plain("hi"))).Input()
	if _, err := b.SendRichMessage(context.Background(), userRef(10, 20), msg); err != nil {
		t.Fatalf("SendRichMessage: %v", err)
	}
	if !inv.called(tg.MessagesSendMessageRequestTypeID) {
		t.Fatal("expected messages.sendMessage")
	}
}

func TestSendRichHTMLAndMarkdown(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMessageRequestTypeID, sendMediaOK())
	b := newMockBot(inv)

	if _, err := b.SendRichHTML(context.Background(), userRef(10, 20), "<p>hi</p>"); err != nil {
		t.Fatalf("SendRichHTML: %v", err)
	}
	if _, err := b.SendRichMarkdown(context.Background(), userRef(10, 20), "**hi**"); err != nil {
		t.Fatalf("SendRichMarkdown: %v", err)
	}
}

func TestSendInvoice(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, sendMediaOK())
	b := newMockBot(inv)

	params := InvoiceParams{
		Title:       "Item",
		Description: "An item",
		Payload:     "payload",
		Currency:    "USD",
		Prices:      []LabeledPrice{{Label: "Item", Amount: 1000}},
	}
	if _, err := b.SendInvoice(context.Background(), userRef(10, 20), params); err != nil {
		t.Fatalf("SendInvoice: %v", err)
	}
	var req tg.MessagesSendMediaRequest
	inv.decode(t, tg.MessagesSendMediaRequestTypeID, &req)
	if _, ok := req.Media.(*tg.InputMediaInvoice); !ok {
		t.Fatalf("media = %#v", req.Media)
	}
}

func TestSendMediaGroup(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUploadMediaRequestTypeID, &tg.MessageMediaPhoto{
		Photo: &tg.Photo{ID: 1, AccessHash: 2, FileReference: []byte{1}, DCID: 2, Sizes: []tg.PhotoSizeClass{&tg.PhotoSize{Type: "x", W: 1, H: 1, Size: 1}}},
	})
	inv.reply(tg.MessagesSendMultiMediaRequestTypeID, &tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateNewMessage{Message: &tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}}},
			&tg.UpdateNewMessage{Message: &tg.Message{ID: 2, PeerID: &tg.PeerUser{UserID: 10}}},
		},
		Users: []tg.UserClass{&tg.User{ID: 10, AccessHash: 20}},
	})
	b := newMockBot(inv)

	media := []InputMedia{
		&InputMediaPhoto{Media: FileFromBytes("a.jpg", []byte("a"))},
		&InputMediaPhoto{Media: FileFromBytes("b.jpg", []byte("b"))},
	}
	msgs, err := b.SendMediaGroup(context.Background(), userRef(10, 20), media)
	if err != nil {
		t.Fatalf("SendMediaGroup: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("messages = %d", len(msgs))
	}
}

func TestPeerRefResolves(t *testing.T) {
	b := newMockBot(newMockInvoker())
	ref, err := b.PeerRef(context.Background(), tdlibChannel(50))
	if err != nil {
		t.Fatalf("PeerRef: %v", err)
	}
	if ref.Kind != peerKindChannel || ref.ID != 50 {
		t.Fatalf("ref = %#v", ref)
	}
}

func TestSetChatPhotoUpload(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsEditPhotoRequestTypeID, okUpdates())
	b := newMockBot(inv)

	err := b.SetChatPhoto(context.Background(), tdlibChannel(50), FileFromBytes("p.jpg", []byte("img")))
	if err != nil {
		t.Fatalf("SetChatPhoto: %v", err)
	}
	if !inv.called(tg.ChannelsEditPhotoRequestTypeID) {
		t.Fatal("expected channels.editPhoto")
	}
}

func TestSetChatPhotoRejectsNonUpload(t *testing.T) {
	b := newMockBot(newMockInvoker())
	if err := b.SetChatPhoto(context.Background(), tdlibChannel(50), FileID("x")); err == nil {
		t.Fatal("SetChatPhoto should reject a non-uploaded file")
	}
}

func TestDownloadFileErrors(t *testing.T) {
	b := newMockBot(newMockInvoker())
	if _, err := b.DownloadFile(context.Background(), "not-a-file-id", nil); err == nil {
		t.Fatal("expected error for invalid file_id")
	}
	if err := b.DownloadFileToPath(context.Background(), "not-a-file-id", "/tmp/x"); err == nil {
		t.Fatal("expected error for invalid file_id")
	}
}

func TestCountWriter(t *testing.T) {
	var buf bytes.Buffer
	cw := &countWriter{w: &buf}
	n, err := cw.Write([]byte("hello"))
	if err != nil || n != 5 || cw.n != 5 {
		t.Fatalf("write: n=%d cw.n=%d err=%v", n, cw.n, err)
	}
	_, _ = cw.Write([]byte("!"))
	if cw.n != 6 || buf.String() != "hello!" {
		t.Fatalf("count = %d buf = %q", cw.n, buf.String())
	}
}

// TestOptionBuildersThroughMethods covers the remaining option constructors by
// passing them through their methods.
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
