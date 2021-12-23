package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func TestBotAPI_convertToBotAPIEntities(t *testing.T) {
	tests := []struct {
		name  string
		input tg.MessageEntityClass
		wantR oas.MessageEntity
	}{
		{
			"Mention",
			&tg.MessageEntityMention{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeMention, Offset: 1, Length: 10},
		},
		{
			"Hashtag",
			&tg.MessageEntityHashtag{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeHashtag, Offset: 1, Length: 10},
		},
		{
			"BotCommand",
			&tg.MessageEntityBotCommand{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeBotCommand, Offset: 1, Length: 10},
		},
		{
			"URL",
			&tg.MessageEntityURL{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeURL, Offset: 1, Length: 10},
		},
		{
			"Email",
			&tg.MessageEntityEmail{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeEmail, Offset: 1, Length: 10},
		},
		{
			"Bold",
			&tg.MessageEntityBold{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeBold, Offset: 1, Length: 10},
		},
		{
			"Italic",
			&tg.MessageEntityItalic{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeItalic, Offset: 1, Length: 10},
		},
		{
			"Code",
			&tg.MessageEntityCode{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeCode, Offset: 1, Length: 10},
		},
		{
			"Pre",
			&tg.MessageEntityPre{Offset: 1, Length: 10, Language: "python"},
			oas.MessageEntity{Type: oas.MessageEntityTypePre, Offset: 1, Length: 10,
				Language: oas.NewOptString("python")},
		},
		{
			"TextURL",
			&tg.MessageEntityTextURL{Offset: 1, Length: 10, URL: "https://ya.ru"},
			oas.MessageEntity{Type: oas.MessageEntityTypeTextLink, Offset: 1, Length: 10,
				URL: oas.NewOptString("https://ya.ru")},
		},
		{
			"MentionName",
			&tg.MessageEntityMentionName{Offset: 1, Length: 10, UserID: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeTextMention, Offset: 1, Length: 10,
				User: oas.NewOptUser(convertRawToBotAPIUser(testUser()))},
		},
		{
			"Phone",
			&tg.MessageEntityPhone{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypePhoneNumber, Offset: 1, Length: 10},
		},
		{
			"Cashtag",
			&tg.MessageEntityCashtag{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeCashtag, Offset: 1, Length: 10},
		},
		{
			"Underline",
			&tg.MessageEntityUnderline{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeUnderline, Offset: 1, Length: 10},
		},
		{
			"Strikethrough",
			&tg.MessageEntityStrike{Offset: 1, Length: 10},
			oas.MessageEntity{Type: oas.MessageEntityTypeStrikethrough, Offset: 1, Length: 10},
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
				a.Equal(
					[]oas.MessageEntity{tt.wantR},
					api.convertToBotAPIEntities(ctx, []tg.MessageEntityClass{tt.input}),
				)
			})
		})
	}
}
