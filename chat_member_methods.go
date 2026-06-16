package botapi

import (
	"context"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

// resolveChannel resolves a ChatID to its MTProto input channel. Member
// management is only available in supergroups and channels.
func (b *Bot) resolveChannel(ctx context.Context, chat ChatID) (tg.InputChannelClass, error) {
	p, err := b.resolvePeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	ch, ok := p.(peers.Channel)
	if !ok {
		return nil, &Error{Code: 400, Description: "Bad Request: method is available only for supergroups and channels"}
	}

	return ch.InputChannel(), nil
}

// BanChatMemberOption configures a BanChatMember call.
type BanChatMemberOption func(*banConfig)

type banConfig struct {
	untilDate     int
	revokeHistory bool
}

// WithBanUntil sets the Unix time when the user will be unbanned (0 or far in
// the future means forever).
func WithBanUntil(untilDate int) BanChatMemberOption {
	return func(c *banConfig) { c.untilDate = untilDate }
}

// WithRevokeMessages deletes all messages from the user being banned.
func WithRevokeMessages() BanChatMemberOption {
	return func(c *banConfig) { c.revokeHistory = true }
}

// BanChatMember bans a user from a supergroup or channel. The bot must be an
// administrator with the can_restrict_members right.
func (b *Bot) BanChatMember(ctx context.Context, chat ChatID, userID int64, opts ...BanChatMemberOption) error {
	var cfg banConfig

	for _, o := range opts {
		o(&cfg)
	}

	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	if _, err := b.raw.ChannelsEditBanned(ctx, &tg.ChannelsEditBannedRequest{
		Channel:     channel,
		Participant: userToInputPeer(user),
		BannedRights: tg.ChatBannedRights{
			ViewMessages: true,
			SendMessages: true,
			SendMedia:    true,
			SendStickers: true,
			SendGifs:     true,
			SendGames:    true,
			SendInline:   true,
			EmbedLinks:   true,
			UntilDate:    cfg.untilDate,
		},
	}); err != nil {
		return asAPIError(err)
	}

	if cfg.revokeHistory {
		if _, err := b.raw.ChannelsDeleteParticipantHistory(ctx, &tg.ChannelsDeleteParticipantHistoryRequest{
			Channel:     channel,
			Participant: userToInputPeer(user),
		}); err != nil {
			return asAPIError(err)
		}
	}

	return nil
}

// UnbanChatMember removes a user's ban from a supergroup or channel, clearing
// all restrictions. It does not re-add the user to the chat.
func (b *Bot) UnbanChatMember(ctx context.Context, chat ChatID, userID int64) error {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	if _, err := b.raw.ChannelsEditBanned(ctx, &tg.ChannelsEditBannedRequest{
		Channel:      channel,
		Participant:  userToInputPeer(user),
		BannedRights: tg.ChatBannedRights{}, // empty rights clear the ban
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// RestrictChatMember restricts a user in a supergroup. permissions is an
// allow-list; denied actions are turned into MTProto banned rights. untilDate is
// a Unix time (0 means forever).
func (b *Bot) RestrictChatMember(ctx context.Context, chat ChatID, userID int64, permissions ChatPermissions, untilDate int) error {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	if _, err := b.raw.ChannelsEditBanned(ctx, &tg.ChannelsEditBannedRequest{
		Channel:      channel,
		Participant:  userToInputPeer(user),
		BannedRights: permissions.toBannedRights(untilDate),
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// PromoteChatMember promotes or demotes a user in a supergroup or channel. The
// granted rights are the true fields of rights; pass a zero value to demote.
func (b *Bot) PromoteChatMember(ctx context.Context, chat ChatID, userID int64, rights ChatAdminRights) error {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	if _, err := b.raw.ChannelsEditAdmin(ctx, &tg.ChannelsEditAdminRequest{
		Channel:     channel,
		UserID:      user,
		AdminRights: rights.toTg(),
		Rank:        rights.CustomTitle,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// GetChatMember returns information about a member of a supergroup or channel.
func (b *Bot) GetChatMember(ctx context.Context, chat ChatID, userID int64) (ChatMember, error) {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return nil, err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.ChannelsGetParticipant(ctx, &tg.ChannelsGetParticipantRequest{
		Channel:     channel,
		Participant: userToInputPeer(user),
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	return chatMemberFromParticipant(res.Participant, usersByID(res.Users)), nil
}

// GetChatAdministrators returns the administrators of a supergroup or channel,
// excluding other bots.
func (b *Bot) GetChatAdministrators(ctx context.Context, chat ChatID) ([]ChatMember, error) {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.ChannelsGetParticipants(ctx, &tg.ChannelsGetParticipantsRequest{
		Channel: channel,
		Filter:  &tg.ChannelParticipantsAdmins{},
		Limit:   200,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	participants, ok := res.AsModified()
	if !ok {
		return nil, &Error{Code: 500, Description: "Internal Server Error: unexpected participants response"}
	}

	users := usersByID(participants.Users)
	out := make([]ChatMember, 0, len(participants.Participants))

	for _, p := range participants.Participants {
		out = append(out, chatMemberFromParticipant(p, users))
	}

	return out, nil
}

// GetChatMemberCount returns the number of members in a supergroup or channel.
func (b *Bot) GetChatMemberCount(ctx context.Context, chat ChatID) (int, error) {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return 0, err
	}

	res, err := b.raw.ChannelsGetParticipants(ctx, &tg.ChannelsGetParticipantsRequest{
		Channel: channel,
		Filter:  &tg.ChannelParticipantsRecent{},
		Limit:   1,
	})
	if err != nil {
		return 0, asAPIError(err)
	}

	participants, ok := res.AsModified()
	if !ok {
		return 0, &Error{Code: 500, Description: "Internal Server Error: unexpected participants response"}
	}

	return participants.Count, nil
}

// userToInputPeer wraps an input user as the input peer the participant-editing
// methods expect.
func userToInputPeer(u tg.InputUserClass) tg.InputPeerClass {
	switch v := u.(type) {
	case *tg.InputUser:
		return &tg.InputPeerUser{UserID: v.UserID, AccessHash: v.AccessHash}
	case *tg.InputUserFromMessage:
		return &tg.InputPeerUserFromMessage{Peer: v.Peer, MsgID: v.MsgID, UserID: v.UserID}
	case *tg.InputUserSelf:
		return &tg.InputPeerSelf{}
	default:
		return &tg.InputPeerEmpty{}
	}
}

// SetChatAdministratorCustomTitle sets a custom title (rank) for an
// administrator that the bot has promoted in a supergroup. It preserves the
// administrator's existing rights.
func (b *Bot) SetChatAdministratorCustomTitle(ctx context.Context, chat ChatID, userID int64, customTitle string) error {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	res, err := b.raw.ChannelsGetParticipant(ctx, &tg.ChannelsGetParticipantRequest{
		Channel:     channel,
		Participant: userToInputPeer(user),
	})
	if err != nil {
		return asAPIError(err)
	}

	admin, ok := res.Participant.(*tg.ChannelParticipantAdmin)
	if !ok {
		return &Error{Code: 400, Description: "Bad Request: user is not an administrator"}
	}

	if _, err := b.raw.ChannelsEditAdmin(ctx, &tg.ChannelsEditAdminRequest{
		Channel:     channel,
		UserID:      user,
		AdminRights: admin.AdminRights,
		Rank:        customTitle,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// SetChatMemberTag sets a custom tag (rank) for a member of a supergroup or
// channel. An empty tag removes it. The bot must be an administrator with the
// appropriate rights.
func (b *Bot) SetChatMemberTag(ctx context.Context, chat ChatID, userID int64, tag string) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	if _, err := b.raw.MessagesEditChatParticipantRank(ctx, &tg.MessagesEditChatParticipantRankRequest{
		Peer:        peer,
		Participant: userToInputPeer(user),
		Rank:        tag,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}
