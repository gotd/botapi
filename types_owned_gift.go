package botapi

// UniqueGiftModel describes the model of a unique gift.
type UniqueGiftModel struct {
	// Name is the name of the model.
	Name string `json:"name"`
	// Sticker is the sticker that represents the model.
	Sticker Sticker `json:"sticker"`
	// RarityPerMille is the number of unique gifts that receive this model for
	// every 1000 gifts upgraded.
	RarityPerMille int `json:"rarity_per_mille"`
	// Rarity is the rarity class of the model, if any: one of "uncommon",
	// "rare", "epic" or "legendary".
	Rarity string `json:"rarity,omitempty"`
}

// UniqueGiftSymbol describes the symbol shown on the pattern of a unique gift.
type UniqueGiftSymbol struct {
	// Name is the name of the symbol.
	Name string `json:"name"`
	// Sticker is the sticker that represents the symbol.
	Sticker Sticker `json:"sticker"`
	// RarityPerMille is the number of unique gifts that receive this symbol for
	// every 1000 gifts upgraded.
	RarityPerMille int `json:"rarity_per_mille"`
}

// UniqueGiftBackdropColors describes the colors of the backdrop of a unique
// gift.
type UniqueGiftBackdropColors struct {
	// CenterColor is the color in the center of the backdrop, in RGB format.
	CenterColor int `json:"center_color"`
	// EdgeColor is the color on the edges of the backdrop, in RGB format.
	EdgeColor int `json:"edge_color"`
	// SymbolColor is the color to be applied to the symbol, in RGB format.
	SymbolColor int `json:"symbol_color"`
	// TextColor is the color for the text on the backdrop, in RGB format.
	TextColor int `json:"text_color"`
}

// UniqueGiftBackdrop describes the backdrop of a unique gift.
type UniqueGiftBackdrop struct {
	// Name is the name of the backdrop.
	Name string `json:"name"`
	// Colors are the colors of the backdrop.
	Colors UniqueGiftBackdropColors `json:"colors"`
	// RarityPerMille is the number of unique gifts that receive this backdrop
	// for every 1000 gifts upgraded.
	RarityPerMille int `json:"rarity_per_mille"`
}

// UniqueGift describes a unique gift that was upgraded from a regular gift.
type UniqueGift struct {
	// GiftID is the unique identifier of the regular gift that was upgraded.
	GiftID string `json:"gift_id,omitempty"`
	// BaseName is the human-readable name of the regular gift from which this
	// unique gift was upgraded.
	BaseName string `json:"base_name"`
	// Name is the unique name of the gift, used to construct its link.
	Name string `json:"name"`
	// Number is the unique number of the upgraded gift among gifts upgraded from
	// the same regular gift.
	Number int `json:"number"`
	// Model is the model of the gift.
	Model UniqueGiftModel `json:"model"`
	// Symbol is the symbol of the gift.
	Symbol UniqueGiftSymbol `json:"symbol"`
	// Backdrop is the backdrop of the gift.
	Backdrop UniqueGiftBackdrop `json:"backdrop"`
	// IsPremium reports whether the gift can be owned only by Premium users.
	IsPremium bool `json:"is_premium,omitempty"`
	// IsBurned reports whether the gift was burned.
	IsBurned bool `json:"is_burned,omitempty"`
	// PublisherChat is the chat that published the gift, if any.
	PublisherChat *Chat `json:"publisher_chat,omitempty"`
}

// OwnedGift is a sealed union describing a gift received and owned by a user or
// a chat.
//
// Concrete variants: OwnedGiftRegular, OwnedGiftUnique.
type OwnedGift interface {
	isOwnedGift()
}

// OwnedGiftRegular describes a regular gift owned by a user or a chat.
type OwnedGiftRegular struct {
	Type string `json:"type"`
	// Gift is the information about the regular gift.
	Gift Gift `json:"gift"`
	// OwnedGiftID is the unique identifier of the gift for use in gift-management
	// methods; only present for gifts received on behalf of business accounts.
	OwnedGiftID string `json:"owned_gift_id,omitempty"`
	// SenderUser is the user that sent the gift, if known.
	SenderUser *User `json:"sender_user,omitempty"`
	// SendDate is the Unix time when the gift was sent.
	SendDate int `json:"send_date"`
	// Text is the text attached to the gift.
	Text string `json:"text,omitempty"`
	// Entities are the special entities that appear in the text.
	Entities []MessageEntity `json:"entities,omitempty"`
	// IsPrivate reports whether the sender and text are shown only to the gift
	// receiver.
	IsPrivate bool `json:"is_private,omitempty"`
	// IsSaved reports whether the gift is displayed on the account's profile.
	IsSaved bool `json:"is_saved,omitempty"`
	// CanBeUpgraded reports whether the gift can be upgraded to a unique gift.
	CanBeUpgraded bool `json:"can_be_upgraded,omitempty"`
	// WasRefunded reports whether the gift was refunded and isn't owned anymore.
	WasRefunded bool `json:"was_refunded,omitempty"`
	// ConvertStarCount is the number of Telegram Stars the gift can be converted
	// to by the owner.
	ConvertStarCount int `json:"convert_star_count,omitempty"`
	// PrepaidUpgradeStarCount is the number of Telegram Stars that were paid for
	// the gift to be upgradable.
	PrepaidUpgradeStarCount int `json:"prepaid_upgrade_star_count,omitempty"`
}

// OwnedGiftUnique describes a unique gift received and owned by a user or a
// chat.
type OwnedGiftUnique struct {
	Type string `json:"type"`
	// Gift is the information about the unique gift.
	Gift UniqueGift `json:"gift"`
	// OwnedGiftID is the unique identifier of the gift for use in gift-management
	// methods; only present for gifts received on behalf of business accounts.
	OwnedGiftID string `json:"owned_gift_id,omitempty"`
	// SenderUser is the user that sent the gift, if known.
	SenderUser *User `json:"sender_user,omitempty"`
	// SendDate is the Unix time when the gift was sent.
	SendDate int `json:"send_date"`
	// IsSaved reports whether the gift is displayed on the account's profile.
	IsSaved bool `json:"is_saved,omitempty"`
	// CanBeTransferred reports whether the gift can be transferred to another
	// owner.
	CanBeTransferred bool `json:"can_be_transferred,omitempty"`
	// TransferStarCount is the number of Telegram Stars that must be paid to
	// transfer the gift.
	TransferStarCount int `json:"transfer_star_count,omitempty"`
	// NextTransferDate is the Unix time when the gift can be transferred; if 0,
	// the gift can be transferred immediately.
	NextTransferDate int `json:"next_transfer_date,omitempty"`
}

func (OwnedGiftRegular) isOwnedGift() {}
func (OwnedGiftUnique) isOwnedGift()  {}

// Owned gift discriminators.
const (
	ownedGiftRegular = "regular"
	ownedGiftUnique  = "unique"
)

// OwnedGifts contains a list of gifts owned by a user or a chat.
type OwnedGifts struct {
	// TotalCount is the total number of gifts owned by the target.
	TotalCount int `json:"total_count"`
	// Gifts is the requested page of owned gifts.
	Gifts []OwnedGift `json:"gifts"`
	// NextOffset is the offset to pass to the next request to get more gifts; an
	// empty string if there are no more gifts.
	NextOffset string `json:"next_offset,omitempty"`
}
