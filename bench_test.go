package botapi

import (
	"testing"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

// Benchmarks for the pure translation hot paths exercised on every incoming and
// outgoing message. Run with: go test -bench . -benchmem

func BenchmarkEntitiesToTg(b *testing.B) {
	entities := []MessageEntity{
		{Type: EntityBold, Offset: 0, Length: 4},
		{Type: EntityItalic, Offset: 5, Length: 6},
		{Type: EntityCode, Offset: 12, Length: 3},
		{Type: EntityTextLink, Offset: 16, Length: 8, URL: "https://example.com"},
		{Type: EntityTextMention, Offset: 25, Length: 5, User: &User{ID: 42}},
	}

	b.ReportAllocs()

	for b.Loop() {
		_ = entitiesToTg(entities)
	}
}

func BenchmarkEntitiesFromTg(b *testing.B) {
	entities := []tg.MessageEntityClass{
		&tg.MessageEntityBold{Offset: 0, Length: 4},
		&tg.MessageEntityItalic{Offset: 5, Length: 6},
		&tg.MessageEntityCode{Offset: 12, Length: 3},
		&tg.MessageEntityTextURL{Offset: 16, Length: 8, URL: "https://example.com"},
		&tg.MessageEntityMentionName{Offset: 25, Length: 5, UserID: 42},
	}

	b.ReportAllocs()

	for b.Loop() {
		_ = entitiesFromTg(entities)
	}
}

func BenchmarkReplyMarkupToTg(b *testing.B) {
	q := "go"
	markup := &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{
		{
			{Text: "site", URL: "https://example.com"},
			{Text: "ok", CallbackData: "ok:1"},
		},
		{
			{Text: "switch", SwitchInlineQuery: &q},
		},
	}}

	b.ReportAllocs()

	for b.Loop() {
		if _, err := replyMarkupToTg(markup); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserFromTgUser(b *testing.B) {
	u := &tg.User{
		ID:        42,
		FirstName: "Ada",
		LastName:  "Lovelace",
		Username:  "ada",
		Bot:       false,
	}

	b.ReportAllocs()

	for b.Loop() {
		_ = userFromTgUser(u)
	}
}

func BenchmarkFileUniqueID(b *testing.B) {
	f := fileid.FileID{Type: fileid.Document, ID: 0x1122334455667788}

	b.ReportAllocs()

	for b.Loop() {
		_ = fileUniqueID(f)
	}
}
