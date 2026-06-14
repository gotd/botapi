package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
)

func convMsg(t *testing.T, m *tg.Message) *Message {
	t.Helper()

	b := newMockBot(newMockInvoker())

	r, err := b.convertMessage(context.Background(), m)
	if err != nil {
		t.Fatalf("convertMessage: %v", err)
	}

	return r
}

func TestConvertMessageExposesRaw(t *testing.T) {
	m := &tg.Message{ID: 42, Message: "hi"}

	m.PeerID = &tg.PeerUser{UserID: 10}

	r := convMsg(t, m)
	if r.Raw() != m {
		t.Fatalf("Raw() = %p, want original %p", r.Raw(), m)
	}

	// Synthesized stubs (e.g. JSON-decoded messages) carry no raw message.
	if (&Message{}).Raw() != nil {
		t.Fatal("Raw() of a zero Message should be nil")
	}
}

func TestConvertTextMessageFromUser(t *testing.T) {
	m := &tg.Message{ID: 7, Message: "hello world"}

	m.PeerID = &tg.PeerUser{UserID: 10}
	m.SetFromID(&tg.PeerUser{UserID: 10})
	m.SetEntities([]tg.MessageEntityClass{&tg.MessageEntityBold{Offset: 0, Length: 5}})
	m.SetReplyTo(&tg.MessageReplyHeader{ReplyToMsgID: 3})
	m.SetReplyMarkup(&tg.ReplyInlineMarkup{Rows: []tg.KeyboardButtonRow{{
		Buttons: []tg.KeyboardButtonClass{&tg.KeyboardButtonCallback{Text: "ok", Data: []byte("d")}},
	}}})

	r := convMsg(t, m)
	if r.Text != "hello world" || r.From == nil || r.From.ID != 10 {
		t.Fatalf("msg from = %#v", r.From)
	}

	if len(r.Entities) != 1 || r.ReplyToMessage == nil || r.ReplyToMessage.MessageID != 3 {
		t.Fatalf("entities/reply = %#v", r)
	}

	if r.ReplyMarkup == nil || len(r.ReplyMarkup.InlineKeyboard) != 1 {
		t.Fatalf("markup = %#v", r.ReplyMarkup)
	}
}

func TestConvertPhotoWithCaption(t *testing.T) {
	m := &tg.Message{ID: 1, Message: "a caption"}

	m.PeerID = &tg.PeerUser{UserID: 10}
	m.SetFromID(&tg.PeerUser{UserID: 10})
	m.SetEntities([]tg.MessageEntityClass{&tg.MessageEntityBold{Offset: 0, Length: 1}})
	m.SetMedia(&tg.MessageMediaPhoto{Photo: &tg.Photo{
		ID: 1, AccessHash: 2, FileReference: []byte{1}, DCID: 2,
		Sizes: []tg.PhotoSizeClass{&tg.PhotoSize{Type: "x", W: 800, H: 600, Size: 1000}},
	}})

	r := convMsg(t, m)
	if len(r.Photo) != 1 || r.Caption != "a caption" || len(r.CaptionEntities) != 1 {
		t.Fatalf("photo msg = %#v", r)
	}
}

func docMessage(attrs ...tg.DocumentAttributeClass) *tg.Message {
	m := &tg.Message{ID: 1}

	m.PeerID = &tg.PeerUser{UserID: 10}
	m.SetFromID(&tg.PeerUser{UserID: 10})
	m.SetMedia(&tg.MessageMediaDocument{Document: &tg.Document{
		ID: 1, AccessHash: 2, FileReference: []byte{1}, DCID: 2, MimeType: "x", Size: 99, Attributes: attrs,
	}})

	return m
}

func TestConvertDocumentVariants(t *testing.T) {
	if r := convMsg(t, docMessage(&tg.DocumentAttributeSticker{Alt: "😀"})); r.Sticker == nil || r.Sticker.Emoji != "😀" {
		t.Fatalf("sticker = %#v", r.Sticker)
	}

	if r := convMsg(t, docMessage(&tg.DocumentAttributeAudio{Voice: true, Duration: 5})); r.Voice == nil {
		t.Fatal("voice not set")
	}

	if r := convMsg(t, docMessage(&tg.DocumentAttributeAudio{Performer: "p", Title: "t"})); r.Audio == nil {
		t.Fatal("audio not set")
	}

	if r := convMsg(t, docMessage(&tg.DocumentAttributeAnimated{}, &tg.DocumentAttributeVideo{W: 4, H: 4})); r.Animation == nil {
		t.Fatal("animation not set")
	}

	if r := convMsg(t, docMessage(&tg.DocumentAttributeVideo{RoundMessage: true, W: 240, H: 240})); r.VideoNote == nil {
		t.Fatal("video note not set")
	}

	if r := convMsg(t, docMessage(&tg.DocumentAttributeVideo{W: 1280, H: 720, Duration: 30})); r.Video == nil {
		t.Fatal("video not set")
	}

	if r := convMsg(t, docMessage(&tg.DocumentAttributeFilename{FileName: "f.bin"})); r.Document == nil {
		t.Fatal("document not set")
	}
}

func TestConvertContactAndPoll(t *testing.T) {
	mc := &tg.Message{ID: 1}

	mc.PeerID = &tg.PeerUser{UserID: 10}
	mc.SetFromID(&tg.PeerUser{UserID: 10})
	mc.SetMedia(&tg.MessageMediaContact{PhoneNumber: "+1", FirstName: "Ada", UserID: 5})

	if rc := convMsg(t, mc); rc.Contact == nil || rc.Contact.FirstName != "Ada" {
		t.Fatalf("contact = %#v", rc.Contact)
	}

	mp := &tg.Message{ID: 1}

	mp.PeerID = &tg.PeerUser{UserID: 10}
	mp.SetFromID(&tg.PeerUser{UserID: 10})
	mp.SetMedia(&tg.MessageMediaPoll{
		Poll: tg.Poll{ID: 1, Question: tg.TextWithEntities{Text: "Q?"}, Answers: []tg.PollAnswerClass{
			&tg.PollAnswer{Text: tg.TextWithEntities{Text: "a"}, Option: []byte{0}},
		}},
		Results: tg.PollResults{},
	})

	if rp := convMsg(t, mp); rp.Poll == nil || rp.Poll.Question != "Q?" {
		t.Fatalf("poll = %#v", rp.Poll)
	}
}

func fwdMessage(t *testing.T, fwd tg.MessageFwdHeader) *Message {
	t.Helper()

	m := &tg.Message{ID: 1, Message: "fwd"}

	m.PeerID = &tg.PeerUser{UserID: 10}
	m.SetFromID(&tg.PeerUser{UserID: 10})
	m.SetFwdFrom(fwd)

	return convMsg(t, m)
}

func TestConvertForwardOrigins(t *testing.T) {
	hiddenFwd := tg.MessageFwdHeader{Date: 1}
	hiddenFwd.SetFromName("Hidden")

	if r := fwdMessage(t, hiddenFwd); r.ForwardOrigin == nil {
		t.Fatal("hidden origin nil")
	} else if _, ok := r.ForwardOrigin.(*MessageOriginHiddenUser); !ok {
		t.Fatalf("origin = %T", r.ForwardOrigin)
	}

	userFwd := tg.MessageFwdHeader{Date: 1}
	userFwd.SetFromID(&tg.PeerUser{UserID: 20})

	if r := fwdMessage(t, userFwd); func() bool { _, ok := r.ForwardOrigin.(*MessageOriginUser); return !ok }() {
		t.Fatalf("expected user origin, got %T", r.ForwardOrigin)
	}

	chFwd := tg.MessageFwdHeader{Date: 1}
	chFwd.SetFromID(&tg.PeerChannel{ChannelID: 50})
	chFwd.SetChannelPost(99)
	chFwd.SetPostAuthor("Editor")

	// Use a broadcast channel so the origin is classified as a channel post.
	inv := newMockInvoker()
	inv.handle(tg.ChannelsGetChannelsRequestTypeID, func(buf *bin.Buffer) (bin.Encoder, error) {
		return &tg.MessagesChats{Chats: []tg.ChatClass{
			&tg.Channel{ID: 50, AccessHash: 1, Title: "ch", Broadcast: true, Photo: &tg.ChatPhotoEmpty{}},
		}}, nil
	})

	b := newMockBot(inv)
	m := &tg.Message{ID: 1, Message: "fwd"}

	m.PeerID = &tg.PeerUser{UserID: 10}
	m.SetFromID(&tg.PeerUser{UserID: 10})
	m.SetFwdFrom(chFwd)

	r, err := b.convertMessage(context.Background(), m)
	if err != nil {
		t.Fatalf("convertMessage: %v", err)
	}

	origin, ok := r.ForwardOrigin.(*MessageOriginChannel)
	if !ok {
		t.Fatalf("expected channel origin, got %T", r.ForwardOrigin)
	}

	if origin.MessageID != 99 || origin.AuthorSignature != "Editor" {
		t.Fatalf("channel origin = %#v", origin)
	}
}

func TestConvertSenderChat(t *testing.T) {
	m := &tg.Message{ID: 1, Message: "post"}

	m.PeerID = &tg.PeerChannel{ChannelID: 50}
	m.SetFromID(&tg.PeerChannel{ChannelID: 50})

	if r := convMsg(t, m); r.SenderChat == nil {
		t.Fatalf("sender chat not set: %#v", r)
	}
}

func TestConvertOutgoingFillsFromSelf(t *testing.T) {
	m := &tg.Message{ID: 1, Message: "mine", Out: true}

	m.PeerID = &tg.PeerUser{UserID: 10}

	if r := convMsg(t, m); r.From == nil || !r.From.IsBot {
		t.Fatalf("outgoing From = %#v", r.From)
	}
}
