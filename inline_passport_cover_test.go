package botapi

import (
	"context"
	"errors"
	"testing"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

// photoFileID builds a valid photo file_id, for cached-photo results.
func photoFileID(t *testing.T, id int64) string {
	t.Helper()
	s, err := fileid.EncodeFileID(fileid.FileID{Type: fileid.Photo, DC: 2, ID: id, AccessHash: 7})
	if err != nil {
		t.Fatalf("encode photo file id: %v", err)
	}
	return s
}

// TestInlineResultsToTg exercises toTg for every inline result variant.
func TestInlineResultsToTg(t *testing.T) {
	b := &Bot{}
	doc := documentFileID(t, 0x55)
	photo := photoFileID(t, 0x66)

	cases := []struct {
		name    string
		result  InlineQueryResult
		wantTyp string
	}{
		{"gif-url", &InlineQueryResultGif{ID: "1", GifURL: "https://e/g.gif", ThumbnailURL: "https://e/t.jpg"}, "gif"},
		{"mpeg4-url", &InlineQueryResultMpeg4Gif{ID: "2", Mpeg4URL: "https://e/m.mp4", ThumbnailURL: "https://e/t.jpg"}, "gif"},
		{"contact", &InlineQueryResultContact{ID: "3", PhoneNumber: "+1", FirstName: "Ada"}, "contact"},
		{"location", &InlineQueryResultLocation{ID: "4", Latitude: 1, Longitude: 2, Title: "L"}, "geo"},
		{"venue", &InlineQueryResultVenue{ID: "5", Latitude: 1, Longitude: 2, Title: "V", Address: "A"}, "venue"},
		{"cached-photo", &InlineQueryResultCachedPhoto{ID: "6", PhotoFileID: photo}, "photo"},
		{"cached-doc", &InlineQueryResultCachedDocument{ID: "7", Title: "T", DocumentFileID: doc}, "file"},
		{"cached-gif", &InlineQueryResultCachedGif{ID: "8", GifFileID: doc}, "gif"},
		{"cached-video", &InlineQueryResultCachedVideo{ID: "9", VideoFileID: doc, Title: "T"}, "video"},
		{"cached-voice", &InlineQueryResultCachedVoice{ID: "10", VoiceFileID: doc, Title: "T"}, "voice"},
		{"cached-audio", &InlineQueryResultCachedAudio{ID: "11", AudioFileID: doc}, "audio"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.result.toTg(context.Background(), b)
			if err != nil {
				t.Fatalf("toTg: %v", err)
			}
			switch r := got.(type) {
			case *tg.InputBotInlineResult:
				if r.Type != c.wantTyp {
					t.Fatalf("type = %q, want %q", r.Type, c.wantTyp)
				}
			case *tg.InputBotInlineResultPhoto:
				if r.Type != c.wantTyp {
					t.Fatalf("type = %q, want %q", r.Type, c.wantTyp)
				}
			case *tg.InputBotInlineResultDocument:
				if r.Type != c.wantTyp {
					t.Fatalf("type = %q, want %q", r.Type, c.wantTyp)
				}
			default:
				t.Fatalf("unexpected %T", got)
			}
		})
	}
}

// TestAnswerInlineQueryDispatch drives AnswerInlineQuery end to end, exercising
// the request assembly over a mix of result types.
func TestAnswerInlineQueryDispatch(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetInlineBotResultsRequestTypeID, &tg.BoolTrue{})
	b := newMockBot(inv)

	results := []InlineQueryResult{
		&InlineQueryResultArticle{ID: "1", Title: "A", InputMessageContent: &InputTextMessageContent{MessageText: "x"}},
		&InlineQueryResultContact{ID: "2", PhoneNumber: "+1", FirstName: "Ada"},
	}
	err := b.AnswerInlineQuery(context.Background(), "777", results,
		WithInlineCacheTime(10), WithInlinePersonal(), WithInlineNextOffset("next"))
	if err != nil {
		t.Fatalf("AnswerInlineQuery: %v", err)
	}
	var req tg.MessagesSetInlineBotResultsRequest
	inv.decode(t, tg.MessagesSetInlineBotResultsRequestTypeID, &req)
	if req.QueryID != 777 || len(req.Results) != 2 {
		t.Fatalf("req = %#v", req)
	}
	if req.CacheTime != 10 || !req.Private || req.NextOffset != "next" {
		t.Fatalf("req opts = %#v", req)
	}
}

func TestAnswerInlineQueryInvalidID(t *testing.T) {
	b := newMockBot(newMockInvoker())
	err := b.AnswerInlineQuery(context.Background(), "not-a-number", nil)
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.Code != 400 {
		t.Fatalf("want 400, got %v", err)
	}
}

func TestInputContentVariantsToTg(t *testing.T) {
	ctx := context.Background()

	loc := &InputLocationMessageContent{Latitude: 1, Longitude: 2}
	if got, err := loc.toTg(ctx, nil, nil); err != nil {
		t.Fatalf("location: %v", err)
	} else if _, ok := got.(*tg.InputBotInlineMessageMediaGeo); !ok {
		t.Fatalf("location = %T", got)
	}

	venue := &InputVenueMessageContent{Latitude: 1, Longitude: 2, Title: "V", Address: "A"}
	if got, err := venue.toTg(ctx, nil, nil); err != nil {
		t.Fatalf("venue: %v", err)
	} else if _, ok := got.(*tg.InputBotInlineMessageMediaVenue); !ok {
		t.Fatalf("venue = %T", got)
	}

	text := &InputTextMessageContent{MessageText: "hi"}
	if got, err := text.toTg(ctx, &Bot{}, nil); err != nil {
		t.Fatalf("text: %v", err)
	} else if _, ok := got.(*tg.InputBotInlineMessageText); !ok {
		t.Fatalf("text = %T", got)
	}
}

// TestPassportErrorsToTg drives SetPassportDataErrors with one of every passport
// error variant, exercising each toTg.
func TestPassportErrorsToTg(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.UsersSetSecureValueErrorsRequestTypeID, &tg.BoolTrue{})
	b := newMockBot(inv)

	const h = "aGFzaGhhc2hoYXNoaGFzaA==" // base64, 16+ bytes
	errs := []PassportElementError{
		&PassportElementErrorDataField{Type: "personal_details", FieldName: "first_name", DataHash: h, Message: "m"},
		&PassportElementErrorFrontSide{Type: "passport", FileHash: h, Message: "m"},
		&PassportElementErrorReverseSide{Type: "driver_license", FileHash: h, Message: "m"},
		&PassportElementErrorSelfie{Type: "passport", FileHash: h, Message: "m"},
		&PassportElementErrorFile{Type: "utility_bill", FileHash: h, Message: "m"},
		&PassportElementErrorFiles{Type: "utility_bill", FileHashes: []string{h}, Message: "m"},
		&PassportElementErrorTranslationFile{Type: "passport", FileHash: h, Message: "m"},
		&PassportElementErrorTranslationFiles{Type: "passport", FileHashes: []string{h}, Message: "m"},
		&PassportElementErrorUnspecified{Type: "passport", ElementHash: h, Message: "m"},
	}
	if err := b.SetPassportDataErrors(context.Background(), 99, errs); err != nil {
		t.Fatalf("SetPassportDataErrors: %v", err)
	}
	var req tg.UsersSetSecureValueErrorsRequest
	inv.decode(t, tg.UsersSetSecureValueErrorsRequestTypeID, &req)
	if len(req.Errors) != len(errs) {
		t.Fatalf("errors = %d, want %d", len(req.Errors), len(errs))
	}
}
