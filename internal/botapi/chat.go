package botapi

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func convertToBotAPIChatPermissions(p tg.ChatBannedRights) oas.ChatPermissions {
	return oas.ChatPermissions{
		CanSendMessages:       oas.NewOptBool(p.SendMessages),
		CanSendMediaMessages:  oas.NewOptBool(p.SendMedia),
		CanSendPolls:          oas.NewOptBool(p.SendPolls),
		CanSendOtherMessages:  oas.NewOptBool(p.SendGames || p.SendStickers || p.SendInline),
		CanAddWebPagePreviews: oas.NewOptBool(p.EmbedLinks),
		CanChangeInfo:         oas.NewOptBool(p.ChangeInfo),
		CanInviteUsers:        oas.NewOptBool(p.InviteUsers),
		CanPinMessages:        oas.NewOptBool(p.PinMessages),
	}
}

// ApproveChatJoinRequest implements oas.Handler.
func (b *BotAPI) ApproveChatJoinRequest(ctx context.Context, req oas.ApproveChatJoinRequest) (oas.Result, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "resolve chatID")
	}
	user, err := b.resolveUserID(ctx, req.UserID)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "resolve userID")
	}

	if err := p.InviteLinks().ApproveJoin(ctx, user.InputUser()); err != nil {
		return oas.Result{}, err
	}
	return resultOK(true), nil
}

// DeclineChatJoinRequest implements oas.Handler.
func (b *BotAPI) DeclineChatJoinRequest(ctx context.Context, req oas.DeclineChatJoinRequest) (oas.Result, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "resolve chatID")
	}
	user, err := b.resolveUserID(ctx, req.UserID)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "resolve userID")
	}

	if err := p.InviteLinks().DeclineJoin(ctx, user.InputUser()); err != nil {
		return oas.Result{}, err
	}
	return resultOK(true), nil
}

// DeleteChatPhoto implements oas.Handler.
func (b *BotAPI) DeleteChatPhoto(ctx context.Context, req oas.DeleteChatPhoto) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// DeleteChatStickerSet implements oas.Handler.
func (b *BotAPI) DeleteChatStickerSet(ctx context.Context, req oas.DeleteChatStickerSet) (oas.Result, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "resolve chatID")
	}

	ch, ok := p.(peers.Channel)
	if !ok {
		return oas.Result{}, &BadRequestError{Message: "Bad Request: method is available only for supergroups"}
	}

	s, ok := ch.ToSupergroup()
	if !ok {
		return oas.Result{}, &BadRequestError{Message: "Bad Request: method is available only for supergroups"}
	}

	if err := s.ResetStickerSet(ctx); err != nil {
		return oas.Result{}, err
	}
	return resultOK(true), nil
}

// GetChat implements oas.Handler.
func (b *BotAPI) GetChat(ctx context.Context, req oas.GetChat) (oas.ResultChat, error) {
	p, err := b.resolveID(ctx, req.ChatID)
	if err != nil {
		return oas.ResultChat{}, errors.Wrap(err, "resolve chatID")
	}

	var (
		id   = p.TDLibPeerID()
		chat oas.Chat
	)
	switch p := p.(type) {
	case peers.User:
		chat = fillBotAPIChatPrivate(p)
		raw := p.Raw()
		full, err := p.FullRaw(ctx)
		if err != nil {
			return oas.ResultChat{}, errors.Wrap(err, "get full")
		}
		chat.Photo = b.setUserPhoto(id, raw.AccessHash, raw.Photo)
		chat.Bio = optString(full.GetAbout)
		if full.PrivateForwardName != "" {
			chat.HasPrivateForwards.SetTo(true)
		}
	case peers.Chat:
		chat = fillBotAPIChatGroup(p)
		full, err := p.FullRaw(ctx)
		if err != nil {
			return oas.ResultChat{}, errors.Wrap(err, "get full")
		}
		chat.Photo = b.setChatPhoto(id, 0, p.Raw().Photo)
		chat.Description = oas.NewOptString(full.GetAbout())
		if invite, ok := full.GetExportedInvite(); ok {
			chat.InviteLink.SetTo(invite.Link)
		}
		// TODO(tdakkota): resolve pinned.
		if v, ok := p.DefaultBannedRights(); ok {
			chat.Permissions.SetTo(convertToBotAPIChatPermissions(v))
		}
		// TODO(tdakkota): set AllMembersAreAdministrators
	case peers.Channel:
		chat = fillBotAPIChatGroup(p)
		raw := p.Raw()
		full, err := p.FullRaw(ctx)
		if err != nil {
			return oas.ResultChat{}, errors.Wrap(err, "get full")
		}
		chat.Photo = b.setChatPhoto(id, raw.AccessHash, raw.Photo)
		chat.Description = oas.NewOptString(full.GetAbout())
		if invite, ok := full.GetExportedInvite(); ok {
			chat.InviteLink.SetTo(invite.Link)
		}
		// TODO(tdakkota): resolve pinned.
		if v, ok := p.DefaultBannedRights(); ok {
			chat.Permissions.SetTo(convertToBotAPIChatPermissions(v))
		}
		if s, ok := full.GetSlowmodeSeconds(); ok {
			chat.SlowModeDelay.SetTo(s)
		}
		if s, ok := full.GetStickerset(); ok {
			chat.StickerSetName.SetTo(s.ShortName)
		}
		chat.LinkedChatID = optInt64(full.GetLinkedChatID)
		if loc, ok := full.Location.(*tg.ChannelLocation); ok {
			if p, ok := loc.GeoPoint.AsNotEmpty(); ok {
				chat.Location.SetTo(oas.ChatLocation{
					Location: convertToBotAPILocation(p),
					Address:  loc.Address,
				})
			}
		}
		// TODO(tdakkota): set AllMembersAreAdministrators
	}

	return oas.ResultChat{
		Result: oas.NewOptChat(chat),
		Ok:     true,
	}, nil
}

// SetChatAdministratorCustomTitle implements oas.Handler.
func (b *BotAPI) SetChatAdministratorCustomTitle(ctx context.Context, req oas.SetChatAdministratorCustomTitle) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// SetChatDescription implements oas.Handler.
func (b *BotAPI) SetChatDescription(ctx context.Context, req oas.SetChatDescription) (oas.Result, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "resolve chatID")
	}
	if err := p.SetDescription(ctx, req.Description.Value); err != nil {
		return oas.Result{}, err
	}
	return resultOK(true), nil
}

// SetChatPermissions implements oas.Handler.
func (b *BotAPI) SetChatPermissions(ctx context.Context, req oas.SetChatPermissions) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// SetChatPhoto implements oas.Handler.
func (b *BotAPI) SetChatPhoto(ctx context.Context, req oas.SetChatPhoto) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// SetChatStickerSet implements oas.Handler.
func (b *BotAPI) SetChatStickerSet(ctx context.Context, req oas.SetChatStickerSet) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// SetChatTitle implements oas.Handler.
func (b *BotAPI) SetChatTitle(ctx context.Context, req oas.SetChatTitle) (oas.Result, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "resolve chatID")
	}
	if err := p.SetTitle(ctx, req.Title); err != nil {
		return oas.Result{}, err
	}
	return resultOK(true), nil
}

// LeaveChat implements oas.Handler.
func (b *BotAPI) LeaveChat(ctx context.Context, req oas.LeaveChat) (oas.Result, error) {
	p, err := b.resolveIDToChat(ctx, req.ChatID)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "resolve chatID")
	}
	if !p.Left() {
		if err := p.Leave(ctx); err != nil {
			return oas.Result{}, err
		}
	}
	return resultOK(true), nil
}
