package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// BotInfoOption configures a bot-profile call (currently the language code the
// value applies to). An empty language code targets the default locale.
type BotInfoOption func(*botInfoConfig)

type botInfoConfig struct {
	langCode string
}

// WithBotInfoLanguage restricts the call to a two-letter IETF language code. The
// empty code (default) sets the value shown to users with no localized value.
func WithBotInfoLanguage(code string) BotInfoOption {
	return func(c *botInfoConfig) { c.langCode = code }
}

func newBotInfoConfig(opts []BotInfoOption) botInfoConfig {
	var cfg botInfoConfig

	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

// setBotInfo updates a single bot-info field. The MTProto bots.setBotInfo edits
// the bot identified by InputUserSelf; only the populated flag fields change.
func (b *Bot) setBotInfo(ctx context.Context, langCode string, mut func(*tg.BotsSetBotInfoRequest)) error {
	req := &tg.BotsSetBotInfoRequest{
		Bot:      &tg.InputUserSelf{},
		LangCode: langCode,
	}
	mut(req)

	if _, err := b.raw.BotsSetBotInfo(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// SetMyName changes the bot's name for the given language. An empty name clears
// the localized value, falling back to the default.
func (b *Bot) SetMyName(ctx context.Context, name string, opts ...BotInfoOption) error {
	cfg := newBotInfoConfig(opts)

	return b.setBotInfo(ctx, cfg.langCode, func(r *tg.BotsSetBotInfoRequest) { r.SetName(name) })
}

// SetMyDescription changes the bot's description, shown in an empty chat with the
// bot, for the given language.
func (b *Bot) SetMyDescription(ctx context.Context, description string, opts ...BotInfoOption) error {
	cfg := newBotInfoConfig(opts)

	return b.setBotInfo(ctx, cfg.langCode, func(r *tg.BotsSetBotInfoRequest) { r.SetDescription(description) })
}

// SetMyShortDescription changes the bot's short description, shown on the bot's
// profile page and in the chat list, for the given language.
func (b *Bot) SetMyShortDescription(ctx context.Context, shortDescription string, opts ...BotInfoOption) error {
	cfg := newBotInfoConfig(opts)

	return b.setBotInfo(ctx, cfg.langCode, func(r *tg.BotsSetBotInfoRequest) { r.SetAbout(shortDescription) })
}

// getBotInfo fetches the bot's localized info via bots.getBotInfo.
func (b *Bot) getBotInfo(ctx context.Context, langCode string) (*tg.BotsBotInfo, error) {
	info, err := b.raw.BotsGetBotInfo(ctx, &tg.BotsGetBotInfoRequest{
		Bot:      &tg.InputUserSelf{},
		LangCode: langCode,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	return info, nil
}

// GetMyName returns the bot's current name for the given language.
func (b *Bot) GetMyName(ctx context.Context, opts ...BotInfoOption) (string, error) {
	cfg := newBotInfoConfig(opts)

	info, err := b.getBotInfo(ctx, cfg.langCode)
	if err != nil {
		return "", err
	}

	return info.Name, nil
}

// GetMyDescription returns the bot's current description for the given language.
func (b *Bot) GetMyDescription(ctx context.Context, opts ...BotInfoOption) (string, error) {
	cfg := newBotInfoConfig(opts)

	info, err := b.getBotInfo(ctx, cfg.langCode)
	if err != nil {
		return "", err
	}

	return info.Description, nil
}

// GetMyShortDescription returns the bot's current short description for the given
// language.
func (b *Bot) GetMyShortDescription(ctx context.Context, opts ...BotInfoOption) (string, error) {
	cfg := newBotInfoConfig(opts)

	info, err := b.getBotInfo(ctx, cfg.langCode)
	if err != nil {
		return "", err
	}

	return info.About, nil
}

// ChatMenuButtonOption configures a chat-menu-button call.
type ChatMenuButtonOption func(*menuButtonConfig)

type menuButtonConfig struct {
	userID    int64
	hasUserID bool
}

// WithMenuButtonChat targets the menu button of a single private chat with the
// given user instead of the bot-wide default.
func WithMenuButtonChat(userID int64) ChatMenuButtonOption {
	return func(c *menuButtonConfig) {
		c.userID = userID
		c.hasUserID = true
	}
}

// menuButtonUser resolves the menu-button target. The bot-wide default uses
// InputUserEmpty; a per-chat button resolves the user.
func (b *Bot) menuButtonUser(ctx context.Context, cfg menuButtonConfig) (tg.InputUserClass, error) {
	if !cfg.hasUserID {
		return &tg.InputUserEmpty{}, nil
	}

	return b.resolveInputUser(ctx, cfg.userID)
}

// menuButtonToTg converts a Bot API menu button to the MTProto representation.
// A nil button is treated as the default.
//
// The switch over the sealed MenuButton union is exhaustive (gochecksumtype).
func menuButtonToTg(button MenuButton) tg.BotMenuButtonClass {
	switch v := button.(type) {
	case MenuButtonCommands:
		return &tg.BotMenuButtonCommands{}
	case MenuButtonWebApp:
		return &tg.BotMenuButton{Text: v.Text, URL: v.WebApp.URL}
	case MenuButtonDefault:
		return &tg.BotMenuButtonDefault{}
	default:
		return &tg.BotMenuButtonDefault{}
	}
}

// menuButtonFromTg converts an MTProto menu button to the Bot API union.
func menuButtonFromTg(button tg.BotMenuButtonClass) MenuButton {
	switch v := button.(type) {
	case *tg.BotMenuButtonCommands:
		return MenuButtonCommands{Type: MenuButtonCommandsType}
	case *tg.BotMenuButton:
		return MenuButtonWebApp{
			Type:   MenuButtonWebAppType,
			Text:   v.Text,
			WebApp: WebAppInfo{URL: v.URL},
		}
	default:
		return MenuButtonDefault{Type: MenuButtonDefaultType}
	}
}

// SetChatMenuButton changes the bot's menu button in a private chat, or the
// default menu button when no chat is targeted.
func (b *Bot) SetChatMenuButton(ctx context.Context, button MenuButton, opts ...ChatMenuButtonOption) error {
	var cfg menuButtonConfig

	for _, o := range opts {
		o(&cfg)
	}

	user, err := b.menuButtonUser(ctx, cfg)
	if err != nil {
		return err
	}

	if _, err := b.raw.BotsSetBotMenuButton(ctx, &tg.BotsSetBotMenuButtonRequest{
		UserID: user,
		Button: menuButtonToTg(button),
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// GetChatMenuButton returns the current menu button of a private chat, or the
// default menu button when no chat is targeted.
func (b *Bot) GetChatMenuButton(ctx context.Context, opts ...ChatMenuButtonOption) (MenuButton, error) {
	var cfg menuButtonConfig

	for _, o := range opts {
		o(&cfg)
	}

	user, err := b.menuButtonUser(ctx, cfg)
	if err != nil {
		return nil, err
	}

	button, err := b.raw.BotsGetBotMenuButton(ctx, user)
	if err != nil {
		return nil, asAPIError(err)
	}

	return menuButtonFromTg(button), nil
}

// SetMyDefaultAdministratorRights changes the default administrator rights
// requested by the bot when it is added to a chat. forChannels selects the
// rights used for channels; otherwise the rights for groups and supergroups are
// changed. Pass a zero ChatAdminRights to clear the defaults.
func (b *Bot) SetMyDefaultAdministratorRights(ctx context.Context, rights ChatAdminRights, forChannels bool) error {
	if forChannels {
		if _, err := b.raw.BotsSetBotBroadcastDefaultAdminRights(ctx, rights.toTg()); err != nil {
			return asAPIError(err)
		}

		return nil
	}

	if _, err := b.raw.BotsSetBotGroupDefaultAdminRights(ctx, rights.toTg()); err != nil {
		return asAPIError(err)
	}

	return nil
}

// GetMyDefaultAdministratorRights returns the bot's current default
// administrator rights. forChannels selects the channel rights; otherwise the
// group rights are returned.
func (b *Bot) GetMyDefaultAdministratorRights(ctx context.Context, forChannels bool) (ChatAdminRights, error) {
	full, err := b.raw.UsersGetFullUser(ctx, &tg.InputUserSelf{})
	if err != nil {
		return ChatAdminRights{}, asAPIError(err)
	}

	if forChannels {
		return chatAdminRightsFromTg(full.FullUser.BotBroadcastAdminRights), nil
	}

	return chatAdminRightsFromTg(full.FullUser.BotGroupAdminRights), nil
}

// chatAdminRightsFromTg converts MTProto admin rights to the Bot API
// representation. It is the inverse of ChatAdminRights.toTg.
func chatAdminRightsFromTg(r tg.ChatAdminRights) ChatAdminRights {
	return ChatAdminRights{
		IsAnonymous:         r.Anonymous,
		CanManageChat:       r.Other,
		CanDeleteMessages:   r.DeleteMessages,
		CanManageVideoChats: r.ManageCall,
		CanRestrictMembers:  r.BanUsers,
		CanPromoteMembers:   r.AddAdmins,
		CanChangeInfo:       r.ChangeInfo,
		CanInviteUsers:      r.InviteUsers,
		CanPostMessages:     r.PostMessages,
		CanEditMessages:     r.EditMessages,
		CanPinMessages:      r.PinMessages,
		CanManageTopics:     r.ManageTopics,
	}
}
