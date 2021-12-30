package botapi

import (
	"context"
	"strconv"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

var maskCoordsNames = []string{"forehead", "eyes", "mouth", "chin"}

func (b *BotAPI) convertToBotAPIEntities(
	ctx context.Context,
	entities []tg.MessageEntityClass,
) (r []oas.MessageEntity) {
	for _, entity := range entities {
		e := oas.MessageEntity{
			Offset: entity.GetOffset(),
			Length: entity.GetLength(),
		}

		switch entity := entity.(type) {
		case *tg.MessageEntityMention:
			e.Type = oas.MessageEntityTypeMention
		case *tg.MessageEntityHashtag:
			e.Type = oas.MessageEntityTypeHashtag
		case *tg.MessageEntityBotCommand:
			e.Type = oas.MessageEntityTypeBotCommand
		case *tg.MessageEntityURL:
			e.Type = oas.MessageEntityTypeURL
		case *tg.MessageEntityEmail:
			e.Type = oas.MessageEntityTypeEmail
		case *tg.MessageEntityBold:
			e.Type = oas.MessageEntityTypeBold
		case *tg.MessageEntityItalic:
			e.Type = oas.MessageEntityTypeItalic
		case *tg.MessageEntityCode:
			e.Type = oas.MessageEntityTypeCode
		case *tg.MessageEntityPre:
			e.Type = oas.MessageEntityTypePre
			e.Language.SetTo(entity.Language)
		case *tg.MessageEntityTextURL:
			e.Type = oas.MessageEntityTypeTextLink
			e.URL.SetTo(entity.URL)
		case *tg.MessageEntityMentionName:
			e.Type = oas.MessageEntityTypeTextMention
			user, err := b.resolveUserID(ctx, entity.UserID)
			if err == nil {
				e.User.SetTo(convertToBotAPIUser(user))
				b.logger.Warn("Resolve user", zap.Int64("user_id", entity.UserID))
			}
		case *tg.MessageEntityPhone:
			e.Type = oas.MessageEntityTypePhoneNumber
		case *tg.MessageEntityCashtag:
			e.Type = oas.MessageEntityTypeCashtag
		case *tg.MessageEntityUnderline:
			e.Type = oas.MessageEntityTypeUnderline
		case *tg.MessageEntityStrike:
			e.Type = oas.MessageEntityTypeStrikethrough
		}
		r = append(r, e)
	}

	return r
}

func (b *BotAPI) convertToBotAPIPhotoSizes(p tg.PhotoClass) (r []oas.PhotoSize) {
	photo, ok := p.AsNotEmpty()
	if !ok {
		return nil
	}

	type sizedPhoto interface {
		GetW() int
		GetH() int
		GetType() string
	}
	for _, sz := range photo.Sizes {
		size, ok := sz.(sizedPhoto)
		if !ok {
			continue
		}

		t := size.GetType()
		if len(t) < 1 {
			continue
		}

		// TODO(tdakkota): compute size if downloaded
		var fileSize oas.OptInt
		switch size := size.(type) {
		case *tg.PhotoSize:
			fileSize.SetTo(size.Size)
		case *tg.PhotoCachedSize:
			fileSize.SetTo(len(size.Bytes))
		}

		fileID, fileUniqueID := b.encodeFileID(fileid.FromPhoto(photo, rune(t[0])))
		r = append(r, oas.PhotoSize{
			FileID:       fileID,
			FileUniqueID: fileUniqueID,
			Width:        size.GetW(),
			Height:       size.GetH(),
			FileSize:     fileSize,
		})
	}

	return r
}

func convertToBotAPILocation(p *tg.GeoPoint) (r oas.Location) {
	r = oas.Location{
		Longitude: p.Long,
		Latitude:  p.Lat,
	}
	if v, ok := p.GetAccuracyRadius(); ok {
		r.HorizontalAccuracy.SetTo(float64(v))
	}
	return r
}

func (b *BotAPI) setDocumentAttachment(ctx context.Context, d *tg.Document, r *oas.Message) error {
	f := fileid.FromDocument(d)
	fileID, fileUniqueID := b.encodeFileID(f)

	var (
		mimeType = oas.OptString{
			Value: d.MimeType,
			Set:   d.MimeType != "",
		}
		fileName oas.OptString
		// TODO(tdakkota): get thumb
		thumb oas.OptPhotoSize

		width    int
		height   int
		duration int

		animated bool
	)
	for _, attr := range d.Attributes {
		switch attr := attr.(type) {
		case *tg.DocumentAttributeFilename:
			fileName.SetTo(attr.FileName)
		case *tg.DocumentAttributeVideo:
			width = attr.W
			height = attr.H
			duration = attr.Duration
		case *tg.DocumentAttributeImageSize:
			width = attr.W
			height = attr.H
		case *tg.DocumentAttributeAnimated:
			animated = true
		}
	}

	for _, attr := range d.Attributes {
		switch attr := attr.(type) {
		case *tg.DocumentAttributeImageSize:
		case *tg.DocumentAttributeAnimated:
			r.Animation.SetTo(oas.Animation{
				FileID:       fileID,
				FileUniqueID: fileUniqueID,
				Width:        width,
				Height:       height,
				Duration:     duration,
				Thumb:        thumb,
				FileName:     fileName,
				MimeType:     mimeType,
				FileSize:     oas.NewOptInt(d.Size),
			})
		case *tg.DocumentAttributeSticker:
			var maskPosition oas.OptMaskPosition
			if coords, ok := attr.GetMaskCoords(); ok &&
				coords.N >= 0 && coords.N < len(maskCoordsNames) {
				maskPosition.SetTo(oas.MaskPosition{
					Point:  maskCoordsNames[coords.N],
					XShift: coords.X,
					YShift: coords.Y,
					Scale:  coords.Zoom,
				})
			}
			result, err := b.getStickerSet(ctx, attr.Stickerset)
			if err != nil {
				return errors.Wrap(err, "get sticker_set")
			}
			r.Sticker.SetTo(oas.Sticker{
				FileID:       fileID,
				FileUniqueID: fileUniqueID,
				Width:        width,
				Height:       height,
				IsAnimated:   result.Set.Animated,
				Thumb:        thumb,
				Emoji:        oas.NewOptString(attr.Alt),
				SetName:      oas.NewOptString(result.Set.ShortName),
				MaskPosition: maskPosition,
				FileSize:     oas.NewOptInt(d.Size),
			})
		case *tg.DocumentAttributeVideo:
			if animated {
				break
			}

			if attr.RoundMessage {
				r.VideoNote.SetTo(oas.VideoNote{
					FileID:       fileID,
					FileUniqueID: fileUniqueID,
					Length:       width,
					Duration:     duration,
					Thumb:        thumb,
					FileSize:     oas.NewOptInt(d.Size),
				})
			} else {
				r.Video.SetTo(oas.Video{
					FileID:       fileID,
					FileUniqueID: fileUniqueID,
					Width:        width,
					Height:       height,
					Duration:     duration,
					Thumb:        thumb,
					FileName:     fileName,
					MimeType:     mimeType,
					FileSize:     oas.NewOptInt(d.Size),
				})
			}
		case *tg.DocumentAttributeAudio:
			if attr.Voice {
				r.Voice.SetTo(oas.Voice{
					FileID:       fileID,
					FileUniqueID: fileUniqueID,
					Duration:     attr.Duration,
					MimeType:     mimeType,
					FileSize:     oas.NewOptInt(d.Size),
				})
			} else {
				r.Audio.SetTo(oas.Audio{
					FileID:       fileID,
					FileUniqueID: fileUniqueID,
					Duration:     attr.Duration,
					Performer:    optString(attr.GetPerformer),
					Title:        optString(attr.GetTitle),
					FileName:     fileName,
					MimeType:     mimeType,
					FileSize:     oas.NewOptInt(d.Size),
					Thumb:        thumb,
				})
			}
		case *tg.DocumentAttributeHasStickers:
		}
	}

	if !r.Sticker.Set &&
		!r.VideoNote.Set && !r.Video.Set &&
		!r.Voice.Set && !r.Audio.Set {
		r.Document.SetTo(oas.Document{
			FileID:       fileID,
			FileUniqueID: fileUniqueID,
			Thumb:        thumb,
			FileName:     fileName,
			MimeType:     mimeType,
			FileSize:     oas.NewOptInt(d.Size),
		})
	}

	return nil
}

func (b *BotAPI) convertMessageMedia(ctx context.Context, media tg.MessageMediaClass, r *oas.Message) error {
	switch media := media.(type) {
	case *tg.MessageMediaPhoto:
		r.Photo = b.convertToBotAPIPhotoSizes(media.Photo)
	case *tg.MessageMediaGeo:
		p, ok := media.Geo.AsNotEmpty()
		if !ok {
			break
		}
		r.Location.SetTo(convertToBotAPILocation(p))
	case *tg.MessageMediaContact:
		r.Contact.SetTo(oas.Contact{
			PhoneNumber: media.PhoneNumber,
			FirstName:   media.FirstName,
			LastName: oas.OptString{
				Value: media.LastName,
				Set:   media.LastName != "",
			},
			UserID: oas.OptInt64{
				Value: media.UserID,
				Set:   media.UserID != 0,
			},
			Vcard: oas.OptString{
				Value: media.Vcard,
				Set:   media.Vcard != "",
			},
		})
	case *tg.MessageMediaDocument:
		d, ok := media.Document.AsNotEmpty()
		if !ok {
			break
		}
		if err := b.setDocumentAttachment(ctx, d, r); err != nil {
			return errors.Wrap(err, "get document")
		}
	case *tg.MessageMediaWebPage:
		// Bots do not receive web page attachments.
	case *tg.MessageMediaVenue:
		p, ok := media.Geo.AsNotEmpty()
		if !ok {
			break
		}
		location := convertToBotAPILocation(p)
		resultVenue := oas.Venue{
			Location:        location,
			Title:           media.Title,
			Address:         media.Address,
			FoursquareID:    oas.OptString{},
			FoursquareType:  oas.OptString{},
			GooglePlaceID:   oas.OptString{},
			GooglePlaceType: oas.OptString{},
		}
		switch media.Provider {
		case "foursquare":
			resultVenue.FoursquareID.SetTo(media.VenueID)
			resultVenue.FoursquareType.SetTo(media.VenueType)
		case "gplaces":
			resultVenue.GooglePlaceID.SetTo(media.VenueID)
			resultVenue.GooglePlaceType.SetTo(media.VenueType)
		}
		r.Venue.SetTo(resultVenue)
		// Set for backward compatibility.
		r.Location.SetTo(location)
	case *tg.MessageMediaGame:
		game := media.Game

		r.Game.SetTo(oas.Game{
			Title:        game.Title,
			Description:  game.Description,
			Photo:        b.convertToBotAPIPhotoSizes(game.Photo),
			Text:         r.Text,
			TextEntities: r.Entities,
			Animation:    oas.OptAnimation{},
		})
	case *tg.MessageMediaInvoice:

	case *tg.MessageMediaGeoLive:
		p, ok := media.Geo.AsNotEmpty()
		if !ok {
			break
		}
		location := convertToBotAPILocation(p)
		location.Heading = optInt(media.GetHeading)
		location.LivePeriod.SetTo(media.Period)
		location.ProximityAlertRadius = optInt(media.GetProximityNotificationRadius)
		r.Location.SetTo(location)
	case *tg.MessageMediaPoll:
		var (
			poll    = media.Poll
			results = media.Results

			typ = oas.PollTypeRegular
		)
		if a, r := len(poll.Answers), len(results.Results); a != r {
			b.logger.Warn("Got poll where len(answers) != len(results)",
				zap.Int("answers", a),
				zap.Int("results", r),
			)
			break
		}

		if poll.Quiz {
			typ = oas.PollTypeQuiz
		}
		resultPoll := oas.Poll{
			ID:                    strconv.FormatInt(poll.ID, 10),
			Question:              poll.Question,
			Options:               nil,
			TotalVoterCount:       results.TotalVoters,
			IsClosed:              poll.Closed,
			IsAnonymous:           !poll.PublicVoters,
			Type:                  typ,
			AllowsMultipleAnswers: poll.MultipleChoice,
			CorrectOptionID:       oas.OptInt{},
			Explanation:           optString(results.GetSolution),
			ExplanationEntities:   nil,
			OpenPeriod:            optInt(poll.GetClosePeriod),
			CloseDate:             optInt(poll.GetCloseDate),
		}

		if e := results.SolutionEntities; len(e) > 0 {
			resultPoll.ExplanationEntities = b.convertToBotAPIEntities(ctx, e)
		}

		// SAFETY: length equality checked above.
		for i, result := range results.Results {
			if result.Correct {
				resultPoll.CorrectOptionID.SetTo(i)
			}
			resultPoll.Options = append(resultPoll.Options, oas.PollOption{
				Text:       poll.Answers[i].Text,
				VoterCount: result.Voters,
			})
		}

		r.Poll.SetTo(resultPoll)
	case *tg.MessageMediaDice:
		r.Dice.SetTo(oas.Dice{
			Emoji: media.Emoticon,
			Value: media.Value,
		})
	}

	return nil
}

func (b *BotAPI) convertPlainMessage(ctx context.Context, m *tg.Message) (r oas.Message, _ error) {
	getFrom := func(fromID tg.PeerClass, user *oas.OptUser, chat *oas.OptChat) error {
		switch fromID := fromID.(type) {
		case *tg.PeerUser:
			u, err := b.resolveUserID(ctx, fromID.UserID)
			if err != nil {
				return errors.Wrap(err, "get user")
			}
			user.SetTo(convertToBotAPIUser(u))
		case *tg.PeerChat, *tg.PeerChannel:
			ch, err := b.getChatByPeer(ctx, fromID)
			if err != nil {
				return errors.Wrap(err, "get chat")
			}
			chat.SetTo(ch)
		}
		return nil
	}

	ch, err := b.getChatByPeer(ctx, m.PeerID)
	if err != nil {
		return oas.Message{}, errors.Wrap(err, "get chat")
	}

	r = oas.Message{
		MessageID:           m.ID,
		Date:                m.Date,
		Chat:                ch,
		EditDate:            optInt(m.GetEditDate),
		HasProtectedContent: ch.HasProtectedContent,
		// TODO(tdakkota): generate media album ids
		MediaGroupID:    oas.OptString{},
		AuthorSignature: optString(m.GetPostAuthor),
	}
	if m.Out {
		self, err := b.peers.Self(ctx)
		if err == nil {
			r.From.SetTo(convertToBotAPIUser(self))
		}
	} else if fromID, ok := m.GetFromID(); ok {
		// FIXME(tdakkota): set service IDs.
		//
		// See https://github.com/tdlib/telegram-bot-api/blob/90f52477814a2d8a08c9ffb1d780fd179815d715/telegram-bot-api/Client.cpp#L9602
		if err := getFrom(fromID, &r.From, &r.SenderChat); err != nil {
			return oas.Message{}, errors.Wrap(err, "get from")
		}
	}

	// See https://github.com/tdlib/telegram-bot-api/blob/90f52477814a2d8a08c9ffb1d780fd179815d715/telegram-bot-api/Client.cpp#L9585-L9587
	if h, ok := m.GetFwdFrom(); ok {
		var isAutomaticForward bool
		if fromID, ok := m.GetFromID(); ok {
			if err := getFrom(fromID, &r.ForwardFrom, &r.ForwardFromChat); err != nil {
				return oas.Message{}, errors.Wrap(err, "get forward_from")
			}

			if from, to := r.ForwardFromChat.Value, r.Chat; r.ForwardFromChat.Set {
				_, isChannelPost := h.GetChannelPost()
				isAutomaticForward = isChannelPost &&
					from.ID != to.ID &&
					to.Type == oas.ChatTypeSupergroup &&
					from.Type == oas.ChatTypeChannel
				if isAutomaticForward {
					r.IsAutomaticForward.SetTo(true)
				}
			}
		}
		r.ForwardFromMessageID = optInt(h.GetChannelPost)
		r.ForwardSignature = optString(h.GetPostAuthor)
		r.ForwardSenderName = optString(h.GetFromName)
		r.ForwardDate = oas.NewOptInt(h.Date)
	}

	if reply, ok := m.GetReplyTo(); ok {
		// TODO(tdakkota): implement reply to resolve.
		r.ReplyToMessage = &oas.Message{
			MessageID: reply.ReplyToMsgID,
		}
	}

	if botID, ok := m.GetViaBotID(); ok {
		u, err := b.resolveUserID(ctx, botID)
		if err != nil {
			return oas.Message{}, errors.Wrap(err, "get via_bot")
		}
		r.ViaBot.SetTo(convertToBotAPIUser(u))
	}

	if text := m.Message; text != "" {
		r.Text.SetTo(text)
	}
	if len(m.Entities) > 0 {
		r.Entities = b.convertToBotAPIEntities(ctx, m.Entities)
	}

	if err := b.convertMessageMedia(ctx, m.Media, &r); err != nil {
		return oas.Message{}, errors.Wrap(err, "get media")
	}

	if mkp, ok := m.ReplyMarkup.(*tg.ReplyInlineMarkup); ok {
		r.ReplyMarkup.SetTo(convertToBotAPIInlineReplyMarkup(mkp))
	}

	return r, nil
}
