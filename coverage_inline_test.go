package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

// inlineToTg is the (ctx, *Bot) conversion every inline-result variant exposes.
type inlineToTg interface {
	toTg(context.Context, *Bot) (tg.InputBotInlineResultClass, error)
}

// badInlineMarkup is a keyboard whose single button carries no action, so
// inlineButtonToTg rejects it — used to drive the markup-error path of each
// result's toTg.
var badInlineMarkup = &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{{Text: "x"}}}}

// goodInlineMarkup is a valid one-button keyboard.
var goodInlineMarkup = &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{{Text: "x", CallbackData: "d"}}}}

func TestInlineResultMarkupErrors(t *testing.T) {
	b := newMockBot(newMockInvoker())
	doc := documentFileID(t, 0x77)
	photo := photoFileID(t, 0x88)
	content := &InputTextMessageContent{MessageText: "hi"}

	cases := []inlineToTg{
		&InlineQueryResultArticle{ID: "1", Title: "t", InputMessageContent: content, ReplyMarkup: badInlineMarkup},
		&InlineQueryResultPhoto{ID: "2", PhotoURL: "u", ReplyMarkup: badInlineMarkup},
		&InlineQueryResultGif{ID: "3", GifURL: "u", ReplyMarkup: badInlineMarkup},
		&InlineQueryResultMpeg4Gif{ID: "4", Mpeg4URL: "u", ReplyMarkup: badInlineMarkup},
		&InlineQueryResultContact{ID: "5", PhoneNumber: "+1", FirstName: "A", ReplyMarkup: badInlineMarkup},
		&InlineQueryResultLocation{ID: "6", ReplyMarkup: badInlineMarkup},
		&InlineQueryResultVenue{ID: "7", ReplyMarkup: badInlineMarkup},
		&InlineQueryResultCachedPhoto{ID: "8", PhotoFileID: photo, ReplyMarkup: badInlineMarkup},
		&InlineQueryResultCachedDocument{ID: "9", DocumentFileID: doc, ReplyMarkup: badInlineMarkup},
		&InlineQueryResultCachedGif{ID: "10", GifFileID: doc, ReplyMarkup: badInlineMarkup},
		&InlineQueryResultCachedVideo{ID: "11", VideoFileID: doc, ReplyMarkup: badInlineMarkup},
		&InlineQueryResultCachedVoice{ID: "12", VoiceFileID: doc, ReplyMarkup: badInlineMarkup},
		&InlineQueryResultCachedAudio{ID: "13", AudioFileID: doc, ReplyMarkup: badInlineMarkup},
		&InlineQueryResultCachedSticker{ID: "14", StickerFileID: doc, ReplyMarkup: badInlineMarkup},
	}
	for _, c := range cases {
		if _, err := c.toTg(context.Background(), b); err == nil {
			t.Errorf("%T: expected reply markup error", c)
		}
	}
}

func TestInlineResultWithContentAndMarkup(t *testing.T) {
	b := newMockBot(newMockInvoker())
	content := &InputTextMessageContent{MessageText: "hi"}

	// Contact/Location/Venue with explicit content take the content branch;
	// with a thumbnail and a valid markup they take the media + thumb branch.
	cases := []inlineToTg{
		&InlineQueryResultContact{ID: "a", PhoneNumber: "+1", FirstName: "A", InputMessageContent: content},
		&InlineQueryResultContact{ID: "b", PhoneNumber: "+1", FirstName: "A", LastName: "B", Vcard: "v", ThumbnailURL: "https://e/t.jpg", ReplyMarkup: goodInlineMarkup},
		&InlineQueryResultLocation{ID: "c", Latitude: 1, Longitude: 2, InputMessageContent: content},
		&InlineQueryResultLocation{ID: "d", Latitude: 1, Longitude: 2, Heading: 90, LivePeriod: 60, ThumbnailURL: "https://e/t.jpg", ReplyMarkup: goodInlineMarkup},
		&InlineQueryResultVenue{ID: "e", Latitude: 1, Longitude: 2, Title: "t", Address: "a", InputMessageContent: content},
		&InlineQueryResultVenue{ID: "f", Latitude: 1, Longitude: 2, Title: "t", Address: "a", FoursquareID: "fid", ThumbnailURL: "https://e/t.jpg", ReplyMarkup: goodInlineMarkup},
	}
	for _, c := range cases {
		if _, err := c.toTg(context.Background(), b); err != nil {
			t.Errorf("%T: %v", c, err)
		}
	}
}

func TestInputMessageContentToTg(t *testing.T) {
	b := newMockBot(newMockInvoker())
	mkp, err := replyMarkupToTg(goodInlineMarkup)
	if err != nil {
		t.Fatal(err)
	}
	contents := []InputMessageContent{
		&InputTextMessageContent{MessageText: "t", Entities: []MessageEntity{{Type: EntityBold, Length: 1}}},
		&InputLocationMessageContent{Latitude: 1, Longitude: 2, HorizontalAccuracy: 5, Heading: 90, LivePeriod: 60},
		&InputVenueMessageContent{Latitude: 1, Longitude: 2, Title: "t", Address: "a", FoursquareID: "f", FoursquareType: "ft"},
		&InputVenueMessageContent{Latitude: 1, Longitude: 2, Title: "t", Address: "a", GooglePlaceID: "g", GooglePlaceType: "gt"},
		&InputVenueMessageContent{Latitude: 1, Longitude: 2, Title: "t", Address: "a"},
		&InputContactMessageContent{PhoneNumber: "+1", FirstName: "A", LastName: "B", Vcard: "v"},
	}
	for _, c := range contents {
		if _, err := c.toTg(context.Background(), b, mkp); err != nil {
			t.Errorf("%T: %v", c, err)
		}
	}
}

// TestInputTextContentParseMode covers the parse-mode styling branch of text
// content conversion.
func TestInputTextContentParseMode(t *testing.T) {
	b := newMockBot(newMockInvoker())
	c := &InputTextMessageContent{MessageText: "<b>bold</b>", ParseMode: ParseModeHTML, DisableWebPagePreview: true}
	got, err := c.toTg(context.Background(), b, nil)
	if err != nil {
		t.Fatal(err)
	}
	text, ok := got.(*tg.InputBotInlineMessageText)
	if !ok || text.Message != "bold" || len(text.Entities) != 1 {
		t.Fatalf("text content = %#v", got)
	}
}

func TestVenueProvider(t *testing.T) {
	cases := []struct {
		content *InputVenueMessageContent
		want    string
	}{
		{&InputVenueMessageContent{FoursquareID: "x"}, "foursquare"},
		{&InputVenueMessageContent{FoursquareType: "x"}, "foursquare"},
		{&InputVenueMessageContent{GooglePlaceID: "x"}, "gplaces"},
		{&InputVenueMessageContent{GooglePlaceType: "x"}, "gplaces"},
		{&InputVenueMessageContent{}, ""},
	}
	for _, c := range cases {
		if got := venueProvider(c.content); got != c.want {
			t.Errorf("venueProvider(%+v) = %q, want %q", c.content, got, c.want)
		}
	}
}
