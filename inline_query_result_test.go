package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

func TestInlineResultArticle(t *testing.T) {
	b := &Bot{}
	r := &InlineQueryResultArticle{
		ID:    "1",
		Title: "Hello",
		InputMessageContent: &InputTextMessageContent{
			MessageText: "world",
		},
		URL:          "https://example.com",
		ThumbnailURL: "https://example.com/t.jpg",
		Description:  "desc",
	}
	got, err := r.toTg(context.Background(), b)
	if err != nil {
		t.Fatal(err)
	}
	art, ok := got.(*tg.InputBotInlineResult)
	if !ok {
		t.Fatalf("got %T", got)
	}
	if art.ID != "1" || art.Type != "article" || art.Title != "Hello" {
		t.Fatalf("bad article: %#v", art)
	}
	if url, hasURL := art.GetURL(); !hasURL || url != "https://example.com" {
		t.Fatalf("url: %q %v", url, ok)
	}
	if _, hasThumb := art.GetThumb(); !hasThumb {
		t.Fatal("thumb not set")
	}
	text, ok := art.SendMessage.(*tg.InputBotInlineMessageText)
	if !ok || text.Message != "world" {
		t.Fatalf("send message: %#v", art.SendMessage)
	}
}

func TestInlineResultArticleRequiresContent(t *testing.T) {
	b := &Bot{}
	r := &InlineQueryResultArticle{ID: "1", Title: "x"}
	if _, err := r.toTg(context.Background(), b); err == nil {
		t.Fatal("expected error when input_message_content is missing")
	}
}

func TestInlineResultCachedSticker(t *testing.T) {
	b := &Bot{}
	// A document file_id round-trips through inputDocumentFromFileID.
	fid := documentFileID(t, 0x42)
	r := &InlineQueryResultCachedSticker{ID: "s", StickerFileID: fid}
	got, err := r.toTg(context.Background(), b)
	if err != nil {
		t.Fatal(err)
	}
	doc, ok := got.(*tg.InputBotInlineResultDocument)
	if !ok || doc.Type != "sticker" {
		t.Fatalf("got %#v", got)
	}
	if d, ok := doc.Document.(*tg.InputDocument); !ok || d.ID != 0x42 {
		t.Fatalf("document: %#v", doc.Document)
	}
	// No caption and no content: defaults to an empty media-auto message.
	if _, ok := doc.SendMessage.(*tg.InputBotInlineMessageMediaAuto); !ok {
		t.Fatalf("send message: %#v", doc.SendMessage)
	}
}

func TestInlineResultPhotoURLWithCaption(t *testing.T) {
	b := &Bot{}
	r := &InlineQueryResultPhoto{
		ID:           "p",
		PhotoURL:     "https://example.com/p.jpg",
		ThumbnailURL: "https://example.com/t.jpg",
		PhotoWidth:   100,
		PhotoHeight:  80,
		captioned:    captioned{Caption: "cap"},
	}
	got, err := r.toTg(context.Background(), b)
	if err != nil {
		t.Fatal(err)
	}
	res, ok := got.(*tg.InputBotInlineResult)
	if !ok || res.Type != "photo" {
		t.Fatalf("got %#v", got)
	}
	content, ok := res.GetContent()
	if !ok || content.URL != "https://example.com/p.jpg" {
		t.Fatalf("content: %#v", content)
	}
	auto, ok := res.SendMessage.(*tg.InputBotInlineMessageMediaAuto)
	if !ok || auto.Message != "cap" {
		t.Fatalf("send: %#v", res.SendMessage)
	}
}

func TestInputContactContent(t *testing.T) {
	c := &InputContactMessageContent{PhoneNumber: "+1", FirstName: "Joe"}
	got, err := c.toTg(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	contact, ok := got.(*tg.InputBotInlineMessageMediaContact)
	if !ok || contact.PhoneNumber != "+1" || contact.FirstName != "Joe" {
		t.Fatalf("contact: %#v", got)
	}
}

// documentFileID builds a valid document file_id with the given media id, for
// tests that exercise file_id-backed results.
func documentFileID(t *testing.T, id int64) string {
	t.Helper()
	s, err := fileid.EncodeFileID(fileid.FileID{
		Type:       fileid.Document,
		DC:         2,
		ID:         id,
		AccessHash: 7,
	})
	if err != nil {
		t.Fatalf("encode file id: %v", err)
	}
	return s
}
