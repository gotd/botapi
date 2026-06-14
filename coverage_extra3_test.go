package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestEntitiesAllTypes(t *testing.T) {
	entities := []MessageEntity{
		{Type: EntityMention, Offset: 0, Length: 1},
		{Type: EntityHashtag, Offset: 1, Length: 1},
		{Type: EntityCashtag, Offset: 2, Length: 1},
		{Type: EntityBotCommand, Offset: 3, Length: 1},
		{Type: EntityURL, Offset: 4, Length: 1},
		{Type: EntityEmail, Offset: 5, Length: 1},
		{Type: EntityPhoneNumber, Offset: 6, Length: 1},
		{Type: EntityBold, Offset: 7, Length: 1},
		{Type: EntityItalic, Offset: 8, Length: 1},
		{Type: EntityUnderline, Offset: 9, Length: 1},
		{Type: EntityStrikethrough, Offset: 10, Length: 1},
		{Type: EntitySpoiler, Offset: 11, Length: 1},
		{Type: EntityBlockquote, Offset: 12, Length: 1},
		{Type: EntityExpandableBlockquote, Offset: 13, Length: 1},
		{Type: EntityCode, Offset: 14, Length: 1},
		{Type: EntityPre, Offset: 15, Length: 1, Language: "go"},
		{Type: EntityTextLink, Offset: 16, Length: 1, URL: "https://e"},
		{Type: EntityTextMention, Offset: 17, Length: 1, User: &User{ID: 42}},
		{Type: EntityCustomEmoji, Offset: 18, Length: 1, CustomEmojiID: "555"},
	}
	tgEntities := entitiesToTg(entities)
	if len(tgEntities) != len(entities) {
		t.Fatalf("toTg produced %d, want %d", len(tgEntities), len(entities))
	}
	back := entitiesFromTg(tgEntities)
	if len(back) != len(entities) {
		t.Fatalf("fromTg produced %d, want %d", len(back), len(entities))
	}
	// Spot-check a few that carry extra data.
	byType := map[MessageEntityType]MessageEntity{}
	for _, e := range back {
		byType[e.Type] = e
	}
	if byType[EntityPre].Language != "go" {
		t.Fatalf("pre lang = %q", byType[EntityPre].Language)
	}
	if byType[EntityTextLink].URL != "https://e" {
		t.Fatalf("text link url = %q", byType[EntityTextLink].URL)
	}
	if byType[EntityTextMention].User == nil || byType[EntityTextMention].User.ID != 42 {
		t.Fatalf("text mention user = %#v", byType[EntityTextMention].User)
	}
	if byType[EntityCustomEmoji].CustomEmojiID != "555" {
		t.Fatalf("custom emoji id = %q", byType[EntityCustomEmoji].CustomEmojiID)
	}
	if _, ok := byType[EntityExpandableBlockquote]; !ok {
		t.Fatal("expandable blockquote lost")
	}

	// Empty input passes through.
	if entitiesToTg(nil) != nil || entitiesFromTg(nil) != nil {
		t.Fatal("nil entities should produce nil")
	}
}

func TestSendChatActionAll(t *testing.T) {
	actions := []ChatAction{
		ChatActionTyping, ChatActionUploadPhoto, ChatActionRecordVideo, ChatActionUploadVideo,
		ChatActionRecordVoice, ChatActionUploadVoice, ChatActionUploadDocument, ChatActionChooseSticker,
		ChatActionFindLocation, ChatActionRecordVideoNote, ChatActionUploadVideoNote,
	}
	for _, action := range actions {
		inv := newMockInvoker()
		inv.reply(tg.MessagesSetTypingRequestTypeID, &tg.BoolTrue{})
		b := newMockBot(inv)
		if err := b.SendChatAction(context.Background(), userRef(10, 20), action); err != nil {
			t.Fatalf("SendChatAction(%q): %v", action, err)
		}
	}
	// Unknown action is rejected before the wire.
	b := newMockBot(newMockInvoker())
	if err := b.SendChatAction(context.Background(), userRef(10, 20), ChatAction("nonsense")); err == nil {
		t.Fatal("unknown action should error")
	}
}

func TestEditMessageMediaVariants(t *testing.T) {
	fid := documentFileID(t, 0x90)
	cases := map[string]InputMedia{
		"video-url":      &InputMediaVideo{Media: FileURL("https://e/v.mp4"), Caption: "c"},
		"animation-id":   &InputMediaAnimation{Media: FileID(fid)},
		"audio-url":      &InputMediaAudio{Media: FileURL("https://e/a.mp3")},
		"document-id":    &InputMediaDocument{Media: FileID(fid)},
		"photo-id":       &InputMediaPhoto{Media: FileID(photoFileID(t, 0x91))},
		"document-bytes": &InputMediaDocument{Media: FileFromBytes("f.bin", []byte("data"))},
	}
	for name, media := range cases {
		t.Run(name, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))
			b := newMockBot(inv)
			if _, err := b.EditMessageMedia(context.Background(), userRef(10, 20), 5, media); err != nil {
				t.Fatalf("EditMessageMedia(%s): %v", name, err)
			}
		})
	}
}

func TestSendMediaGroupDocuments(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUploadMediaRequestTypeID, &tg.MessageMediaDocument{
		Document: &tg.Document{ID: 1, AccessHash: 2, FileReference: []byte{1}, DCID: 2, MimeType: "application/pdf", Size: 10},
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
		&InputMediaDocument{Media: FileFromBytes("a.pdf", []byte("a")), Caption: "first"},
		&InputMediaVideo{Media: FileFromBytes("b.mp4", []byte("b"))},
	}
	msgs, err := b.SendMediaGroup(context.Background(), userRef(10, 20), media)
	if err != nil {
		t.Fatalf("SendMediaGroup docs: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("messages = %d", len(msgs))
	}
}

func TestSendMediaGroupRejectsNonUpload(t *testing.T) {
	b := newMockBot(newMockInvoker())
	// Media groups only support uploads, not URLs/file_ids.
	_, err := b.SendMediaGroup(context.Background(), userRef(10, 20), []InputMedia{
		&InputMediaPhoto{Media: FileURL("https://e/p.jpg")},
	})
	if err == nil {
		t.Fatal("media group with URL media should error")
	}
}
