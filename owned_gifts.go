package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// giftListConfig holds the optional filters of the gift-listing methods.
type giftListConfig struct {
	offset               string
	limit                int
	sortByPrice          bool
	excludeUnsaved       bool
	excludeSaved         bool
	excludeUnlimited     bool
	excludeUnique        bool
	excludeUpgradable    bool
	excludeNonUpgradable bool
}

// GiftListOption customizes the gift-listing methods (GetUserGifts,
// GetChatGifts, GetBusinessAccountGifts).
type GiftListOption func(*giftListConfig)

// WithGiftListOffset sets the pagination offset, taken from a previous result's
// NextOffset.
func WithGiftListOffset(offset string) GiftListOption {
	return func(c *giftListConfig) { c.offset = offset }
}

// WithGiftListLimit caps the number of gifts returned.
func WithGiftListLimit(limit int) GiftListOption {
	return func(c *giftListConfig) { c.limit = limit }
}

// WithGiftListSortByPrice sorts the gifts by price instead of by reception date.
func WithGiftListSortByPrice() GiftListOption {
	return func(c *giftListConfig) { c.sortByPrice = true }
}

// WithGiftListExcludeUnsaved excludes gifts not displayed on the profile.
func WithGiftListExcludeUnsaved() GiftListOption {
	return func(c *giftListConfig) { c.excludeUnsaved = true }
}

// WithGiftListExcludeSaved excludes gifts displayed on the profile.
func WithGiftListExcludeSaved() GiftListOption {
	return func(c *giftListConfig) { c.excludeSaved = true }
}

// WithGiftListExcludeUnlimited excludes gifts that can be bought unlimited times.
func WithGiftListExcludeUnlimited() GiftListOption {
	return func(c *giftListConfig) { c.excludeUnlimited = true }
}

// WithGiftListExcludeUnique excludes unique (collectible) gifts.
func WithGiftListExcludeUnique() GiftListOption {
	return func(c *giftListConfig) { c.excludeUnique = true }
}

// WithGiftListExcludeUpgradable excludes limited gifts that can be upgraded.
func WithGiftListExcludeUpgradable() GiftListOption {
	return func(c *giftListConfig) { c.excludeUpgradable = true }
}

// WithGiftListExcludeNonUpgradable excludes limited gifts that cannot be
// upgraded.
func WithGiftListExcludeNonUpgradable() GiftListOption {
	return func(c *giftListConfig) { c.excludeNonUpgradable = true }
}

// request builds the MTProto request for the given peer from the filters.
func (c giftListConfig) request(peer tg.InputPeerClass) *tg.PaymentsGetSavedStarGiftsRequest {
	return &tg.PaymentsGetSavedStarGiftsRequest{
		Peer:                peer,
		Offset:              c.offset,
		Limit:               c.limit,
		SortByValue:         c.sortByPrice,
		ExcludeUnsaved:      c.excludeUnsaved,
		ExcludeSaved:        c.excludeSaved,
		ExcludeUnlimited:    c.excludeUnlimited,
		ExcludeUnique:       c.excludeUnique,
		ExcludeUpgradable:   c.excludeUpgradable,
		ExcludeUnupgradable: c.excludeNonUpgradable,
	}
}

// GetUserGifts returns the gifts owned and hosted by a user.
func (b *Bot) GetUserGifts(ctx context.Context, userID int64, opts ...GiftListOption) (OwnedGifts, error) {
	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return OwnedGifts{}, err
	}

	peer := inputPeerFromUser(user)
	if peer == nil {
		return OwnedGifts{}, &Error{Code: 400, Description: "Bad Request: can't resolve user"}
	}

	res, err := b.raw.PaymentsGetSavedStarGifts(ctx, giftListFrom(opts).request(peer))
	if err != nil {
		return OwnedGifts{}, asAPIError(err)
	}

	return b.ownedGiftsFromTg(res, 0, 0), nil
}

// GetChatGifts returns the gifts owned by a chat.
func (b *Bot) GetChatGifts(ctx context.Context, chat ChatID, opts ...GiftListOption) (OwnedGifts, error) {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return OwnedGifts{}, err
	}

	channel, ok := peer.(*tg.InputPeerChannel)
	if !ok {
		return OwnedGifts{}, &Error{Code: 400, Description: "Bad Request: chat is not a supergroup or channel"}
	}

	res, err := b.raw.PaymentsGetSavedStarGifts(ctx, giftListFrom(opts).request(peer))
	if err != nil {
		return OwnedGifts{}, asAPIError(err)
	}

	return b.ownedGiftsFromTg(res, channel.ChannelID, channel.AccessHash), nil
}

// GetBusinessAccountGifts returns the gifts received and owned by a managed
// business account. Requires the can_view_gifts_and_stars business bot right.
func (b *Bot) GetBusinessAccountGifts(ctx context.Context, businessConnectionID string, opts ...GiftListOption) (OwnedGifts, error) {
	var res tg.PaymentsSavedStarGifts

	req := giftListFrom(opts).request(&tg.InputPeerSelf{})
	if err := b.invokeBusiness(ctx, businessConnectionID, req, &res); err != nil {
		return OwnedGifts{}, asAPIError(err)
	}

	return b.ownedGiftsFromTg(&res, 0, 0), nil
}

// inputPeerFromUser adapts a resolved input user to an input peer for the
// saved-gifts request.
func inputPeerFromUser(u tg.InputUserClass) tg.InputPeerClass {
	switch u := u.(type) {
	case *tg.InputUser:
		return &tg.InputPeerUser{UserID: u.UserID, AccessHash: u.AccessHash}
	case *tg.InputUserSelf:
		return &tg.InputPeerSelf{}
	default:
		return nil
	}
}

// giftListFrom collapses the options into a config.
func giftListFrom(opts []GiftListOption) giftListConfig {
	var cfg giftListConfig

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// ownedGiftsFromTg converts a saved-star-gifts response into the Bot API
// OwnedGifts. ownerChannelID/ownerHash, when non-zero, are used to mint the
// owned_gift_id of chat-owned gifts.
func (b *Bot) ownedGiftsFromTg(res *tg.PaymentsSavedStarGifts, ownerChannelID, ownerHash int64) OwnedGifts {
	users := usersByID(res.Users)
	chats := chatsByID(res.Chats)

	out := OwnedGifts{
		TotalCount: res.Count,
		NextOffset: res.NextOffset,
		Gifts:      make([]OwnedGift, 0, len(res.Gifts)),
	}

	for i := range res.Gifts {
		out.Gifts = append(out.Gifts, ownedGiftFromTg(res.Gifts[i], users, chats, ownerChannelID, ownerHash))
	}

	return out
}

// ownedGiftFromTg converts a single saved star gift into an OwnedGift.
func ownedGiftFromTg(
	sg tg.SavedStarGift, users map[int64]*tg.User, chats map[int64]tg.ChatClass, ownerChannelID, ownerHash int64,
) OwnedGift {
	id := mintOwnedGiftID(sg, ownerChannelID, ownerHash)
	sender := senderUser(sg, users)

	switch gift := sg.Gift.(type) {
	case *tg.StarGiftUnique:
		owned := OwnedGiftUnique{
			Type:        ownedGiftUnique,
			Gift:        uniqueGiftFromTg(gift, chats),
			OwnedGiftID: id,
			SenderUser:  sender,
			SendDate:    sg.Date,
			IsSaved:     !sg.Unsaved,
		}

		if v, ok := sg.GetTransferStars(); ok {
			owned.TransferStarCount = int(v)
		}

		if v, ok := sg.GetCanTransferAt(); ok {
			owned.CanBeTransferred = true
			owned.NextTransferDate = v
		}

		return owned
	default:
		return regularOwnedGift(sg, id, sender, chats)
	}
}

// regularOwnedGift builds an OwnedGiftRegular from a saved gift. A non-StarGift
// payload yields an empty Gift.
func regularOwnedGift(sg tg.SavedStarGift, id string, sender *User, chats map[int64]tg.ChatClass) OwnedGift {
	owned := OwnedGiftRegular{
		Type:          ownedGiftRegular,
		OwnedGiftID:   id,
		SenderUser:    sender,
		SendDate:      sg.Date,
		IsPrivate:     sg.NameHidden,
		IsSaved:       !sg.Unsaved,
		CanBeUpgraded: sg.CanUpgrade,
		WasRefunded:   sg.Refunded,
	}

	if g, ok := sg.Gift.(*tg.StarGift); ok {
		owned.Gift = giftFromTg(g, chats)
	}

	if msg, ok := sg.GetMessage(); ok {
		owned.Text = msg.Text
		owned.Entities = entitiesFromTg(msg.Entities)
	}

	if v, ok := sg.GetConvertStars(); ok {
		owned.ConvertStarCount = int(v)
	}

	if v, ok := sg.GetUpgradeStars(); ok {
		owned.PrepaidUpgradeStarCount = int(v)
	}

	return owned
}

// mintOwnedGiftID derives the manageable owned_gift_id for a saved gift, or ""
// when the gift cannot be addressed (e.g. a user gift without a message id).
func mintOwnedGiftID(sg tg.SavedStarGift, ownerChannelID, ownerHash int64) string {
	if msgID, ok := sg.GetMsgID(); ok {
		return OwnedGiftFromMessage(msgID)
	}

	if savedID, ok := sg.GetSavedID(); ok && ownerChannelID != 0 {
		return ownedGiftFromChannel(ownerChannelID, ownerHash, savedID)
	}

	return ""
}

// senderUser resolves the gift sender into a Bot API User.
func senderUser(sg tg.SavedStarGift, users map[int64]*tg.User) *User {
	from, ok := sg.GetFromID()
	if !ok {
		return nil
	}

	pu, ok := from.(*tg.PeerUser)
	if !ok {
		return nil
	}

	var user User

	if u, ok := users[pu.UserID]; ok {
		user = userFromTgUser(u)
	} else {
		user = User{ID: pu.UserID}
	}

	return &user
}

// uniqueGiftFromTg converts a collectible star gift into the Bot API UniqueGift.
func uniqueGiftFromTg(g *tg.StarGiftUnique, chats map[int64]tg.ChatClass) UniqueGift {
	out := UniqueGift{
		GiftID:    strconv.FormatInt(g.GiftID, 10),
		BaseName:  g.Title,
		Name:      g.Slug,
		Number:    g.Num,
		IsPremium: g.RequirePremium,
		IsBurned:  g.Burned,
	}

	for _, attr := range g.Attributes {
		switch a := attr.(type) {
		case *tg.StarGiftAttributeModel:
			out.Model = giftModelFromTg(a)
		case *tg.StarGiftAttributePattern:
			out.Symbol = giftSymbolFromTg(a)
		case *tg.StarGiftAttributeBackdrop:
			out.Backdrop = giftBackdropFromTg(a)
		}
	}

	if peer, ok := g.GetReleasedBy(); ok {
		if chat, ok := chatFromPublisher(peer, chats); ok {
			out.PublisherChat = &chat
		}
	}

	return out
}

// giftModelFromTg converts a model attribute, including its sticker and rarity.
func giftModelFromTg(a *tg.StarGiftAttributeModel) UniqueGiftModel {
	perMille, rarity := rarityValues(a.Rarity)

	m := UniqueGiftModel{Name: a.Name, RarityPerMille: perMille, Rarity: rarity}
	if doc, ok := a.Document.(*tg.Document); ok {
		m.Sticker = stickerFromDocument(doc, "", StickerRegular)
	}

	return m
}

// giftSymbolFromTg converts a pattern attribute into a unique gift symbol.
func giftSymbolFromTg(a *tg.StarGiftAttributePattern) UniqueGiftSymbol {
	perMille, _ := rarityValues(a.Rarity)

	s := UniqueGiftSymbol{Name: a.Name, RarityPerMille: perMille}
	if doc, ok := a.Document.(*tg.Document); ok {
		s.Sticker = stickerFromDocument(doc, "", StickerRegular)
	}

	return s
}

// giftBackdropFromTg converts a backdrop attribute into a unique gift backdrop.
func giftBackdropFromTg(a *tg.StarGiftAttributeBackdrop) UniqueGiftBackdrop {
	perMille, _ := rarityValues(a.Rarity)

	return UniqueGiftBackdrop{
		Name:           a.Name,
		RarityPerMille: perMille,
		Colors: UniqueGiftBackdropColors{
			CenterColor: a.CenterColor,
			EdgeColor:   a.EdgeColor,
			SymbolColor: a.PatternColor,
			TextColor:   a.TextColor,
		},
	}
}

// rarityValues extracts the per-mille count and rarity name from a gift rarity.
// The per-mille variant carries the count; the named variants carry the name.
func rarityValues(r tg.StarGiftAttributeRarityClass) (perMille int, name string) {
	switch r := r.(type) {
	case *tg.StarGiftAttributeRarity:
		return r.Permille, ""
	case *tg.StarGiftAttributeRarityUncommon:
		return 0, "uncommon"
	case *tg.StarGiftAttributeRarityRare:
		return 0, "rare"
	case *tg.StarGiftAttributeRarityEpic:
		return 0, "epic"
	case *tg.StarGiftAttributeRarityLegendary:
		return 0, "legendary"
	default:
		return 0, ""
	}
}
