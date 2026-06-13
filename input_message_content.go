package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// InputMessageContent is a sealed union describing the content of a message to
// be sent as the result of an inline query.
//
// Concrete variants: *InputTextMessageContent, *InputLocationMessageContent,
// *InputVenueMessageContent, *InputContactMessageContent.
type InputMessageContent interface {
	isInputMessageContent()
	// toTg builds the MTProto inline message, attaching the (already resolved)
	// reply markup from the enclosing result.
	toTg(ctx context.Context, b *Bot, markup tg.ReplyMarkupClass) (tg.InputBotInlineMessageClass, error)
}

// InputTextMessageContent is the content of a text message.
type InputTextMessageContent struct {
	MessageText           string          `json:"message_text"`
	ParseMode             ParseMode       `json:"parse_mode,omitempty"`
	Entities              []MessageEntity `json:"entities,omitempty"`
	DisableWebPagePreview bool            `json:"disable_web_page_preview,omitempty"`
}

// InputLocationMessageContent is the content of a location message.
type InputLocationMessageContent struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
}

// InputVenueMessageContent is the content of a venue message.
type InputVenueMessageContent struct {
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Title           string  `json:"title"`
	Address         string  `json:"address"`
	FoursquareID    string  `json:"foursquare_id,omitempty"`
	FoursquareType  string  `json:"foursquare_type,omitempty"`
	GooglePlaceID   string  `json:"google_place_id,omitempty"`
	GooglePlaceType string  `json:"google_place_type,omitempty"`
}

// InputContactMessageContent is the content of a contact message.
type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	Vcard       string `json:"vcard,omitempty"`
}

func (*InputTextMessageContent) isInputMessageContent()     {}
func (*InputLocationMessageContent) isInputMessageContent() {}
func (*InputVenueMessageContent) isInputMessageContent()    {}
func (*InputContactMessageContent) isInputMessageContent()  {}

func (c *InputTextMessageContent) toTg(ctx context.Context, b *Bot, markup tg.ReplyMarkupClass) (tg.InputBotInlineMessageClass, error) {
	msg := &tg.InputBotInlineMessageText{
		NoWebpage: c.DisableWebPagePreview,
		Message:   c.MessageText,
	}
	if len(c.Entities) > 0 {
		msg.Entities = entitiesToTg(c.Entities)
	} else if c.ParseMode != ParseModeNone {
		text, entities, err := b.styledMessage(ctx, c.MessageText, c.ParseMode)
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

func (c *InputLocationMessageContent) toTg(_ context.Context, _ *Bot, markup tg.ReplyMarkupClass) (tg.InputBotInlineMessageClass, error) { //nolint:unparam,lll // signature satisfies the InputMessageContent interface
	msg := &tg.InputBotInlineMessageMediaGeo{
		GeoPoint: &tg.InputGeoPoint{
			Lat:            c.Latitude,
			Long:           c.Longitude,
			AccuracyRadius: int(c.HorizontalAccuracy),
		},
		Heading:                     c.Heading,
		Period:                      c.LivePeriod,
		ProximityNotificationRadius: c.ProximityAlertRadius,
	}
	if markup != nil {
		msg.SetReplyMarkup(markup)
	}
	return msg, nil
}

func (c *InputVenueMessageContent) toTg(_ context.Context, _ *Bot, markup tg.ReplyMarkupClass) (tg.InputBotInlineMessageClass, error) { //nolint:unparam,lll // signature satisfies the InputMessageContent interface
	msg := &tg.InputBotInlineMessageMediaVenue{
		GeoPoint:  &tg.InputGeoPoint{Lat: c.Latitude, Long: c.Longitude},
		Title:     c.Title,
		Address:   c.Address,
		Provider:  venueProvider(c),
		VenueID:   firstNonEmpty(c.FoursquareID, c.GooglePlaceID),
		VenueType: firstNonEmpty(c.FoursquareType, c.GooglePlaceType),
	}
	if markup != nil {
		msg.SetReplyMarkup(markup)
	}
	return msg, nil
}

func (c *InputContactMessageContent) toTg(_ context.Context, _ *Bot, markup tg.ReplyMarkupClass) (tg.InputBotInlineMessageClass, error) { //nolint:unparam,lll // signature satisfies the InputMessageContent interface
	msg := &tg.InputBotInlineMessageMediaContact{
		PhoneNumber: c.PhoneNumber,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Vcard:       c.Vcard,
	}
	if markup != nil {
		msg.SetReplyMarkup(markup)
	}
	return msg, nil
}

// venueProvider reports which venue provider the content uses.
func venueProvider(c *InputVenueMessageContent) string {
	if c.FoursquareID != "" || c.FoursquareType != "" {
		return "foursquare"
	}
	if c.GooglePlaceID != "" || c.GooglePlaceType != "" {
		return "gplaces"
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
