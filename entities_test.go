package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestEntitiesRoundTrip(t *testing.T) {
	in := []MessageEntity{
		{Type: EntityBold, Offset: 0, Length: 4},
		{Type: EntityTextLink, Offset: 5, Length: 3, URL: "https://example.com"},
		{Type: EntityPre, Offset: 9, Length: 2, Language: "go"},
		{Type: EntityExpandableBlockquote, Offset: 12, Length: 1},
		{Type: EntityTextMention, Offset: 14, Length: 2, User: &User{ID: 777}},
		{Type: EntityCustomEmoji, Offset: 17, Length: 1, CustomEmojiID: "12345"},
	}

	tgEnts := entitiesToTg(in)
	if len(tgEnts) != len(in) {
		t.Fatalf("to-tg length: got %d want %d", len(tgEnts), len(in))
	}
	if bq, ok := tgEnts[3].(*tg.MessageEntityBlockquote); !ok || !bq.Collapsed {
		t.Fatalf("expandable blockquote should set Collapsed: %#v", tgEnts[3])
	}
	if mn, ok := tgEnts[4].(*tg.MessageEntityMentionName); !ok || mn.UserID != 777 {
		t.Fatalf("text mention user id lost: %#v", tgEnts[4])
	}
	if ce, ok := tgEnts[5].(*tg.MessageEntityCustomEmoji); !ok || ce.DocumentID != 12345 {
		t.Fatalf("custom emoji doc id lost: %#v", tgEnts[5])
	}

	out := entitiesFromTg(tgEnts)
	if len(out) != len(in) {
		t.Fatalf("from-tg length: got %d want %d", len(out), len(in))
	}
	for i := range in {
		if out[i].Type != in[i].Type || out[i].Offset != in[i].Offset || out[i].Length != in[i].Length {
			t.Fatalf("entity %d mismatch: got %+v want %+v", i, out[i], in[i])
		}
	}
	if out[1].URL != "https://example.com" || out[2].Language != "go" {
		t.Fatalf("attribute fields lost: %+v %+v", out[1], out[2])
	}
	if out[4].User == nil || out[4].User.ID != 777 {
		t.Fatalf("text mention user lost: %+v", out[4])
	}
}

func TestEntitiesEmpty(t *testing.T) {
	if entitiesToTg(nil) != nil || entitiesFromTg(nil) != nil {
		t.Fatal("empty input should yield nil")
	}
}

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
}
