package botapi

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram/peers"

	"github.com/gotd/botapi/internal/oas"
)

func (b *BotAPI) convertInviteLink(ctx context.Context, link peers.InviteLink) (oas.ChatInviteLink, error) {
	creator, err := link.Creator(ctx)
	if err != nil {
		return oas.ChatInviteLink{}, errors.Wrap(err, "get creator")
	}

	raw := link.Raw()
	result := oas.ChatInviteLink{
		InviteLink:              link.Link(),
		Creator:                 convertToBotAPIUser(creator),
		CreatesJoinRequest:      link.RequestNeeded(),
		IsPrimary:               link.Permanent(),
		IsRevoked:               link.Revoked(),
		Name:                    optString(link.Title),
		ExpireDate:              optInt(raw.GetExpireDate),
		MemberLimit:             optInt(link.UsageLimit),
		PendingJoinRequestCount: optInt(link.Requested),
	}
	return result, nil
}

// CreateChatInviteLink implements oas.Handler.
func (b *BotAPI) CreateChatInviteLink(ctx context.Context, req oas.CreateChatInviteLink) (oas.ResultChatInviteLink, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.ResultChatInviteLink{}, errors.Wrap(err, "resolve chatID")
	}
	opts := peers.ExportLinkOptions{
		RequestNeeded: req.CreatesJoinRequest.Value,
		ExpireDate:    time.Time{},
		UsageLimit:    req.MemberLimit.Value,
		Title:         req.Name.Value,
	}
	if u, ok := req.ExpireDate.Get(); ok {
		opts.ExpireDate = time.Unix(int64(u), 0)
	}

	link, err := p.InviteLinks().AddNew(ctx, opts)
	if err != nil {
		return oas.ResultChatInviteLink{}, errors.Wrap(err, "create link")
	}

	result, err := b.convertInviteLink(ctx, link)
	if err != nil {
		return oas.ResultChatInviteLink{}, errors.Wrap(err, "convert link")
	}

	return oas.ResultChatInviteLink{
		Result: oas.NewOptChatInviteLink(result),
		Ok:     true,
	}, nil
}

// EditChatInviteLink implements oas.Handler.
func (b *BotAPI) EditChatInviteLink(ctx context.Context, req oas.EditChatInviteLink) (oas.ResultChatInviteLink, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.ResultChatInviteLink{}, errors.Wrap(err, "resolve chatID")
	}

	opts := peers.ExportLinkOptions{
		RequestNeeded: req.CreatesJoinRequest.Value,
		ExpireDate:    time.Time{},
		UsageLimit:    req.MemberLimit.Value,
		Title:         req.Name.Value,
	}
	if u, ok := req.ExpireDate.Get(); ok {
		opts.ExpireDate = time.Unix(int64(u), 0)
	}

	link, err := p.InviteLinks().Edit(ctx, req.InviteLink, opts)
	if err != nil {
		return oas.ResultChatInviteLink{}, errors.Wrap(err, "edit link")
	}

	result, err := b.convertInviteLink(ctx, link)
	if err != nil {
		return oas.ResultChatInviteLink{}, errors.Wrap(err, "convert link")
	}

	return oas.ResultChatInviteLink{
		Result: oas.NewOptChatInviteLink(result),
		Ok:     true,
	}, nil
}

// ExportChatInviteLink implements oas.Handler.
func (b *BotAPI) ExportChatInviteLink(ctx context.Context, req oas.ExportChatInviteLink) (oas.ResultString, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.ResultString{}, errors.Wrap(err, "resolve chatID")
	}

	link, err := p.InviteLinks().ExportNew(ctx, peers.ExportLinkOptions{})
	if err != nil {
		return oas.ResultString{}, errors.Wrap(err, "export link")
	}

	return oas.ResultString{
		Result: oas.NewOptString(link.Link()),
		Ok:     true,
	}, nil
}

// RevokeChatInviteLink implements oas.Handler.
func (b *BotAPI) RevokeChatInviteLink(ctx context.Context, req oas.RevokeChatInviteLink) (oas.ResultChatInviteLink, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.ResultChatInviteLink{}, errors.Wrap(err, "resolve chatID")
	}

	link, err := p.InviteLinks().Revoke(ctx, req.InviteLink)
	if err != nil {
		return oas.ResultChatInviteLink{}, errors.Wrap(err, "export link")
	}

	result, err := b.convertInviteLink(ctx, link)
	if err != nil {
		return oas.ResultChatInviteLink{}, errors.Wrap(err, "convert link")
	}

	return oas.ResultChatInviteLink{
		Result: oas.NewOptChatInviteLink(result),
		Ok:     true,
	}, nil
}
