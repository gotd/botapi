package botapi

import (
	"context"
	"time"

	"github.com/gotd/td/telegram/peers"
)

// resolveInviteLinks resolves a ChatID to its invite-link manager.
func (b *Bot) resolveInviteLinks(ctx context.Context, chat ChatID) (peers.InviteLinks, error) {
	p, err := b.resolvePeer(ctx, chat)
	if err != nil {
		return peers.InviteLinks{}, err
	}
	type linkable interface {
		InviteLinks() peers.InviteLinks
	}
	l, ok := p.(linkable)
	if !ok {
		return peers.InviteLinks{}, &Error{Code: 400, Description: "Bad Request: method is not available in private chats"}
	}
	return l.InviteLinks(), nil
}

// convertInviteLink converts a gotd invite link into the Bot API type. It
// resolves the creator, which may require a network round trip.
func (b *Bot) convertInviteLink(ctx context.Context, link peers.InviteLink) (*ChatInviteLink, error) {
	out := &ChatInviteLink{
		InviteLink:         link.Link(),
		CreatesJoinRequest: link.RequestNeeded(),
		IsPrimary:          link.Permanent(),
		IsRevoked:          link.Revoked(),
	}
	if creator, err := link.Creator(ctx); err == nil {
		out.Creator = userFromTgUser(creator.Raw())
	}
	if name, ok := link.Title(); ok {
		out.Name = name
	}
	if expire, ok := link.ExpireDate(); ok {
		out.ExpireDate = int(expire.Unix())
	}
	if limit, ok := link.UsageLimit(); ok {
		out.MemberLimit = limit
	}
	if requested, ok := link.Requested(); ok {
		out.PendingJoinRequestCount = requested
	}
	return out, nil
}

// InviteLinkOption configures a create/edit invite-link call.
type InviteLinkOption func(*peers.ExportLinkOptions)

// WithInviteLinkName sets the invite link name (0-32 characters).
func WithInviteLinkName(name string) InviteLinkOption {
	return func(o *peers.ExportLinkOptions) { o.Title = name }
}

// WithInviteLinkExpire sets the Unix time when the link will expire.
func WithInviteLinkExpire(unixTime int) InviteLinkOption {
	return func(o *peers.ExportLinkOptions) { o.ExpireDate = time.Unix(int64(unixTime), 0) }
}

// WithInviteLinkMemberLimit caps how many users may join via this link.
func WithInviteLinkMemberLimit(limit int) InviteLinkOption {
	return func(o *peers.ExportLinkOptions) { o.UsageLimit = limit }
}

// WithInviteLinkJoinRequest requires administrators to approve users joining via
// this link.
func WithInviteLinkJoinRequest() InviteLinkOption {
	return func(o *peers.ExportLinkOptions) { o.RequestNeeded = true }
}

// ExportChatInviteLink generates a new primary invite link for a chat, revoking
// any previous one, and returns it.
func (b *Bot) ExportChatInviteLink(ctx context.Context, chat ChatID) (string, error) {
	links, err := b.resolveInviteLinks(ctx, chat)
	if err != nil {
		return "", err
	}
	link, err := links.ExportNew(ctx, peers.ExportLinkOptions{})
	if err != nil {
		return "", asAPIError(err)
	}
	return link.Link(), nil
}

// CreateChatInviteLink creates an additional invite link for a chat.
func (b *Bot) CreateChatInviteLink(ctx context.Context, chat ChatID, opts ...InviteLinkOption) (*ChatInviteLink, error) {
	links, err := b.resolveInviteLinks(ctx, chat)
	if err != nil {
		return nil, err
	}
	var o peers.ExportLinkOptions
	for _, opt := range opts {
		opt(&o)
	}
	link, err := links.AddNew(ctx, o)
	if err != nil {
		return nil, asAPIError(err)
	}
	return b.convertInviteLink(ctx, link)
}

// EditChatInviteLink edits a non-primary invite link created by the bot.
func (b *Bot) EditChatInviteLink(ctx context.Context, chat ChatID, inviteLink string, opts ...InviteLinkOption) (*ChatInviteLink, error) {
	links, err := b.resolveInviteLinks(ctx, chat)
	if err != nil {
		return nil, err
	}
	var o peers.ExportLinkOptions
	for _, opt := range opts {
		opt(&o)
	}
	link, err := links.Edit(ctx, inviteLink, o)
	if err != nil {
		return nil, asAPIError(err)
	}
	return b.convertInviteLink(ctx, link)
}

// RevokeChatInviteLink revokes an invite link created by the bot.
func (b *Bot) RevokeChatInviteLink(ctx context.Context, chat ChatID, inviteLink string) (*ChatInviteLink, error) {
	links, err := b.resolveInviteLinks(ctx, chat)
	if err != nil {
		return nil, err
	}
	link, err := links.Revoke(ctx, inviteLink)
	if err != nil {
		return nil, asAPIError(err)
	}
	return b.convertInviteLink(ctx, link)
}
