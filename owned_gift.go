package botapi

import (
	"strconv"
	"strings"

	"github.com/gotd/td/tg"
)

// Owned gift id encoding.
//
// The Bot API addresses a gift owned by an account with an opaque
// owned_gift_id string. botapi has no TDLib to mint TDLib's format, so it
// defines its own: a prefixed, ":"-delimited encoding of the MTProto
// inputSavedStarGift the gift maps to. The encoding round-trips through
// ownedGiftToTg and is produced by the OwnedGiftFrom* constructors (and, in the
// future, by the gift-listing methods).
const (
	ownedGiftMessagePrefix = "msg"
	ownedGiftChatPrefix    = "chat"
	ownedGiftSlugPrefix    = "slug"
)

// OwnedGiftFromMessage builds an owned_gift_id for a gift received by the
// account from the id of the service message that carries the gift (the
// messageActionStarGift). This is the form for gifts owned by a user, including
// a connected business account.
func OwnedGiftFromMessage(messageID int) string {
	return ownedGiftMessagePrefix + ":" + strconv.Itoa(messageID)
}

// OwnedGiftFromSlug builds an owned_gift_id for a unique (collectible) gift
// addressed by its public link slug.
func OwnedGiftFromSlug(slug string) string {
	return ownedGiftSlugPrefix + ":" + slug
}

// ownedGiftToTg decodes an owned_gift_id into the MTProto inputSavedStarGift it
// addresses.
func ownedGiftToTg(id string) (tg.InputSavedStarGiftClass, error) {
	prefix, rest, ok := strings.Cut(id, ":")
	if !ok {
		return nil, errInvalidOwnedGift()
	}

	switch prefix {
	case ownedGiftMessagePrefix:
		msgID, err := strconv.Atoi(rest)
		if err != nil {
			return nil, errInvalidOwnedGift()
		}

		return &tg.InputSavedStarGiftUser{MsgID: msgID}, nil
	case ownedGiftSlugPrefix:
		if rest == "" {
			return nil, errInvalidOwnedGift()
		}

		return &tg.InputSavedStarGiftSlug{Slug: rest}, nil
	case ownedGiftChatPrefix:
		return ownedGiftChatToTg(rest)
	default:
		return nil, errInvalidOwnedGift()
	}
}

// ownedGiftChatToTg decodes the "chat:<channelID>:<accessHash>:<savedID>" form
// into an inputSavedStarGiftChat.
func ownedGiftChatToTg(rest string) (tg.InputSavedStarGiftClass, error) {
	parts := strings.Split(rest, ":")
	if len(parts) != 3 {
		return nil, errInvalidOwnedGift()
	}

	channelID, err1 := strconv.ParseInt(parts[0], 10, 64)
	accessHash, err2 := strconv.ParseInt(parts[1], 10, 64)
	savedID, err3 := strconv.ParseInt(parts[2], 10, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return nil, errInvalidOwnedGift()
	}

	return &tg.InputSavedStarGiftChat{
		Peer:    &tg.InputPeerChannel{ChannelID: channelID, AccessHash: accessHash},
		SavedID: savedID,
	}, nil
}

// errInvalidOwnedGift is returned for a malformed owned_gift_id.
func errInvalidOwnedGift() error {
	return &Error{Code: 400, Description: "Bad Request: invalid owned_gift_id"}
}
