package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// GiftBackground describes the background palette of a gift.
type GiftBackground struct {
	// CenterColor is the center color of the background, in RGB format.
	CenterColor int `json:"center_color"`
	// EdgeColor is the edge color of the background, in RGB format.
	EdgeColor int `json:"edge_color"`
	// TextColor is the color to use for text on the background, in RGB format.
	TextColor int `json:"text_color"`
}

// Gift describes a gift that can be sent by the bot.
type Gift struct {
	// ID is the unique identifier of the gift.
	ID string `json:"id"`
	// Sticker is the sticker that represents the gift.
	Sticker Sticker `json:"sticker"`
	// StarCount is the number of Telegram Stars that must be paid to send the
	// sticker.
	StarCount int `json:"star_count"`
	// UpgradeStarCount is the number of Telegram Stars that must be paid to
	// upgrade the gift to a unique one.
	UpgradeStarCount int `json:"upgrade_star_count,omitempty"`
	// IsPremium reports whether the gift can be purchased only by Telegram
	// Premium subscribers.
	IsPremium bool `json:"is_premium,omitempty"`
	// HasColors reports whether the gift can be used to generate a message color
	// palette.
	HasColors bool `json:"has_colors,omitempty"`
	// TotalCount is the total number of the gifts of this type that can be sent;
	// only for limited gifts.
	TotalCount int `json:"total_count,omitempty"`
	// RemainingCount is the number of remaining gifts of this type that can be
	// sent; only for limited gifts.
	RemainingCount int `json:"remaining_count,omitempty"`
	// PersonalTotalCount is the total number of the gifts of this type that can
	// be owned by a single user; only for limited-per-user gifts.
	PersonalTotalCount int `json:"personal_total_count,omitempty"`
	// PersonalRemainingCount is the number of remaining gifts of this type that
	// can be owned by the current user; only for limited-per-user gifts.
	PersonalRemainingCount int `json:"personal_remaining_count,omitempty"`
	// UniqueGiftVariantCount is the number of possible unique variants the gift
	// can be upgraded to.
	UniqueGiftVariantCount int `json:"unique_gift_variant_count,omitempty"`
	// Background is the default background of the gift.
	Background *GiftBackground `json:"background,omitempty"`
	// PublisherChat is the chat that published the gift, if any.
	PublisherChat *Chat `json:"publisher_chat,omitempty"`
}

// GetAvailableGifts returns the list of gifts that can be sent by the bot.
func (b *Bot) GetAvailableGifts(ctx context.Context) ([]Gift, error) {
	res, err := b.raw.PaymentsGetStarGifts(ctx, 0)
	if err != nil {
		return nil, asAPIError(err)
	}

	gifts, ok := res.(*tg.PaymentsStarGifts)
	if !ok {
		// PaymentsStarGiftsNotModified only appears when a cache hash is sent,
		// which this method never does.
		return nil, &Error{Code: 500, Description: "Internal Server Error: unexpected gifts response"}
	}

	chats := chatsByID(gifts.Chats)

	out := make([]Gift, 0, len(gifts.Gifts))
	for _, g := range gifts.Gifts {
		// getStarGifts also returns upgraded (unique) collectibles; the Bot API
		// surface only describes regular, sendable gifts.
		raw, ok := g.(*tg.StarGift)
		if !ok {
			continue
		}

		out = append(out, giftFromTg(raw, chats))
	}

	return out, nil
}

// giftFromTg converts a raw MTProto star gift into the Bot API Gift. chats
// resolves the publisher peer.
func giftFromTg(g *tg.StarGift, chats map[int64]tg.ChatClass) Gift {
	gift := Gift{
		ID:        strconv.FormatInt(g.ID, 10),
		StarCount: int(g.Stars),
		IsPremium: g.RequirePremium,
		HasColors: g.PeerColorAvailable,
	}

	if doc, ok := g.Sticker.(*tg.Document); ok {
		// Gift stickers carry no enclosing set; they are plain regular stickers.
		gift.Sticker = stickerFromDocument(doc, "", StickerRegular)
	}

	if v, ok := g.GetUpgradeStars(); ok {
		gift.UpgradeStarCount = int(v)
	}

	if v, ok := g.GetAvailabilityTotal(); ok {
		gift.TotalCount = v
	}

	if v, ok := g.GetAvailabilityRemains(); ok {
		gift.RemainingCount = v
	}

	if v, ok := g.GetPerUserTotal(); ok {
		gift.PersonalTotalCount = v
	}

	if v, ok := g.GetPerUserRemains(); ok {
		gift.PersonalRemainingCount = v
	}

	if v, ok := g.GetUpgradeVariants(); ok {
		gift.UniqueGiftVariantCount = v
	}

	if bg, ok := g.GetBackground(); ok {
		gift.Background = &GiftBackground{
			CenterColor: bg.CenterColor,
			EdgeColor:   bg.EdgeColor,
			TextColor:   bg.TextColor,
		}
	}

	if peer, ok := g.GetReleasedBy(); ok {
		if chat, ok := chatFromPublisher(peer, chats); ok {
			gift.PublisherChat = &chat
		}
	}

	return gift
}

// chatFromPublisher resolves a gift publisher peer into a Bot API Chat. Only
// chat/channel publishers map onto publisher_chat.
func chatFromPublisher(p tg.PeerClass, chats map[int64]tg.ChatClass) (Chat, bool) {
	switch p.(type) {
	case *tg.PeerChat, *tg.PeerChannel:
		return chatFromRaw(p, chats), true
	default:
		return Chat{}, false
	}
}
