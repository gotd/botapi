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
