package botapi

import (
	"context"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
)

// BusinessBotRights lists the rights a business bot has over the connected
// business account.
type BusinessBotRights struct {
	CanReply                   bool `json:"can_reply,omitempty"`
	CanReadMessages            bool `json:"can_read_messages,omitempty"`
	CanDeleteSentMessages      bool `json:"can_delete_sent_messages,omitempty"`
	CanDeleteAllMessages       bool `json:"can_delete_all_messages,omitempty"`
	CanEditName                bool `json:"can_edit_name,omitempty"`
	CanEditBio                 bool `json:"can_edit_bio,omitempty"`
	CanEditProfilePhoto        bool `json:"can_edit_profile_photo,omitempty"`
	CanEditUsername            bool `json:"can_edit_username,omitempty"`
	CanChangeGiftSettings      bool `json:"can_change_gift_settings,omitempty"`
	CanViewGiftsAndStars       bool `json:"can_view_gifts_and_stars,omitempty"`
	CanConvertGiftsToStars     bool `json:"can_convert_gifts_to_stars,omitempty"`
	CanTransferAndUpgradeGifts bool `json:"can_transfer_and_upgrade_gifts,omitempty"`
	CanTransferStars           bool `json:"can_transfer_stars,omitempty"`
	CanManageStories           bool `json:"can_manage_stories,omitempty"`
}

// businessBotRightsFromTg converts MTProto business bot rights to the Bot API
// representation.
func businessBotRightsFromTg(r tg.BusinessBotRights) BusinessBotRights {
	return BusinessBotRights{
		CanReply:                   r.Reply,
		CanReadMessages:            r.ReadMessages,
		CanDeleteSentMessages:      r.DeleteSentMessages,
		CanDeleteAllMessages:       r.DeleteReceivedMessages,
		CanEditName:                r.EditName,
		CanEditBio:                 r.EditBio,
		CanEditProfilePhoto:        r.EditProfilePhoto,
		CanEditUsername:            r.EditUsername,
		CanChangeGiftSettings:      r.ChangeGiftSettings,
		CanViewGiftsAndStars:       r.ViewGifts,
		CanConvertGiftsToStars:     r.SellGifts,
		CanTransferAndUpgradeGifts: r.TransferAndUpgradeGifts,
		CanTransferStars:           r.TransferStars,
		CanManageStories:           r.ManageStories,
	}
}

// BusinessConnection describes the connection of the bot with a business
// account.
type BusinessConnection struct {
	// ID is the unique identifier of the business connection.
	ID string `json:"id"`
	// User is the business account user that connected to the bot.
	User User `json:"user"`
	// UserChatID is the identifier of the private chat with the user.
	UserChatID int64 `json:"user_chat_id"`
	// Date is the Unix time when the connection was established.
	Date int `json:"date"`
	// Rights are the rights of the business bot.
	Rights *BusinessBotRights `json:"rights,omitempty"`
	// IsEnabled reports whether the connection is active.
	IsEnabled bool `json:"is_enabled"`
}

// invokeBusiness sends an inner RPC on behalf of a connected business account,
// wrapping it in invokeWithBusinessConnection. *tg.Client has no generated
// helper for this wrapper, so it goes through the raw invoker.
func (b *Bot) invokeBusiness(ctx context.Context, connectionID string, query bin.Object, output bin.Decoder) error {
	return b.invoker.Invoke(ctx, &tg.InvokeWithBusinessConnectionRequest{
		ConnectionID: connectionID,
		Query:        query,
	}, output)
}

// GetBusinessConnection returns information about the connection of the bot with
// a business account by the connection id received in updates.
func (b *Bot) GetBusinessConnection(ctx context.Context, businessConnectionID string) (*BusinessConnection, error) {
	res, err := b.raw.AccountGetBotBusinessConnection(ctx, businessConnectionID)
	if err != nil {
		return nil, asAPIError(err)
	}

	updates, users := updatesAndUsers(res)

	for _, upd := range updates {
		connect, ok := upd.(*tg.UpdateBotBusinessConnect)
		if !ok {
			continue
		}

		return b.businessConnectionFromTg(connect.Connection, users), nil
	}

	return nil, &Error{Code: 400, Description: "Bad Request: business connection not found"}
}

// businessConnectionFromTg converts an MTProto business connection into the Bot
// API type, resolving the connected user from the update's users.
func (b *Bot) businessConnectionFromTg(c tg.BotBusinessConnection, users map[int64]*tg.User) *BusinessConnection {
	out := &BusinessConnection{
		ID:         c.ConnectionID,
		UserChatID: c.UserID,
		Date:       c.Date,
		IsEnabled:  !c.Disabled,
	}

	if u, ok := users[c.UserID]; ok {
		out.User = userFromTgUser(u)
	} else {
		out.User = User{ID: c.UserID}
	}

	rights := businessBotRightsFromTg(c.Rights)

	out.Rights = &rights

	return out
}

// SetBusinessAccountName changes the first and last name of a managed business
// account. The bot must have the can_edit_name business bot right.
func (b *Bot) SetBusinessAccountName(ctx context.Context, businessConnectionID, firstName, lastName string) error {
	req := &tg.AccountUpdateProfileRequest{}
	req.SetFirstName(firstName)
	req.SetLastName(lastName)

	return b.invokeBusiness(ctx, businessConnectionID, req, &tg.UserBox{})
}

// SetBusinessAccountBio changes the bio of a managed business account. The bot
// must have the can_edit_bio business bot right.
func (b *Bot) SetBusinessAccountBio(ctx context.Context, businessConnectionID, bio string) error {
	req := &tg.AccountUpdateProfileRequest{}
	req.SetAbout(bio)

	return b.invokeBusiness(ctx, businessConnectionID, req, &tg.UserBox{})
}

// SetBusinessAccountUsername changes the username of a managed business account.
// The bot must have the can_edit_username business bot right.
func (b *Bot) SetBusinessAccountUsername(ctx context.Context, businessConnectionID, username string) error {
	return b.invokeBusiness(ctx, businessConnectionID, &tg.AccountUpdateUsernameRequest{Username: username}, &tg.UserBox{})
}

// GetBusinessAccountStarBalance returns the amount of Telegram Stars owned by a
// managed business account. The bot must have the can_view_gifts_and_stars
// business bot right.
func (b *Bot) GetBusinessAccountStarBalance(ctx context.Context, businessConnectionID string) (StarAmount, error) {
	var status tg.PaymentsStarsStatus

	err := b.invokeBusiness(ctx, businessConnectionID, &tg.PaymentsGetStarsStatusRequest{
		Peer: &tg.InputPeerSelf{},
	}, &status)
	if err != nil {
		return StarAmount{}, asAPIError(err)
	}

	amount, ok := status.Balance.(*tg.StarsAmount)
	if !ok {
		return StarAmount{}, &Error{Code: 500, Description: "Internal Server Error: unexpected stars balance"}
	}

	return StarAmount{Amount: int(amount.Amount), NanostarAmount: amount.Nanos}, nil
}

// DeleteBusinessMessages deletes messages on behalf of a managed business
// account. The messages must all belong to the same chat. The bot must have the
// can_delete_sent_messages right (for its own messages) or can_delete_all_messages
// right.
func (b *Bot) DeleteBusinessMessages(ctx context.Context, businessConnectionID string, messageIDs []int) error {
	if len(messageIDs) == 0 {
		return nil
	}

	var affected tg.MessagesAffectedMessages

	err := b.invokeBusiness(ctx, businessConnectionID, &tg.MessagesDeleteMessagesRequest{
		Revoke: true,
		ID:     messageIDs,
	}, &affected)
	if err != nil {
		return asAPIError(err)
	}

	return nil
}

// updatesAndUsers unpacks the update list and a user-by-id index from an
// Updates response.
func updatesAndUsers(resp tg.UpdatesClass) (updates []tg.UpdateClass, users map[int64]*tg.User) {
	switch u := resp.(type) {
	case *tg.Updates:
		return u.Updates, usersByID(u.Users)
	case *tg.UpdatesCombined:
		return u.Updates, usersByID(u.Users)
	default:
		return nil, nil
	}
}
