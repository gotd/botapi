package botapi

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

// SendContact implements oas.Handler.
func (b *BotAPI) SendContact(ctx context.Context, req oas.SendContact) (oas.ResultMessage, error) {
	s, p, err := b.prepareSend(
		ctx,
		sendOpts{
			To:                       req.ChatID,
			DisableNotification:      req.DisableNotification,
			ProtectContent:           req.ProtectContent,
			ReplyToMessageID:         req.ReplyToMessageID,
			AllowSendingWithoutReply: req.AllowSendingWithoutReply,
			ReplyMarkup:              req.ReplyMarkup,
		},
	)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "prepare send")
	}
	resp, err := s.Media(ctx, message.Contact(tg.InputMediaContact{
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName.Value,
		Vcard:       req.Vcard.Value,
	}))
	return b.sentMessage(ctx, p, resp, err)
}

// SendDice implements oas.Handler.
func (b *BotAPI) SendDice(ctx context.Context, req oas.SendDice) (oas.ResultMessage, error) {
	s, p, err := b.prepareSend(
		ctx,
		sendOpts{
			To:                       req.ChatID,
			DisableNotification:      req.DisableNotification,
			ProtectContent:           req.ProtectContent,
			ReplyToMessageID:         req.ReplyToMessageID,
			AllowSendingWithoutReply: req.AllowSendingWithoutReply,
			ReplyMarkup:              req.ReplyMarkup,
		},
	)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "prepare send")
	}
	resp, err := s.Media(ctx, message.MediaDice(req.Emoji.Or("🎲")))
	return b.sentMessage(ctx, p, resp, err)
}

func parseInlineKeyboardMarkup(r oas.OptInlineKeyboardMarkup) oas.OptSendReplyMarkup {
	var markup oas.OptSendReplyMarkup
	if m, ok := r.Get(); ok {
		markup.SetTo(oas.SendReplyMarkup{
			Type:                 oas.InlineKeyboardMarkupSendReplyMarkup,
			InlineKeyboardMarkup: m,
		})
	}
	return markup
}

// SendGame implements oas.Handler.
func (b *BotAPI) SendGame(ctx context.Context, req oas.SendGame) (oas.ResultMessage, error) {
	s, p, err := b.prepareSend(
		ctx,
		sendOpts{
			To:                       oas.NewInt64ID(req.ChatID),
			DisableNotification:      req.DisableNotification,
			ProtectContent:           req.ProtectContent,
			ReplyToMessageID:         req.ReplyToMessageID,
			AllowSendingWithoutReply: req.AllowSendingWithoutReply,
			ReplyMarkup:              parseInlineKeyboardMarkup(req.ReplyMarkup),
		},
	)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "prepare send")
	}
	resp, err := s.Media(ctx, message.Game(&tg.InputGameShortName{
		// TDLib uses self user.
		//
		// See https://github.com/tdlib/td/blob/fa8feefed70d64271945e9d5fd010b957d93c8cd/td/telegram/Game.cpp#L93.
		BotID:     &tg.InputUserSelf{},
		ShortName: req.GameShortName,
	}))
	return b.sentMessage(ctx, p, resp, err)
}

// SendInvoice implements oas.Handler.
func (b *BotAPI) SendInvoice(ctx context.Context, req oas.SendInvoice) (oas.ResultMessage, error) {
	s, p, err := b.prepareSend(
		ctx,
		sendOpts{
			To:                       req.ChatID,
			DisableNotification:      req.DisableNotification,
			ProtectContent:           req.ProtectContent,
			ReplyToMessageID:         req.ReplyToMessageID,
			AllowSendingWithoutReply: req.AllowSendingWithoutReply,
			ReplyMarkup:              parseInlineKeyboardMarkup(req.ReplyMarkup),
		},
	)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "prepare send")
	}
	invoice := tg.Invoice{
		Test:                     false,
		NameRequested:            req.NeedName.Value,
		PhoneRequested:           req.NeedPhoneNumber.Value,
		EmailRequested:           req.NeedEmail.Value,
		ShippingAddressRequested: req.NeedShippingAddress.Value,
		Flexible:                 req.IsFlexible.Value,
		PhoneToProvider:          req.SendPhoneNumberToProvider.Value,
		EmailToProvider:          req.SendEmailToProvider.Value,
		Currency:                 req.Currency,
		Prices:                   make([]tg.LabeledPrice, len(req.Prices)),
		MaxTipAmount:             0,
		SuggestedTipAmounts:      req.SuggestedTipAmounts,
	}
	invoice.SetFlags()
	{
		to := invoice.Prices
		from := req.Prices

		for i := range from {
			to[i] = tg.LabeledPrice{
				Label:  from[i].Label,
				Amount: int64(from[i].Amount),
			}
		}
		invoice.Prices = to
	}

	if v, ok := req.MaxTipAmount.Get(); ok {
		invoice.SetMaxTipAmount(int64(v))
	}
	media := &tg.InputMediaInvoice{
		Title:       req.Title,
		Description: req.Description,
		Photo:       tg.InputWebDocument{},
		Invoice:     invoice,
		Payload:     []byte(req.Payload),
		Provider:    req.ProviderToken,
		ProviderData: tg.DataJSON{
			Data: req.ProviderData.Value,
		},
	}
	if u, ok := req.PhotoURL.Get(); ok {
		doc := tg.InputWebDocument{
			URL:  u,
			Size: req.PhotoSize.Value,
			// TODO(tdakkota): Port TDLib extension parser.
			//
			// See https://github.com/tdlib/td/blob/fa8feefed70d64271945e9d5fd010b957d93c8cd/td/telegram/Payments.cpp#L877
			MimeType: "image/jpeg",
		}
		if w, h := req.PhotoWidth.Value, req.PhotoHeight.Value; w != 0 && h != 0 {
			doc.Attributes = append(doc.Attributes, &tg.DocumentAttributeImageSize{
				W: req.PhotoWidth.Value,
				H: req.PhotoHeight.Value,
			})
		}
		media.SetPhoto(doc)
	}
	if v, ok := req.StartParameter.Get(); ok {
		media.SetStartParam(v)
	}
	resp, err := s.Media(ctx, message.Media(media))
	return b.sentMessage(ctx, p, resp, err)
}

// SendLocation implements oas.Handler.
func (b *BotAPI) SendLocation(ctx context.Context, req oas.SendLocation) (oas.ResultMessage, error) {
	s, p, err := b.prepareSend(
		ctx,
		sendOpts{
			To:                       req.ChatID,
			DisableNotification:      req.DisableNotification,
			ProtectContent:           req.ProtectContent,
			ReplyToMessageID:         req.ReplyToMessageID,
			AllowSendingWithoutReply: req.AllowSendingWithoutReply,
			ReplyMarkup:              req.ReplyMarkup,
		},
	)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "prepare send")
	}

	point := &tg.InputGeoPoint{
		Lat:            req.Latitude,
		Long:           req.Longitude,
		AccuracyRadius: 0,
	}
	if v, ok := req.HorizontalAccuracy.Get(); ok {
		point.SetAccuracyRadius(int(v))
	}
	var media tg.InputMediaClass
	if livePeriod, ok := req.LivePeriod.Get(); ok {
		live := &tg.InputMediaGeoLive{
			Stopped:                     false,
			GeoPoint:                    point,
			Heading:                     0,
			Period:                      0,
			ProximityNotificationRadius: 0,
		}
		live.SetPeriod(livePeriod)
		if v, ok := req.Heading.Get(); ok {
			live.SetHeading(v)
		}
		if v, ok := req.ProximityAlertRadius.Get(); ok {
			live.SetProximityNotificationRadius(v)
		}
		media = live
	} else {
		media = &tg.InputMediaGeoPoint{
			GeoPoint: point,
		}
	}

	resp, err := s.Media(ctx, message.Media(media))
	return b.sentMessage(ctx, p, resp, err)
}

// SendVenue implements oas.Handler.
func (b *BotAPI) SendVenue(ctx context.Context, req oas.SendVenue) (oas.ResultMessage, error) {
	s, p, err := b.prepareSend(
		ctx,
		sendOpts{
			To:                       req.ChatID,
			DisableNotification:      req.DisableNotification,
			ProtectContent:           req.ProtectContent,
			ReplyToMessageID:         req.ReplyToMessageID,
			AllowSendingWithoutReply: req.AllowSendingWithoutReply,
			ReplyMarkup:              req.ReplyMarkup,
		},
	)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "prepare send")
	}

	point := &tg.InputGeoPoint{
		Lat:            req.Latitude,
		Long:           req.Longitude,
		AccuracyRadius: 0,
	}
	media := &tg.InputMediaVenue{
		GeoPoint: point,
		Title:    req.Title,
		Address:  req.Address,
	}
	if id, typ := req.FoursquareID.Value, req.FoursquareType.Value; id != "" || typ != "" {
		media.Provider = "foursquare"
		media.VenueID = id
		media.VenueType = typ
	}
	if id, typ := req.GooglePlaceID.Value, req.GooglePlaceType.Value; id != "" || typ != "" {
		media.Provider = "gplaces"
		media.VenueID = id
		media.VenueType = typ
	}

	resp, err := s.Media(ctx, message.Media(media))
	return b.sentMessage(ctx, p, resp, err)
}

// SendPoll implements oas.Handler.
func (b *BotAPI) SendPoll(ctx context.Context, req oas.SendPoll) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}
