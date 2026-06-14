package botapi

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-faster/jx"
)

// strptr is a helper for optional string fields in test fixtures.
func strptr(s string) *string { return &s }

// fullMessage returns a Message with every field populated, including the
// polymorphic forward_origin, nested messages, all media variants and a poll
// with explanation entities. It exercises every branch of Message.Encode and
// Message.Decode.
func fullMessage() *Message {
	user := &User{
		ID:                      1,
		IsBot:                   true,
		FirstName:               "Ada",
		LastName:                "Lovelace",
		Username:                "ada",
		LanguageCode:            "en",
		IsPremium:               true,
		AddedToAttachmentMenu:   true,
		CanJoinGroups:           true,
		CanReadAllGroupMessages: true,
		SupportsInlineQueries:   true,
	}
	chat := Chat{
		ID:        -100,
		Type:      ChatTypeSupergroup,
		Title:     "Math",
		Username:  "math",
		FirstName: "f",
		LastName:  "l",
		IsForum:   true,
	}
	thumb := &PhotoSize{FileID: "t", FileUniqueID: "tu", Width: 1, Height: 2, FileSize: 3}
	return &Message{
		MessageID:           7,
		MessageThreadID:     3,
		From:                user,
		SenderChat:          &chat,
		Date:                1700000000,
		Chat:                chat,
		ForwardOrigin:       &MessageOriginChannel{Type: OriginChannel, Date: 5, Chat: chat, MessageID: 9, AuthorSignature: "sig"},
		ReplyToMessage:      &Message{MessageID: 1, Date: 1, Chat: chat, Text: "parent"},
		ViaBot:              user,
		EditDate:            1700000100,
		HasProtectedContent: true,
		MediaGroupID:        "mg",
		AuthorSignature:     "auth",
		Text:                "hello",
		Entities:            []MessageEntity{{Type: EntityBold, Offset: 0, Length: 5}, {Type: EntityTextMention, Offset: 1, Length: 2, User: user, URL: "u", Language: "go", CustomEmojiID: "e"}},
		Caption:             "cap",
		CaptionEntities:     []MessageEntity{{Type: EntityItalic, Offset: 0, Length: 3}},
		Animation:           &Animation{FileID: "a", FileUniqueID: "au", Width: 1, Height: 2, Duration: 3, Thumbnail: thumb, FileName: "a.gif", MIMEType: "image/gif", FileSize: 4},
		Audio:               &Audio{FileID: "b", FileUniqueID: "bu", Duration: 10, Performer: "p", Title: "t", FileName: "b.mp3", MIMEType: "audio/mpeg", FileSize: 5, Thumbnail: thumb},
		Document:            &Document{FileID: "c", FileUniqueID: "cu", Thumbnail: thumb, FileName: "c.pdf", MIMEType: "application/pdf", FileSize: 6},
		Photo:               []PhotoSize{{FileID: "p1", FileUniqueID: "p1u", Width: 100, Height: 200, FileSize: 7}},
		Sticker:             &Sticker{FileID: "s", FileUniqueID: "su", Type: StickerRegular, Width: 512, Height: 512, IsAnimated: true, IsVideo: true, Thumbnail: thumb, Emoji: "🙂", SetName: "set", FileSize: 8},
		Video:               &Video{FileID: "v", FileUniqueID: "vu", Width: 640, Height: 480, Duration: 12, Thumbnail: thumb, FileName: "v.mp4", MIMEType: "video/mp4", FileSize: 9},
		VideoNote:           &VideoNote{FileID: "vn", FileUniqueID: "vnu", Length: 240, Duration: 6, Thumbnail: thumb, FileSize: 11},
		Voice:               &Voice{FileID: "vo", FileUniqueID: "vou", Duration: 4, MIMEType: "audio/ogg", FileSize: 12},
		Contact:             &Contact{PhoneNumber: "+1", FirstName: "Ada", LastName: "L", UserID: 1, VCard: "vc"},
		Dice:                &Dice{Emoji: DiceDart, Value: 6},
		Poll: &Poll{
			ID:                    "poll1",
			Question:              "q?",
			Options:               []PollOption{{Text: "a", VoterCount: 1}, {Text: "b", VoterCount: 2}},
			TotalVoterCount:       3,
			IsClosed:              true,
			IsAnonymous:           true,
			Type:                  PollQuiz,
			AllowsMultipleAnswers: true,
			CorrectOptionID:       1,
			Explanation:           "because",
			ExplanationEntities:   []MessageEntity{{Type: EntityCode, Offset: 0, Length: 1}},
			OpenPeriod:            60,
		},
		Venue:          &Venue{Location: Location{Longitude: 1.5, Latitude: 2.25}, Title: "v", Address: "addr", FoursquareID: "fid", FoursquareType: "ft", GooglePlaceID: "gid", GooglePlaceType: "gt"},
		Location:       &Location{Longitude: 3.5, Latitude: 4.75, HorizontalAccuracy: 1.5, LivePeriod: 60, Heading: 90, ProximityAlertRadius: 100},
		NewChatMembers: []User{{ID: 2, FirstName: "Bob"}},
		LeftChatMember: &User{ID: 3, FirstName: "Carl"},
		NewChatTitle:   "New",
		PinnedMessage:  &Message{MessageID: 2, Date: 2, Chat: chat, Text: "pinned"},
		ReplyMarkup: &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{
			{Text: "url", URL: "https://x"},
			{Text: "cb", CallbackData: "d"},
			{Text: "wa", WebApp: &WebAppInfo{URL: "https://app"}},
			{Text: "si", SwitchInlineQuery: strptr("q"), SwitchInlineQueryCurrentChat: strptr("c"), Pay: true},
		}}},
	}
}

func TestMessageJXRoundTrip(t *testing.T) {
	in := fullMessage()

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out Message
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, &out) {
		t.Fatalf("round trip mismatch:\n in: %+v\nout: %+v", in, &out)
	}

	// Re-marshaling the decoded value must reproduce identical bytes.
	again, err := json.Marshal(&out)
	if err != nil {
		t.Fatalf("re-marshal: %v", err)
	}
	if !bytes.Equal(again, data) {
		t.Fatalf("non-idempotent encoding:\n first: %s\nsecond: %s", data, again)
	}
}

// jsonRoundTrip marshals in through encoding/json (exercising its MarshalJSON),
// parses it back (exercising UnmarshalJSON) and asserts the value is preserved.
func jsonRoundTrip[T any](t *testing.T, in T) {
	t.Helper()
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal %T: %v", in, err)
	}
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal %T: %v", in, err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("%T round trip mismatch:\n in: %+v\nout: %+v", in, in, out)
	}
}

// TestLeafEntitiesJSON round-trips every receivable entity directly through
// encoding/json, covering each type's MarshalJSON/UnmarshalJSON wrappers.
func TestLeafEntitiesJSON(t *testing.T) {
	thumb := &PhotoSize{FileID: "t", FileUniqueID: "tu", Width: 1, Height: 2, FileSize: 3}
	jsonRoundTrip(t, User{ID: 1, FirstName: "A", LastName: "B", Username: "u", LanguageCode: "en", IsBot: true, IsPremium: true, AddedToAttachmentMenu: true, CanJoinGroups: true, CanReadAllGroupMessages: true, SupportsInlineQueries: true})
	jsonRoundTrip(t, Chat{ID: -1, Type: ChatTypeGroup, Title: "t", Username: "u", FirstName: "f", LastName: "l", IsForum: true})
	jsonRoundTrip(t, PhotoSize{FileID: "p", FileUniqueID: "pu", Width: 1, Height: 2, FileSize: 3})
	jsonRoundTrip(t, Animation{FileID: "a", FileUniqueID: "au", Width: 1, Height: 2, Duration: 3, Thumbnail: thumb, FileName: "a", MIMEType: "m", FileSize: 4})
	jsonRoundTrip(t, Audio{FileID: "b", FileUniqueID: "bu", Duration: 5, Performer: "p", Title: "t", FileName: "f", MIMEType: "m", FileSize: 6, Thumbnail: thumb})
	jsonRoundTrip(t, Document{FileID: "c", FileUniqueID: "cu", Thumbnail: thumb, FileName: "f", MIMEType: "m", FileSize: 7})
	jsonRoundTrip(t, Video{FileID: "v", FileUniqueID: "vu", Width: 1, Height: 2, Duration: 3, Thumbnail: thumb, FileName: "f", MIMEType: "m", FileSize: 8})
	jsonRoundTrip(t, VideoNote{FileID: "vn", FileUniqueID: "vnu", Length: 1, Duration: 2, Thumbnail: thumb, FileSize: 9})
	jsonRoundTrip(t, Voice{FileID: "vo", FileUniqueID: "vou", Duration: 3, MIMEType: "m", FileSize: 10})
	jsonRoundTrip(t, Sticker{FileID: "s", FileUniqueID: "su", Type: StickerMask, Width: 1, Height: 2, IsAnimated: true, IsVideo: true, Thumbnail: thumb, Emoji: "x", SetName: "set", FileSize: 11})
	jsonRoundTrip(t, MessageEntity{Type: EntityTextLink, Offset: 1, Length: 2, URL: "u", User: &User{ID: 1, FirstName: "A"}, Language: "go", CustomEmojiID: "e"})
	jsonRoundTrip(t, Contact{PhoneNumber: "+1", FirstName: "A", LastName: "B", UserID: 1, VCard: "v"})
	jsonRoundTrip(t, Dice{Emoji: DiceBasketball, Value: 5})
	jsonRoundTrip(t, Location{Longitude: 1.5, Latitude: 2.25, HorizontalAccuracy: 0.5, LivePeriod: 60, Heading: 90, ProximityAlertRadius: 10})
	jsonRoundTrip(t, Venue{Location: Location{Longitude: 1.5, Latitude: 2.25}, Title: "t", Address: "a", FoursquareID: "f", FoursquareType: "ft", GooglePlaceID: "g", GooglePlaceType: "gt"})
	jsonRoundTrip(t, PollOption{Text: "a", VoterCount: 3})
	jsonRoundTrip(t, Poll{ID: "p", Question: "q", Options: []PollOption{{Text: "a", VoterCount: 1}}, TotalVoterCount: 1, Type: PollRegular, IsClosed: true, IsAnonymous: true, AllowsMultipleAnswers: true, CorrectOptionID: 0, Explanation: "e", ExplanationEntities: []MessageEntity{{Type: EntityCode, Length: 1}}, OpenPeriod: 30})
	jsonRoundTrip(t, WebAppInfo{URL: "https://x"})
	jsonRoundTrip(t, InlineKeyboardButton{Text: "t", URL: "u", CallbackData: "d", WebApp: &WebAppInfo{URL: "w"}, SwitchInlineQuery: strptr("q"), SwitchInlineQueryCurrentChat: strptr("c"), Pay: true})
	jsonRoundTrip(t, InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{{Text: "a"}}}})

	// Pointer-receiver entities (always used via pointer on the wire).
	jsonRoundTrip(t, &MessageOriginUser{Type: OriginUser, Date: 1, SenderUser: User{ID: 1, FirstName: "A"}})
	jsonRoundTrip(t, &MessageOriginHiddenUser{Type: OriginHiddenUser, Date: 2, SenderUserName: "g"})
	jsonRoundTrip(t, &MessageOriginChat{Type: OriginChat, Date: 3, SenderChat: Chat{ID: 1, Type: ChatTypeChannel}, AuthorSignature: "s"})
	jsonRoundTrip(t, &MessageOriginChannel{Type: OriginChannel, Date: 4, Chat: Chat{ID: 1, Type: ChatTypeChannel}, MessageID: 5, AuthorSignature: "s"})
	jsonRoundTrip(t, &Message{MessageID: 1, Date: 2, Chat: Chat{ID: 1, Type: ChatTypePrivate}, Text: "hi"})
}

// TestDecodeTruncated feeds every prefix of a fully-populated Message document
// to Decode. Each truncation point cuts the input inside a different field, so
// collectively the sweep exercises the per-field error-return paths across the
// whole transitive type set. The contract under test: a truncated document
// always yields an error and never a panic.
func TestDecodeTruncated(t *testing.T) {
	data, err := json.Marshal(fullMessage())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for i := range len(data) {
		var m Message
		// Decode directly (not via encoding/json, which rejects a truncated
		// outer object before delegating). Must error, never panic.
		if err := unmarshalJX(data[:i], &m); err == nil {
			t.Fatalf("prefix len %d decoded without error", i)
		}
	}
}

// TestDecodeTypeMismatch checks that a wrong JSON type for a field surfaces an
// error rather than being silently coerced, across a representative sample of
// scalar field kinds.
func TestDecodeTypeMismatch(t *testing.T) {
	cases := []string{
		`{"message_id":"x"}`,                     // int field, string value
		`{"date":true}`,                          // int field, bool value
		`{"has_protected_content":1}`,            // bool field, number value
		`{"text":123}`,                           // string field, number value
		`{"from":[]}`,                            // object field, array value
		`{"chat":7}`,                             // nested object, number value
		`{"entities":{}}`,                        // array field, object value
		`{"photo":5}`,                            // array field, number value
		`{"reply_markup":{"inline_keyboard":7}}`, // nested array, number value
		`{"forward_origin":{"type":"user","date":"x"}}`, // union variant, bad field
	}
	for _, c := range cases {
		var m Message
		if err := json.Unmarshal([]byte(c), &m); err == nil {
			t.Fatalf("expected error for %s", c)
		}
	}
}

// fieldErrorProbe pairs a fully-populated value (to enumerate its wire field
// names) with a constructor for a fresh decode target.
type fieldErrorProbe struct {
	full   any
	newDec func() jsonDecoder
}

// TestDecodeFieldErrors feeds an invalid value token to each wire field of each
// entity and asserts the field's decoder reports an error. This deterministically
// exercises the per-field error-return branch of every Decode method.
func TestDecodeFieldErrors(t *testing.T) {
	thumb := &PhotoSize{FileID: "t", FileUniqueID: "tu", Width: 1, Height: 2, FileSize: 3}
	probes := []fieldErrorProbe{
		{User{ID: 1, FirstName: "A", LastName: "B", Username: "u", LanguageCode: "en", IsBot: true, IsPremium: true, AddedToAttachmentMenu: true, CanJoinGroups: true, CanReadAllGroupMessages: true, SupportsInlineQueries: true}, func() jsonDecoder { return &User{} }},
		{Chat{ID: 1, Type: ChatTypeGroup, Title: "t", Username: "u", FirstName: "f", LastName: "l", IsForum: true}, func() jsonDecoder { return &Chat{} }},
		{PhotoSize{FileID: "p", FileUniqueID: "pu", Width: 1, Height: 2, FileSize: 3}, func() jsonDecoder { return &PhotoSize{} }},
		{Animation{FileID: "a", FileUniqueID: "au", Width: 1, Height: 2, Duration: 3, Thumbnail: thumb, FileName: "a", MIMEType: "m", FileSize: 4}, func() jsonDecoder { return &Animation{} }},
		{Audio{FileID: "b", FileUniqueID: "bu", Duration: 5, Performer: "p", Title: "t", FileName: "f", MIMEType: "m", FileSize: 6, Thumbnail: thumb}, func() jsonDecoder { return &Audio{} }},
		{Document{FileID: "c", FileUniqueID: "cu", Thumbnail: thumb, FileName: "f", MIMEType: "m", FileSize: 7}, func() jsonDecoder { return &Document{} }},
		{Video{FileID: "v", FileUniqueID: "vu", Width: 1, Height: 2, Duration: 3, Thumbnail: thumb, FileName: "f", MIMEType: "m", FileSize: 8}, func() jsonDecoder { return &Video{} }},
		{VideoNote{FileID: "vn", FileUniqueID: "vnu", Length: 1, Duration: 2, Thumbnail: thumb, FileSize: 9}, func() jsonDecoder { return &VideoNote{} }},
		{Voice{FileID: "vo", FileUniqueID: "vou", Duration: 3, MIMEType: "m", FileSize: 10}, func() jsonDecoder { return &Voice{} }},
		{Sticker{FileID: "s", FileUniqueID: "su", Type: StickerMask, Width: 1, Height: 2, IsAnimated: true, IsVideo: true, Thumbnail: thumb, Emoji: "x", SetName: "set", FileSize: 11}, func() jsonDecoder { return &Sticker{} }},
		{MessageEntity{Type: EntityTextLink, Offset: 1, Length: 2, URL: "u", User: &User{ID: 1, FirstName: "A"}, Language: "go", CustomEmojiID: "e"}, func() jsonDecoder { return &MessageEntity{} }},
		{Contact{PhoneNumber: "+1", FirstName: "A", LastName: "B", UserID: 1, VCard: "v"}, func() jsonDecoder { return &Contact{} }},
		{Dice{Emoji: DiceBasketball, Value: 5}, func() jsonDecoder { return &Dice{} }},
		{Location{Longitude: 1.5, Latitude: 2.25, HorizontalAccuracy: 0.5, LivePeriod: 60, Heading: 90, ProximityAlertRadius: 10}, func() jsonDecoder { return &Location{} }},
		{Venue{Location: Location{Longitude: 1.5, Latitude: 2.25}, Title: "t", Address: "a", FoursquareID: "f", FoursquareType: "ft", GooglePlaceID: "g", GooglePlaceType: "gt"}, func() jsonDecoder { return &Venue{} }},
		{PollOption{Text: "a", VoterCount: 3}, func() jsonDecoder { return &PollOption{} }},
		{Poll{ID: "p", Question: "q", Options: []PollOption{{Text: "a", VoterCount: 1}}, TotalVoterCount: 1, Type: PollRegular, IsClosed: true, IsAnonymous: true, AllowsMultipleAnswers: true, CorrectOptionID: 1, Explanation: "e", ExplanationEntities: []MessageEntity{{Type: EntityCode, Length: 1}}, OpenPeriod: 30}, func() jsonDecoder { return &Poll{} }},
		{WebAppInfo{URL: "https://x"}, func() jsonDecoder { return &WebAppInfo{} }},
		{InlineKeyboardButton{Text: "t", URL: "u", CallbackData: "d", WebApp: &WebAppInfo{URL: "w"}, SwitchInlineQuery: strptr("q"), SwitchInlineQueryCurrentChat: strptr("c"), Pay: true}, func() jsonDecoder { return &InlineKeyboardButton{} }},
		{InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{{Text: "a"}}}}, func() jsonDecoder { return &InlineKeyboardMarkup{} }},
		{&MessageOriginUser{Type: OriginUser, Date: 1, SenderUser: User{ID: 1, FirstName: "A"}}, func() jsonDecoder { return &MessageOriginUser{} }},
		{&MessageOriginHiddenUser{Type: OriginHiddenUser, Date: 2, SenderUserName: "g"}, func() jsonDecoder { return &MessageOriginHiddenUser{} }},
		{&MessageOriginChat{Type: OriginChat, Date: 3, SenderChat: Chat{ID: 1, Type: ChatTypeChannel}, AuthorSignature: "s"}, func() jsonDecoder { return &MessageOriginChat{} }},
		{&MessageOriginChannel{Type: OriginChannel, Date: 4, Chat: Chat{ID: 1, Type: ChatTypeChannel}, MessageID: 5, AuthorSignature: "s"}, func() jsonDecoder { return &MessageOriginChannel{} }},
		{fullMessage(), func() jsonDecoder { return &Message{} }},
	}
	for _, p := range probes {
		data, err := json.Marshal(p.full)
		if err != nil {
			t.Fatalf("marshal %T: %v", p.full, err)
		}
		// An unknown field must be skipped, not rejected (covers the default
		// branch of every Decode).
		if err := unmarshalJX([]byte(`{"__unknown__":{"a":[1,2,3]}}`), p.newDec()); err != nil {
			t.Errorf("%T: unknown field not skipped: %v", p.full, err)
		}
		var keyed map[string]json.RawMessage
		if err := json.Unmarshal(data, &keyed); err != nil {
			t.Fatalf("rekey %T: %v", p.full, err)
		}
		for field := range keyed {
			// "@" is not a valid JSON value start, so whichever decoder the
			// field dispatches to (scalar, nested object or array) fails.
			bad := []byte(`{"` + field + `":@}`)
			if err := unmarshalJX(bad, p.newDec()); err == nil {
				t.Errorf("%T field %q: expected decode error", p.full, field)
			}
		}
	}
}

// TestEncodeMinimal encodes media entities with all optional fields zeroed,
// covering the absent-thumbnail and absent-optional encode paths.
func TestEncodeMinimal(t *testing.T) {
	jsonRoundTrip(t, Animation{FileID: "a", FileUniqueID: "au", Width: 1, Height: 2, Duration: 3})
	jsonRoundTrip(t, Audio{FileID: "b", FileUniqueID: "bu", Duration: 5})
	jsonRoundTrip(t, Document{FileID: "c", FileUniqueID: "cu"})
	jsonRoundTrip(t, Video{FileID: "v", FileUniqueID: "vu", Width: 1, Height: 2, Duration: 3})
	jsonRoundTrip(t, VideoNote{FileID: "vn", FileUniqueID: "vnu", Length: 1, Duration: 2})
	jsonRoundTrip(t, Voice{FileID: "vo", FileUniqueID: "vou", Duration: 3})
	jsonRoundTrip(t, Sticker{FileID: "s", FileUniqueID: "su", Type: StickerRegular, Width: 1, Height: 2})
	jsonRoundTrip(t, PhotoSize{FileID: "p", FileUniqueID: "pu", Width: 1, Height: 2})
	jsonRoundTrip(t, MessageEntity{Type: EntityBold, Offset: 0, Length: 1})
	jsonRoundTrip(t, Contact{PhoneNumber: "+1", FirstName: "A"})
	jsonRoundTrip(t, Location{Longitude: 1.5, Latitude: 2.25})
	jsonRoundTrip(t, Venue{Location: Location{Longitude: 1.5, Latitude: 2.25}, Title: "t", Address: "a"})
	jsonRoundTrip(t, Poll{ID: "p", Question: "q", Options: []PollOption{{Text: "a", VoterCount: 1}}, TotalVoterCount: 1, Type: PollRegular})
	jsonRoundTrip(t, InlineKeyboardButton{Text: "t"})
	jsonRoundTrip(t, &MessageOriginChat{Type: OriginChat, Date: 1, SenderChat: Chat{ID: 1, Type: ChatTypeChannel}})
	jsonRoundTrip(t, &MessageOriginChannel{Type: OriginChannel, Date: 1, Chat: Chat{ID: 1, Type: ChatTypeChannel}, MessageID: 2})
	jsonRoundTrip(t, &Message{MessageID: 1, Date: 1, Chat: Chat{ID: 1, Type: ChatTypePrivate}})
}

func TestMessageOriginVariantsRoundTrip(t *testing.T) {
	chat := Chat{ID: 1, Type: ChatTypeChannel}
	cases := []MessageOrigin{
		&MessageOriginUser{Type: OriginUser, Date: 1, SenderUser: User{ID: 1, FirstName: "A"}},
		&MessageOriginHiddenUser{Type: OriginHiddenUser, Date: 2, SenderUserName: "ghost"},
		&MessageOriginChat{Type: OriginChat, Date: 3, SenderChat: chat, AuthorSignature: "s"},
		&MessageOriginChannel{Type: OriginChannel, Date: 4, Chat: chat, MessageID: 5},
	}
	for _, want := range cases {
		var e jx.Encoder
		want.Encode(&e)
		got, err := decodeMessageOrigin(jx.DecodeBytes(e.Bytes()))
		if err != nil {
			t.Fatalf("decode %T: %v", want, err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("origin %T mismatch:\nwant %+v\n got %+v", want, want, got)
		}
	}
}

func TestDecodeMessageOriginUnknown(t *testing.T) {
	if _, err := decodeMessageOrigin(jx.DecodeStr(`{"type":"bogus"}`)); err == nil {
		t.Fatal("expected error for unknown origin type")
	}
}

// TestEntitiesUnmarshalJSON checks that the standalone entity types parse from
// raw Bot API JSON through their UnmarshalJSON methods.
func TestEntitiesUnmarshalJSON(t *testing.T) {
	t.Run("User", func(t *testing.T) {
		var u User
		if err := json.Unmarshal([]byte(`{"id":42,"first_name":"Ada","is_bot":true}`), &u); err != nil {
			t.Fatal(err)
		}
		if u.ID != 42 || u.FirstName != "Ada" || !u.IsBot {
			t.Fatalf("user = %+v", u)
		}
	})
	t.Run("Chat", func(t *testing.T) {
		var c Chat
		if err := json.Unmarshal([]byte(`{"id":-1,"type":"private","first_name":"X"}`), &c); err != nil {
			t.Fatal(err)
		}
		if c.ID != -1 || c.Type != ChatTypePrivate {
			t.Fatalf("chat = %+v", c)
		}
	})
	t.Run("Poll", func(t *testing.T) {
		var p Poll
		if err := json.Unmarshal([]byte(`{"id":"x","question":"q","options":[{"text":"a","voter_count":1}],"total_voter_count":1,"type":"regular"}`), &p); err != nil {
			t.Fatal(err)
		}
		if len(p.Options) != 1 || p.Options[0].Text != "a" {
			t.Fatalf("poll = %+v", p)
		}
	})
	t.Run("UnknownFieldSkipped", func(t *testing.T) {
		var u User
		if err := json.Unmarshal([]byte(`{"id":1,"first_name":"A","unexpected":{"a":[1,2]}}`), &u); err != nil {
			t.Fatal(err)
		}
		if u.ID != 1 {
			t.Fatalf("user = %+v", u)
		}
	})
}

// TestEntityDecodeErrors ensures malformed input surfaces a decode error rather
// than silently succeeding.
func TestEntityDecodeErrors(t *testing.T) {
	cases := []struct {
		name string
		dec  jsonDecoder
		data string
	}{
		{"user-id-type", &User{}, `{"id":"notnumber"}`},
		{"chat-truncated", &Chat{}, `{"id":1`},
		{"message-bad-entities", &Message{}, `{"entities":[{"offset":"x"}]}`},
		{"poll-bad-options", &Poll{}, `{"options":[{"voter_count":"x"}]}`},
		{"markup-bad", &InlineKeyboardMarkup{}, `{"inline_keyboard":[[{"text":1}]]}`},
		{"origin-bad", &MessageOriginUser{}, `{"date":"x"}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := unmarshalJX([]byte(c.data), c.dec); err == nil {
				t.Fatalf("expected error decoding %s", c.data)
			}
		})
	}
}
