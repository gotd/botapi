package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// LiveLocationOption configures a live-location edit.
type LiveLocationOption func(*liveLocationConfig)

type liveLocationConfig struct {
	heading              int
	proximityAlertRadius int
	accuracy             float64
	markup               ReplyMarkup
}

// WithHeading sets the direction in which the user is moving, 1-360 degrees.
func WithHeading(degrees int) LiveLocationOption {
	return func(c *liveLocationConfig) { c.heading = degrees }
}

// WithProximityAlertRadius sets the maximum distance, in meters, for proximity
// alerts about another chat member.
func WithProximityAlertRadius(meters int) LiveLocationOption {
	return func(c *liveLocationConfig) { c.proximityAlertRadius = meters }
}

// WithHorizontalAccuracy sets the radius of uncertainty for the location, in
// meters (0-1500).
func WithHorizontalAccuracy(meters float64) LiveLocationOption {
	return func(c *liveLocationConfig) { c.accuracy = meters }
}

// WithLiveLocationMarkup attaches or replaces the inline keyboard on the edited
// message.
func WithLiveLocationMarkup(markup ReplyMarkup) LiveLocationOption {
	return func(c *liveLocationConfig) { c.markup = markup }
}

// EditMessageLiveLocation edits a live location message, moving the point to the
// given coordinates.
func (b *Bot) EditMessageLiveLocation(
	ctx context.Context, chat ChatID, messageID int, latitude, longitude float64, opts ...LiveLocationOption,
) (*Message, error) {
	var cfg liveLocationConfig
	for _, o := range opts {
		o(&cfg)
	}

	media := &tg.InputMediaGeoLive{
		GeoPoint: &tg.InputGeoPoint{
			Lat:  latitude,
			Long: longitude,
		},
	}
	if cfg.accuracy > 0 {
		media.GeoPoint.(*tg.InputGeoPoint).AccuracyRadius = int(cfg.accuracy)
	}
	if cfg.heading > 0 {
		media.Heading = cfg.heading
	}
	if cfg.proximityAlertRadius > 0 {
		media.ProximityNotificationRadius = cfg.proximityAlertRadius
	}
	return b.editLiveLocation(ctx, chat, messageID, media, cfg.markup)
}

// StopMessageLiveLocation stops updating a live location message before its
// live period expires.
func (b *Bot) StopMessageLiveLocation(ctx context.Context, chat ChatID, messageID int, markup ReplyMarkup) (*Message, error) {
	media := &tg.InputMediaGeoLive{
		Stopped:  true,
		GeoPoint: &tg.InputGeoPointEmpty{},
	}
	return b.editLiveLocation(ctx, chat, messageID, media, markup)
}

// editLiveLocation issues the editMessage with a geo-live media payload.
func (b *Bot) editLiveLocation(
	ctx context.Context, chat ChatID, messageID int, media tg.InputMediaClass, markup ReplyMarkup,
) (*Message, error) {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	req := &tg.MessagesEditMessageRequest{Peer: peer, ID: messageID}
	req.SetMedia(media)
	if markup != nil {
		mkp, err := replyMarkupToTg(markup)
		if err != nil {
			return nil, err
		}
		req.SetReplyMarkup(mkp)
	}

	resp, err := b.raw.MessagesEditMessage(ctx, req)
	return b.sentMessage(ctx, peer, resp, err)
}
