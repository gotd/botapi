package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// ChatBoostSource is a sealed union describing how a chat boost was obtained.
//
// Concrete variants: ChatBoostSourcePremium, ChatBoostSourceGiftCode,
// ChatBoostSourceGiveaway.
type ChatBoostSource interface {
	isChatBoostSource()
}

// ChatBoostSourcePremium is a boost from a user subscribing to Telegram Premium
// or gifting a Premium subscription to another user.
type ChatBoostSourcePremium struct {
	Source string `json:"source"`
	User   User   `json:"user"`
}

// ChatBoostSourceGiftCode is a boost from a user using a gift code made by the
// chat (for example as part of a giveaway).
type ChatBoostSourceGiftCode struct {
	Source string `json:"source"`
	User   User   `json:"user"`
}

// ChatBoostSourceGiveaway is a boost from a giveaway prize created by the chat.
type ChatBoostSourceGiveaway struct {
	Source string `json:"source"`
	// GiveawayMessageID is the id of the message with the giveaway; it may be 0
	// if the message was deleted.
	GiveawayMessageID int `json:"giveaway_message_id"`
	// User is the user that won the prize, if any.
	User *User `json:"user,omitempty"`
	// IsUnclaimed reports whether the giveaway prize was not claimed.
	IsUnclaimed bool `json:"is_unclaimed,omitempty"`
}

func (ChatBoostSourcePremium) isChatBoostSource()  {}
func (ChatBoostSourceGiftCode) isChatBoostSource() {}
func (ChatBoostSourceGiveaway) isChatBoostSource() {}

// Chat boost source discriminators.
const (
	chatBoostSourcePremium  = "premium"
	chatBoostSourceGiftCode = "gift_code"
	chatBoostSourceGiveaway = "giveaway"
)

// ChatBoost describes a boost applied to a chat.
type ChatBoost struct {
	// BoostID is the unique identifier of the boost.
	BoostID string `json:"boost_id"`
	// AddDate is the Unix time when the chat was boosted.
	AddDate int `json:"add_date"`
	// ExpirationDate is the Unix time when the boost will automatically expire.
	ExpirationDate int `json:"expiration_date"`
	// Source is the source of the added boost.
	Source ChatBoostSource `json:"source"`
}

// chatBoostFromTg converts an MTProto boost into the Bot API ChatBoost. users
// resolves the booster by id.
func chatBoostFromTg(boost tg.Boost, users map[int64]*tg.User) ChatBoost {
	var booster User

	if u, ok := users[boost.UserID]; ok {
		booster = userFromTgUser(u)
	} else {
		booster = User{ID: boost.UserID}
	}

	var source ChatBoostSource

	switch {
	case boost.Giveaway:
		g := ChatBoostSourceGiveaway{
			Source:            chatBoostSourceGiveaway,
			GiveawayMessageID: boost.GiveawayMsgID,
			IsUnclaimed:       boost.Unclaimed,
		}
		if boost.UserID != 0 {
			source = withGiveawayUser(g, booster)
		} else {
			source = g
		}
	case boost.Gift:
		source = ChatBoostSourceGiftCode{Source: chatBoostSourceGiftCode, User: booster}
	default:
		source = ChatBoostSourcePremium{Source: chatBoostSourcePremium, User: booster}
	}

	return ChatBoost{
		BoostID:        boost.ID,
		AddDate:        boost.Date,
		ExpirationDate: boost.Expires,
		Source:         source,
	}
}

// withGiveawayUser attaches the winning user to a giveaway boost source.
func withGiveawayUser(g ChatBoostSourceGiveaway, user User) ChatBoostSourceGiveaway {
	u := user

	g.User = &u

	return g
}

// GetUserChatBoosts returns the list of boosts added to a chat by a user. The
// bot must be an administrator in the chat.
func (b *Bot) GetUserChatBoosts(ctx context.Context, chat ChatID, userID int64) ([]ChatBoost, error) {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.PremiumGetUserBoosts(ctx, &tg.PremiumGetUserBoostsRequest{
		Peer:   peer,
		UserID: user,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	users := usersByID(res.Users)
	out := make([]ChatBoost, 0, len(res.Boosts))

	for _, boost := range res.Boosts {
		out = append(out, chatBoostFromTg(boost, users))
	}

	return out, nil
}
