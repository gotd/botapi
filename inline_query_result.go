package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// mimeImageJPEG is the default thumbnail MIME type for inline results.
const mimeImageJPEG = "image/jpeg"

// InlineQueryResult is a sealed union describing one result of an inline query.
//
// Concrete variants cover articles, fresh-by-URL media (photo/gif/mpeg4
// gif/video), cached-by-file_id media (photo/gif/mpeg4 gif/sticker/document/
// video/audio/voice) and contact/location/venue results.
type InlineQueryResult interface {
	isInlineQueryResult()
	toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error)
}

// captioned is the caption fields shared by media results.
type captioned struct {
	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
}

// InlineQueryResultArticle is a link to an article or web page.
type InlineQueryResultArticle struct {
	ID                  string              `json:"id"`
	Title               string              `json:"title"`
	InputMessageContent InputMessageContent `json:"input_message_content"`
	ReplyMarkup         *InlineKeyboardMarkup
	URL                 string `json:"url,omitempty"`
	Description         string `json:"description,omitempty"`
	ThumbnailURL        string `json:"thumbnail_url,omitempty"`
}

// InlineQueryResultPhoto is a link to a photo, sent by URL.
type InlineQueryResultPhoto struct {
	ID           string `json:"id"`
	PhotoURL     string `json:"photo_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	PhotoWidth   int    `json:"photo_width,omitempty"`
	PhotoHeight  int    `json:"photo_height,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	captioned
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultGif is a link to an animated GIF, sent by URL.
type InlineQueryResultGif struct {
	ID           string `json:"id"`
	GifURL       string `json:"gif_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	GifWidth     int    `json:"gif_width,omitempty"`
	GifHeight    int    `json:"gif_height,omitempty"`
	Title        string `json:"title,omitempty"`
	captioned
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultMpeg4Gif is a link to an MPEG-4 animation (video without
// sound), sent by URL.
type InlineQueryResultMpeg4Gif struct {
	ID           string `json:"id"`
	Mpeg4URL     string `json:"mpeg4_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Title        string `json:"title,omitempty"`
	captioned
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultCachedPhoto is a link to a photo stored on Telegram's
// servers, referenced by file_id.
type InlineQueryResultCachedPhoto struct {
	ID          string `json:"id"`
	PhotoFileID string `json:"photo_file_id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	captioned
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultCachedSticker is a link to a sticker stored on Telegram's
// servers, referenced by file_id.
type InlineQueryResultCachedSticker struct {
	ID                  string `json:"id"`
	StickerFileID       string `json:"sticker_file_id"`
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultCachedDocument is a link to a file stored on Telegram's
// servers, referenced by file_id.
type InlineQueryResultCachedDocument struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	DocumentFileID string `json:"document_file_id"`
	Description    string `json:"description,omitempty"`
	captioned
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultCachedGif is a link to an animated GIF stored on Telegram's
// servers, referenced by file_id.
type InlineQueryResultCachedGif struct {
	ID        string `json:"id"`
	GifFileID string `json:"gif_file_id"`
	Title     string `json:"title,omitempty"`
	captioned
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultCachedVideo is a link to a video stored on Telegram's
// servers, referenced by file_id.
type InlineQueryResultCachedVideo struct {
	ID          string `json:"id"`
	VideoFileID string `json:"video_file_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	captioned
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultCachedVoice is a link to a voice message stored on Telegram's
// servers, referenced by file_id.
type InlineQueryResultCachedVoice struct {
	ID          string `json:"id"`
	VoiceFileID string `json:"voice_file_id"`
	Title       string `json:"title"`
	captioned
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultCachedAudio is a link to an audio file stored on Telegram's
// servers, referenced by file_id.
type InlineQueryResultCachedAudio struct {
	ID          string `json:"id"`
	AudioFileID string `json:"audio_file_id"`
	captioned
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultContact is a contact with a phone number.
type InlineQueryResultContact struct {
	ID                  string `json:"id"`
	PhoneNumber         string `json:"phone_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name,omitempty"`
	Vcard               string `json:"vcard,omitempty"`
	ThumbnailURL        string `json:"thumbnail_url,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

// InlineQueryResultLocation is a location on a map.
type InlineQueryResultLocation struct {
	ID                   string  `json:"id"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Title                string  `json:"title"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
	ThumbnailURL         string  `json:"thumbnail_url,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup
	InputMessageContent  InputMessageContent
}

// InlineQueryResultVenue is a venue.
type InlineQueryResultVenue struct {
	ID                  string  `json:"id"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	Title               string  `json:"title"`
	Address             string  `json:"address"`
	FoursquareID        string  `json:"foursquare_id,omitempty"`
	FoursquareType      string  `json:"foursquare_type,omitempty"`
	GooglePlaceID       string  `json:"google_place_id,omitempty"`
	GooglePlaceType     string  `json:"google_place_type,omitempty"`
	ThumbnailURL        string  `json:"thumbnail_url,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup
	InputMessageContent InputMessageContent
}

func (*InlineQueryResultContact) isInlineQueryResult()        {}
func (*InlineQueryResultLocation) isInlineQueryResult()       {}
func (*InlineQueryResultVenue) isInlineQueryResult()          {}
func (*InlineQueryResultArticle) isInlineQueryResult()        {}
func (*InlineQueryResultPhoto) isInlineQueryResult()          {}
func (*InlineQueryResultGif) isInlineQueryResult()            {}
func (*InlineQueryResultMpeg4Gif) isInlineQueryResult()       {}
func (*InlineQueryResultCachedPhoto) isInlineQueryResult()    {}
func (*InlineQueryResultCachedSticker) isInlineQueryResult()  {}
func (*InlineQueryResultCachedDocument) isInlineQueryResult() {}
func (*InlineQueryResultCachedGif) isInlineQueryResult()      {}
func (*InlineQueryResultCachedVideo) isInlineQueryResult()    {}
func (*InlineQueryResultCachedVoice) isInlineQueryResult()    {}
func (*InlineQueryResultCachedAudio) isInlineQueryResult()    {}

// resultMarkup resolves an optional inline keyboard to its MTProto form.
func resultMarkup(markup *InlineKeyboardMarkup) (tg.ReplyMarkupClass, error) {
	if markup == nil {
		return nil, nil
	}

	return replyMarkupToTg(markup)
}

// sendMessage builds the inline result's message: the explicit message content
// if provided, otherwise a media-auto message carrying the caption.
func (b *Bot) sendMessage(
	ctx context.Context,
	content InputMessageContent,
	capt captioned,
	markup tg.ReplyMarkupClass,
) (tg.InputBotInlineMessageClass, error) {
	if content != nil {
		return content.toTg(ctx, b, markup)
	}

	msg := &tg.InputBotInlineMessageMediaAuto{Message: capt.Caption}
	if len(capt.CaptionEntities) > 0 {
		msg.Entities = entitiesToTg(capt.CaptionEntities)
	} else if capt.ParseMode != ParseModeNone && capt.Caption != "" {
		text, entities, err := b.styledMessage(ctx, capt.Caption, capt.ParseMode)
		if err != nil {
			return nil, err
		}

		msg.Message = text
		msg.Entities = entities
	}

	if markup != nil {
		msg.SetReplyMarkup(markup)
	}

	return msg, nil
}

func (r *InlineQueryResultArticle) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	if r.InputMessageContent == nil {
		return nil, &Error{Code: 400, Description: "Bad Request: article result requires input_message_content"}
	}

	markup, err := resultMarkup(r.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	send, err := r.InputMessageContent.toTg(ctx, b, markup)
	if err != nil {
		return nil, err
	}

	res := &tg.InputBotInlineResult{
		ID:          r.ID,
		Type:        "article",
		Title:       r.Title,
		Description: r.Description,
		SendMessage: send,
	}
	if r.URL != "" {
		res.SetURL(r.URL)
	}

	if r.ThumbnailURL != "" {
		res.SetThumb(tg.InputWebDocument{URL: r.ThumbnailURL, MimeType: mimeImageJPEG})
	}

	return res, nil
}

// webResult builds a fresh-by-URL inline result with a web-document body.
func (b *Bot) webResult(
	ctx context.Context,
	id, typ, title, description, contentURL, mimeType, thumbURL string,
	width, height int,
	content InputMessageContent,
	capt captioned,
	markup *InlineKeyboardMarkup,
) (tg.InputBotInlineResultClass, error) {
	mkp, err := resultMarkup(markup)
	if err != nil {
		return nil, err
	}

	send, err := b.sendMessage(ctx, content, capt, mkp)
	if err != nil {
		return nil, err
	}

	res := &tg.InputBotInlineResult{
		ID:          id,
		Type:        typ,
		Title:       title,
		Description: description,
		SendMessage: send,
	}

	if contentURL != "" {
		doc := tg.InputWebDocument{URL: contentURL, MimeType: mimeType}
		if width > 0 || height > 0 {
			doc.Attributes = []tg.DocumentAttributeClass{
				&tg.DocumentAttributeImageSize{W: width, H: height},
			}
		}

		res.SetContent(doc)
	}

	if thumbURL != "" {
		res.SetThumb(tg.InputWebDocument{URL: thumbURL, MimeType: mimeImageJPEG})
	}

	return res, nil
}

func (r *InlineQueryResultPhoto) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	return b.webResult(ctx, r.ID, "photo", r.Title, r.Description,
		r.PhotoURL, mimeImageJPEG, r.ThumbnailURL, r.PhotoWidth, r.PhotoHeight,
		r.InputMessageContent, r.captioned, r.ReplyMarkup)
}

func (r *InlineQueryResultGif) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	return b.webResult(ctx, r.ID, "gif", r.Title, "",
		r.GifURL, "image/gif", r.ThumbnailURL, r.GifWidth, r.GifHeight,
		r.InputMessageContent, r.captioned, r.ReplyMarkup)
}

func (r *InlineQueryResultMpeg4Gif) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	return b.webResult(ctx, r.ID, "gif", r.Title, "",
		r.Mpeg4URL, "video/mp4", r.ThumbnailURL, 0, 0,
		r.InputMessageContent, r.captioned, r.ReplyMarkup)
}

// cachedDocument builds a cached (file_id) inline result backed by a document.
func (b *Bot) cachedDocument(
	ctx context.Context,
	id, typ, title, description, fileID string,
	content InputMessageContent,
	capt captioned,
	markup *InlineKeyboardMarkup,
) (tg.InputBotInlineResultClass, error) {
	doc, err := inputDocumentFromFileID(fileID)
	if err != nil {
		return nil, err
	}

	mkp, err := resultMarkup(markup)
	if err != nil {
		return nil, err
	}

	send, err := b.sendMessage(ctx, content, capt, mkp)
	if err != nil {
		return nil, err
	}

	return &tg.InputBotInlineResultDocument{
		ID:          id,
		Type:        typ,
		Title:       title,
		Description: description,
		Document:    doc,
		SendMessage: send,
	}, nil
}

func (r *InlineQueryResultCachedPhoto) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	photo, err := inputPhotoFromFileID(r.PhotoFileID)
	if err != nil {
		return nil, err
	}

	mkp, err := resultMarkup(r.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	send, err := b.sendMessage(ctx, r.InputMessageContent, r.captioned, mkp)
	if err != nil {
		return nil, err
	}

	return &tg.InputBotInlineResultPhoto{
		ID:          r.ID,
		Type:        "photo",
		Photo:       photo,
		SendMessage: send,
	}, nil
}

func (r *InlineQueryResultCachedSticker) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	return b.cachedDocument(ctx, r.ID, "sticker", "", "", r.StickerFileID,
		r.InputMessageContent, captioned{}, r.ReplyMarkup)
}

func (r *InlineQueryResultCachedDocument) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	return b.cachedDocument(ctx, r.ID, "file", r.Title, r.Description, r.DocumentFileID,
		r.InputMessageContent, r.captioned, r.ReplyMarkup)
}

func (r *InlineQueryResultCachedGif) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	return b.cachedDocument(ctx, r.ID, "gif", r.Title, "", r.GifFileID,
		r.InputMessageContent, r.captioned, r.ReplyMarkup)
}

func (r *InlineQueryResultCachedVideo) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	return b.cachedDocument(ctx, r.ID, "video", r.Title, r.Description, r.VideoFileID,
		r.InputMessageContent, r.captioned, r.ReplyMarkup)
}

func (r *InlineQueryResultCachedVoice) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	return b.cachedDocument(ctx, r.ID, "voice", r.Title, "", r.VoiceFileID,
		r.InputMessageContent, r.captioned, r.ReplyMarkup)
}

func (r *InlineQueryResultCachedAudio) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	return b.cachedDocument(ctx, r.ID, "audio", "", "", r.AudioFileID,
		r.InputMessageContent, r.captioned, r.ReplyMarkup)
}

func (r *InlineQueryResultContact) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	mkp, err := resultMarkup(r.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	var send tg.InputBotInlineMessageClass

	if r.InputMessageContent != nil {
		if send, err = r.InputMessageContent.toTg(ctx, b, mkp); err != nil {
			return nil, err
		}
	} else {
		msg := &tg.InputBotInlineMessageMediaContact{
			PhoneNumber: r.PhoneNumber,
			FirstName:   r.FirstName,
			LastName:    r.LastName,
			Vcard:       r.Vcard,
		}
		if mkp != nil {
			msg.SetReplyMarkup(mkp)
		}

		send = msg
	}

	res := &tg.InputBotInlineResult{
		ID:          r.ID,
		Type:        "contact",
		Title:       r.FirstName,
		Description: r.PhoneNumber,
		SendMessage: send,
	}
	if r.ThumbnailURL != "" {
		res.SetThumb(tg.InputWebDocument{URL: r.ThumbnailURL, MimeType: mimeImageJPEG})
	}

	return res, nil
}

func (r *InlineQueryResultLocation) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	mkp, err := resultMarkup(r.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	var send tg.InputBotInlineMessageClass

	if r.InputMessageContent != nil {
		if send, err = r.InputMessageContent.toTg(ctx, b, mkp); err != nil {
			return nil, err
		}
	} else {
		msg := &tg.InputBotInlineMessageMediaGeo{
			GeoPoint: &tg.InputGeoPoint{
				Lat:            r.Latitude,
				Long:           r.Longitude,
				AccuracyRadius: int(r.HorizontalAccuracy),
			},
			Heading:                     r.Heading,
			Period:                      r.LivePeriod,
			ProximityNotificationRadius: r.ProximityAlertRadius,
		}
		if mkp != nil {
			msg.SetReplyMarkup(mkp)
		}

		send = msg
	}

	res := &tg.InputBotInlineResult{
		ID:          r.ID,
		Type:        "geo",
		Title:       r.Title,
		SendMessage: send,
	}
	if r.ThumbnailURL != "" {
		res.SetThumb(tg.InputWebDocument{URL: r.ThumbnailURL, MimeType: mimeImageJPEG})
	}

	return res, nil
}

func (r *InlineQueryResultVenue) toTg(ctx context.Context, b *Bot) (tg.InputBotInlineResultClass, error) {
	mkp, err := resultMarkup(r.ReplyMarkup)
	if err != nil {
		return nil, err
	}

	var send tg.InputBotInlineMessageClass

	if r.InputMessageContent != nil {
		if send, err = r.InputMessageContent.toTg(ctx, b, mkp); err != nil {
			return nil, err
		}
	} else {
		venue := &InputVenueMessageContent{
			Latitude:        r.Latitude,
			Longitude:       r.Longitude,
			Title:           r.Title,
			Address:         r.Address,
			FoursquareID:    r.FoursquareID,
			FoursquareType:  r.FoursquareType,
			GooglePlaceID:   r.GooglePlaceID,
			GooglePlaceType: r.GooglePlaceType,
		}
		if send, err = venue.toTg(ctx, b, mkp); err != nil {
			return nil, err
		}
	}

	res := &tg.InputBotInlineResult{
		ID:          r.ID,
		Type:        "venue",
		Title:       r.Title,
		Description: r.Address,
		SendMessage: send,
	}
	if r.ThumbnailURL != "" {
		res.SetThumb(tg.InputWebDocument{URL: r.ThumbnailURL, MimeType: mimeImageJPEG})
	}

	return res, nil
}
