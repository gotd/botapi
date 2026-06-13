package botapi

import (
	"encoding/json"
	"testing"
)

// Compile-time proof that every concrete variant satisfies its sealed union.
var (
	_ ReplyMarkup = (*InlineKeyboardMarkup)(nil)
	_ ReplyMarkup = (*ReplyKeyboardMarkup)(nil)
	_ ReplyMarkup = (*ReplyKeyboardRemove)(nil)
	_ ReplyMarkup = (*ForceReply)(nil)

	_ ReactionType = ReactionTypeEmoji{}
	_ ReactionType = ReactionTypeCustomEmoji{}
	_ ReactionType = ReactionTypePaid{}

	_ MenuButton = MenuButtonCommands{}
	_ MenuButton = MenuButtonWebApp{}
	_ MenuButton = MenuButtonDefault{}

	_ MessageOrigin = (*MessageOriginUser)(nil)
	_ MessageOrigin = (*MessageOriginHiddenUser)(nil)
	_ MessageOrigin = (*MessageOriginChat)(nil)
	_ MessageOrigin = (*MessageOriginChannel)(nil)

	_ ChatMember = (*ChatMemberOwner)(nil)
	_ ChatMember = (*ChatMemberAdministrator)(nil)
	_ ChatMember = (*ChatMemberMember)(nil)
	_ ChatMember = (*ChatMemberRestricted)(nil)
	_ ChatMember = (*ChatMemberLeft)(nil)
	_ ChatMember = (*ChatMemberBanned)(nil)

	_ InputMedia = (*InputMediaPhoto)(nil)
	_ InputMedia = (*InputMediaVideo)(nil)
	_ InputMedia = (*InputMediaAnimation)(nil)
	_ InputMedia = (*InputMediaAudio)(nil)
	_ InputMedia = (*InputMediaDocument)(nil)
)

func TestUpdateJSONRoundTrip(t *testing.T) {
	in := Update{
		UpdateID: 42,
		Message: &Message{
			MessageID: 7,
			Date:      1700000000,
			From:      &User{ID: 1, FirstName: "Ada", Username: "ada"},
			Chat:      Chat{ID: -100123, Type: ChatTypeSupergroup, Title: "g"},
			Text:      "hi",
			Entities:  []MessageEntity{{Type: EntityBold, Offset: 0, Length: 2}},
		},
	}

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	var out Update
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if out.UpdateID != 42 || out.Message == nil {
		t.Fatalf("round trip lost data: %+v", out)
	}
	if out.Message.Chat.Type != ChatTypeSupergroup || out.Message.Text != "hi" {
		t.Fatalf("unexpected message: %+v", out.Message)
	}
	if len(out.Message.Entities) != 1 || out.Message.Entities[0].Type != EntityBold {
		t.Fatalf("unexpected entities: %+v", out.Message.Entities)
	}
}
