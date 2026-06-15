package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestInlineMessageIDRoundTrip(t *testing.T) {
	cases := []tg.InputBotInlineMessageIDClass{
		&tg.InputBotInlineMessageID{DCID: 2, ID: 123456789, AccessHash: -987654321},
		&tg.InputBotInlineMessageID64{DCID: 4, OwnerID: 42, ID: 7, AccessHash: 555},
	}

	for _, want := range cases {
		enc, err := encodeInlineMessageID(want)
		if err != nil {
			t.Fatalf("encode %#v: %v", want, err)
		}

		if enc == "" {
			t.Fatalf("empty encoding for %#v", want)
		}

		got, err := decodeInlineMessageID(enc)
		if err != nil {
			t.Fatalf("decode %q: %v", enc, err)
		}

		if got.String() != want.String() {
			t.Fatalf("round trip: got %#v, want %#v", got, want)
		}
	}
}

func TestEncodeInlineMessageIDNil(t *testing.T) {
	enc, err := encodeInlineMessageID(nil)
	if err != nil {
		t.Fatalf("encode nil: %v", err)
	}

	if enc != "" {
		t.Fatalf("nil encoding = %q, want empty", enc)
	}
}

func TestDecodeInlineMessageIDInvalid(t *testing.T) {
	for _, s := range []string{"", "not base64 !!!", "AAAA"} {
		if _, err := decodeInlineMessageID(s); err == nil {
			t.Fatalf("expected error for %q", s)
		}
	}
}

func TestChosenInlineResultFromTg(t *testing.T) {
	mid := &tg.InputBotInlineMessageID64{DCID: 2, OwnerID: 10, ID: 3, AccessHash: 99}

	u := &tg.UpdateBotInlineSend{
		UserID: 10,
		Query:  "q",
		ID:     "result-1",
	}
	u.SetMsgID(mid)
	u.SetGeo(&tg.GeoPoint{Lat: 1.5, Long: 2.5, AccuracyRadius: 7})

	e := tg.Entities{Users: map[int64]*tg.User{10: {ID: 10, FirstName: "Picker"}}}

	got := chosenInlineResultFromTg(e, u)

	if got.ResultID != "result-1" || got.Query != "q" || got.From.FirstName != "Picker" {
		t.Fatalf("chosen result = %#v", got)
	}

	if got.Location == nil || got.Location.Latitude != 1.5 || got.Location.Longitude != 2.5 {
		t.Fatalf("location = %#v", got.Location)
	}

	decoded, err := decodeInlineMessageID(got.InlineMessageID)
	if err != nil {
		t.Fatalf("decode inline message id: %v", err)
	}

	if decoded.String() != mid.String() {
		t.Fatalf("inline message id = %#v, want %#v", decoded, mid)
	}
}

func TestInlineCallbackQueryFromTg(t *testing.T) {
	mid := &tg.InputBotInlineMessageID{DCID: 1, ID: 5, AccessHash: 6}

	u := &tg.UpdateInlineBotCallbackQuery{
		QueryID:      77,
		UserID:       10,
		MsgID:        mid,
		ChatInstance: 88,
		Data:         []byte("payload"),
	}

	e := tg.Entities{Users: map[int64]*tg.User{10: {ID: 10, FirstName: "Tapper"}}}

	got := inlineCallbackQueryFromTg(e, u)

	if got.ID != "77" || got.Data != "payload" || got.From.FirstName != "Tapper" {
		t.Fatalf("callback = %#v", got)
	}

	if got.Message != nil {
		t.Fatalf("inline callback should have no Message")
	}

	decoded, err := decodeInlineMessageID(got.InlineMessageID)
	if err != nil {
		t.Fatalf("decode inline message id: %v", err)
	}

	if decoded.String() != mid.String() {
		t.Fatalf("inline message id = %#v, want %#v", decoded, mid)
	}
}
