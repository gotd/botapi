package botapi

import (
	"context"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

// errUnexpectedInvite is returned when an invite RPC answers with an unexpected
// constructor.
func errUnexpectedInvite() *Error {
	return &Error{Code: 500, Description: "Internal Server Error: unexpected invite response"}
}

// inviteLinkName resolves the optional name from invite-link options, ignoring
// the other fields that subscription links do not support.
func inviteLinkName(opts []InviteLinkOption) string {
	var o peers.ExportLinkOptions

	for _, opt := range opts {
		opt(&o)
	}

	return o.Title
}

// inviteLinkFromExported converts an MTProto exported invite into the Bot API
// type. The bot is always the creator of the links it makes.
func (b *Bot) inviteLinkFromExported(inv *tg.ChatInviteExported) *ChatInviteLink {
	out := &ChatInviteLink{
		InviteLink:              inv.Link,
		Creator:                 userFromTgUser(b.self),
		CreatesJoinRequest:      inv.RequestNeeded,
		IsPrimary:               inv.Permanent,
		IsRevoked:               inv.Revoked,
		Name:                    inv.Title,
		ExpireDate:              inv.ExpireDate,
		MemberLimit:             inv.UsageLimit,
		PendingJoinRequestCount: inv.Requested,
	}

	if pricing, ok := inv.GetSubscriptionPricing(); ok {
		out.SubscriptionPeriod = pricing.Period
		out.SubscriptionPrice = int(pricing.Amount)
	}

	return out
}

// CreateChatSubscriptionInviteLink creates a subscription invite link for a
// channel chat. Users joining through the link pay subscriptionPrice Telegram
// Stars every subscriptionPeriod seconds (currently only 2592000, 30 days, is
// supported). The bot must have the can_invite_users administrator right.
func (b *Bot) CreateChatSubscriptionInviteLink(
	ctx context.Context, chat ChatID, subscriptionPeriod, subscriptionPrice int, opts ...InviteLinkOption,
) (*ChatInviteLink, error) {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	req := &tg.MessagesExportChatInviteRequest{Peer: peer}
	req.SetSubscriptionPricing(tg.StarsSubscriptionPricing{
		Period: subscriptionPeriod,
		Amount: int64(subscriptionPrice),
	})

	if name := inviteLinkName(opts); name != "" {
		req.SetTitle(name)
	}

	res, err := b.raw.MessagesExportChatInvite(ctx, req)
	if err != nil {
		return nil, asAPIError(err)
	}

	exported, ok := res.(*tg.ChatInviteExported)
	if !ok {
		return nil, errUnexpectedInvite()
	}

	return b.inviteLinkFromExported(exported), nil
}

// EditChatSubscriptionInviteLink edits the name of a subscription invite link
// created by the bot. The price and period of a subscription link cannot be
// changed.
func (b *Bot) EditChatSubscriptionInviteLink(
	ctx context.Context, chat ChatID, inviteLink string, opts ...InviteLinkOption,
) (*ChatInviteLink, error) {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	req := &tg.MessagesEditExportedChatInviteRequest{Peer: peer, Link: inviteLink}
	req.SetTitle(inviteLinkName(opts))

	res, err := b.raw.MessagesEditExportedChatInvite(ctx, req)
	if err != nil {
		return nil, asAPIError(err)
	}

	edited, ok := res.(*tg.MessagesExportedChatInvite)
	if !ok {
		return nil, errUnexpectedInvite()
	}

	exported, ok := edited.Invite.(*tg.ChatInviteExported)
	if !ok {
		return nil, errUnexpectedInvite()
	}

	return b.inviteLinkFromExported(exported), nil
}
