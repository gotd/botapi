package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func testSentMedia(a *require.Assertions, mock *tgmock.Mock, media tg.InputMediaClass) {
	mock.ExpectFunc(func(b bin.Encoder) {
		r := b.(*tg.MessagesSendMediaRequest)
		a.Equal(&tg.InputPeerChat{ChatID: testChat().ID}, r.Peer)

		setFlags(media)
		setFlags(r.Media)
		a.Equal(media, r.Media)
	}).ThenResult(&tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateNewMessage{
				Message: &tg.Message{
					Out:    false,
					ID:     10,
					PeerID: &tg.PeerChat{ChatID: testChat().ID},
				},
			},
		},
	})
}

func TestBotAPI_SendContact(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		testSentMedia(a, mock, &tg.InputMediaContact{
			PhoneNumber: "aboba",
			FirstName:   "aboba",
			LastName:    "aboba",
			Vcard:       "aboba",
		})

		msg, err := api.SendContact(ctx, &oas.SendContact{
			ChatID:      oas.NewInt64ID(testChatID()),
			PhoneNumber: "aboba",
			FirstName:   "aboba",
			LastName:    oas.NewOptString("aboba"),
			Vcard:       oas.NewOptString("aboba"),
		})
		a.NoError(err)
		a.Equal(10, msg.Result.Value.MessageID)
	})
}

func TestBotAPI_SendDice(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		testSendDice := func(expect string, input oas.OptString) {
			testSentMedia(a, mock, &tg.InputMediaDice{
				Emoticon: expect,
			})

			msg, err := api.SendDice(ctx, &oas.SendDice{
				ChatID: oas.NewInt64ID(testChatID()),
				Emoji:  input,
			})
			a.NoError(err)
			a.Equal(10, msg.Result.Value.MessageID)
		}

		// Ensure setting default.
		testSendDice(message.DiceEmoticon, oas.OptString{})
		testSendDice(message.BowlingEmoticon, oas.NewOptString(message.BowlingEmoticon))
	})
}

func TestBotAPI_SendInvoice(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		invoice := tg.Invoice{
			Test:                     false,
			NameRequested:            true,
			PhoneRequested:           true,
			EmailRequested:           true,
			ShippingAddressRequested: true,
			Flexible:                 true,
			PhoneToProvider:          true,
			EmailToProvider:          true,
			Currency:                 "currency",
			Prices: []tg.LabeledPrice{
				{
					Label:  "label",
					Amount: 10,
				},
			},
			MaxTipAmount:        10,
			SuggestedTipAmounts: []int64{1, 2, 3},
		}
		invoice.SetFlags()
		testSentMedia(a, mock, &tg.InputMediaInvoice{
			Title:       "title",
			Description: "description",
			Photo: tg.InputWebDocument{
				URL:      "photo URL",
				Size:     10,
				MimeType: "image/jpeg",
				Attributes: []tg.DocumentAttributeClass{
					&tg.DocumentAttributeImageSize{
						W: 10,
						H: 10,
					},
				},
			},
			Invoice:  invoice,
			Payload:  []byte(`payload`),
			Provider: "provider",
			ProviderData: tg.DataJSON{
				Data: "provider data",
			},
			StartParam: "start parameter",
		})

		msg, err := api.SendInvoice(ctx, &oas.SendInvoice{
			ChatID:        oas.NewInt64ID(testChatID()),
			Title:         "title",
			Description:   "description",
			Payload:       "payload",
			ProviderToken: "provider",
			Currency:      "currency",
			Prices: []oas.LabeledPrice{
				{
					Label:  "label",
					Amount: 10,
				},
			},
			MaxTipAmount:              oas.NewOptInt(10),
			SuggestedTipAmounts:       []int64{1, 2, 3},
			StartParameter:            oas.NewOptString("start parameter"),
			ProviderData:              oas.NewOptString("provider data"),
			PhotoURL:                  oas.NewOptString("photo URL"),
			PhotoSize:                 oas.NewOptInt(10),
			PhotoWidth:                oas.NewOptInt(10),
			PhotoHeight:               oas.NewOptInt(10),
			NeedName:                  oas.NewOptBool(true),
			NeedPhoneNumber:           oas.NewOptBool(true),
			NeedEmail:                 oas.NewOptBool(true),
			NeedShippingAddress:       oas.NewOptBool(true),
			SendPhoneNumberToProvider: oas.NewOptBool(true),
			SendEmailToProvider:       oas.NewOptBool(true),
			IsFlexible:                oas.NewOptBool(true),
			DisableNotification:       oas.OptBool{},
			ProtectContent:            oas.OptBool{},
			ReplyToMessageID:          oas.OptInt{},
			AllowSendingWithoutReply:  oas.OptBool{},
			ReplyMarkup:               oas.OptInlineKeyboardMarkup{},
		})
		a.NoError(err)
		a.Equal(10, msg.Result.Value.MessageID)
	})
}

func TestBotAPI_SendLocation(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		testSendLocation := func(expect tg.InputMediaClass, input oas.SendLocation) {
			testSentMedia(a, mock, expect)

			msg, err := api.SendLocation(ctx, &input)
			a.NoError(err)
			a.Equal(10, msg.Result.Value.MessageID)
		}

		p := &tg.InputGeoPoint{
			Lat:            10,
			Long:           10,
			AccuracyRadius: 10,
		}
		p.SetFlags()
		testSendLocation(&tg.InputMediaGeoPoint{
			GeoPoint: p,
		}, oas.SendLocation{
			ChatID:             oas.NewInt64ID(testChatID()),
			Latitude:           10,
			Longitude:          10,
			HorizontalAccuracy: oas.NewOptFloat64(10),
		})
		testSendLocation(&tg.InputMediaGeoLive{
			Stopped:                     false,
			GeoPoint:                    p,
			Heading:                     10,
			Period:                      10,
			ProximityNotificationRadius: 10,
		}, oas.SendLocation{
			ChatID:                   oas.NewInt64ID(testChatID()),
			Latitude:                 10,
			Longitude:                10,
			HorizontalAccuracy:       oas.NewOptFloat64(10),
			LivePeriod:               oas.NewOptInt(10),
			Heading:                  oas.NewOptInt(10),
			ProximityAlertRadius:     oas.NewOptInt(10),
			DisableNotification:      oas.OptBool{},
			ProtectContent:           oas.OptBool{},
			ReplyToMessageID:         oas.OptInt{},
			AllowSendingWithoutReply: oas.OptBool{},
			ReplyMarkup:              oas.OptSendReplyMarkup{},
		})
	})
}

func TestBotAPI_SendVenue(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		type testOptions struct {
			FoursquareID    oas.OptString
			FoursquareType  oas.OptString
			GooglePlaceID   oas.OptString
			GooglePlaceType oas.OptString
		}
		testSendVenue := func(expect *tg.InputMediaVenue, opts testOptions) {
			testSentMedia(a, mock, expect)

			msg, err := api.SendVenue(ctx, &oas.SendVenue{
				ChatID:                   oas.NewInt64ID(testChatID()),
				Latitude:                 10,
				Longitude:                10,
				Title:                    "title",
				Address:                  "address",
				FoursquareID:             opts.FoursquareID,
				FoursquareType:           opts.FoursquareType,
				GooglePlaceID:            opts.GooglePlaceID,
				GooglePlaceType:          opts.GooglePlaceType,
				DisableNotification:      oas.OptBool{},
				ProtectContent:           oas.OptBool{},
				ReplyToMessageID:         oas.OptInt{},
				AllowSendingWithoutReply: oas.OptBool{},
				ReplyMarkup:              oas.OptSendReplyMarkup{},
			})
			a.NoError(err)
			a.Equal(10, msg.Result.Value.MessageID)
		}

		p := &tg.InputGeoPoint{
			Lat:  10,
			Long: 10,
		}
		p.SetFlags()
		testSendVenue(&tg.InputMediaVenue{
			GeoPoint:  p,
			Title:     "title",
			Address:   "address",
			Provider:  "foursquare",
			VenueID:   "venue_id",
			VenueType: "venue_type",
		}, testOptions{
			FoursquareID:   oas.NewOptString("venue_id"),
			FoursquareType: oas.NewOptString("venue_type"),
		})
		testSendVenue(&tg.InputMediaVenue{
			GeoPoint:  p,
			Title:     "title",
			Address:   "address",
			Provider:  "gplaces",
			VenueID:   "venue_id",
			VenueType: "venue_type",
		}, testOptions{
			GooglePlaceID:   oas.NewOptString("venue_id"),
			GooglePlaceType: oas.NewOptString("venue_type"),
		})
		testSendVenue(&tg.InputMediaVenue{
			GeoPoint:  p,
			Title:     "title",
			Address:   "address",
			Provider:  "gplaces",
			VenueID:   "venue_id",
			VenueType: "venue_type",
		}, testOptions{
			FoursquareID:    oas.NewOptString("venue_id"),
			FoursquareType:  oas.NewOptString("venue_type"),
			GooglePlaceID:   oas.NewOptString("venue_id"),
			GooglePlaceType: oas.NewOptString("venue_type"),
		})
	})
}
