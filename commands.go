package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// BotCommand represents a bot command shown in the menu.
type BotCommand struct {
	// Command is the text of the command, 1-32 characters, lowercase English
	// letters, digits and underscores.
	Command string `json:"command"`
	// Description is the command description, 1-256 characters.
	Description string `json:"description"`
}

// BotCommandScope is a sealed union describing the scope to which a list of bot
// commands applies.
//
// Construct with BotCommandScopeDefault, BotCommandScopeAllPrivateChats,
// BotCommandScopeAllGroupChats, BotCommandScopeAllChatAdministrators,
// BotCommandScopeChat, BotCommandScopeChatAdministrators or
// BotCommandScopeChatMember.
type BotCommandScope interface {
	isBotCommandScope()
	// resolve converts the scope into the MTProto representation. The bot needs
	// peer resolution for the chat-targeted variants, hence the receiver.
	resolve(ctx context.Context, b *Bot) (tg.BotCommandScopeClass, error)
}

type botCommandScopeDefault struct{}
type botCommandScopeAllPrivateChats struct{}
type botCommandScopeAllGroupChats struct{}
type botCommandScopeAllChatAdministrators struct{}
type botCommandScopeChat struct{ chat ChatID }
type botCommandScopeChatAdministrators struct{ chat ChatID }
type botCommandScopeChatMember struct {
	chat   ChatID
	userID int64
}

func (botCommandScopeDefault) isBotCommandScope()               {}
func (botCommandScopeAllPrivateChats) isBotCommandScope()       {}
func (botCommandScopeAllGroupChats) isBotCommandScope()         {}
func (botCommandScopeAllChatAdministrators) isBotCommandScope() {}
func (botCommandScopeChat) isBotCommandScope()                  {}
func (botCommandScopeChatAdministrators) isBotCommandScope()    {}
func (botCommandScopeChatMember) isBotCommandScope()            {}

// BotCommandScopeDefault covers all chats with no narrower scope set.
func BotCommandScopeDefault() BotCommandScope { return botCommandScopeDefault{} }

// BotCommandScopeAllPrivateChats covers all private chats.
func BotCommandScopeAllPrivateChats() BotCommandScope { return botCommandScopeAllPrivateChats{} }

// BotCommandScopeAllGroupChats covers all group and supergroup chats.
func BotCommandScopeAllGroupChats() BotCommandScope { return botCommandScopeAllGroupChats{} }

// BotCommandScopeAllChatAdministrators covers all group and supergroup chat
// administrators.
func BotCommandScopeAllChatAdministrators() BotCommandScope {
	return botCommandScopeAllChatAdministrators{}
}

// BotCommandScopeChat covers a specific chat.
func BotCommandScopeChat(chat ChatID) BotCommandScope { return botCommandScopeChat{chat: chat} }

// BotCommandScopeChatAdministrators covers the administrators of a specific
// group or supergroup chat.
func BotCommandScopeChatAdministrators(chat ChatID) BotCommandScope {
	return botCommandScopeChatAdministrators{chat: chat}
}

// BotCommandScopeChatMember covers a specific member of a group or supergroup
// chat.
func BotCommandScopeChatMember(chat ChatID, userID int64) BotCommandScope {
	return botCommandScopeChatMember{chat: chat, userID: userID}
}

func (botCommandScopeDefault) resolve(context.Context, *Bot) (tg.BotCommandScopeClass, error) {
	return &tg.BotCommandScopeDefault{}, nil
}

func (botCommandScopeAllPrivateChats) resolve(context.Context, *Bot) (tg.BotCommandScopeClass, error) {
	return &tg.BotCommandScopeUsers{}, nil
}

func (botCommandScopeAllGroupChats) resolve(context.Context, *Bot) (tg.BotCommandScopeClass, error) {
	return &tg.BotCommandScopeChats{}, nil
}

func (botCommandScopeAllChatAdministrators) resolve(context.Context, *Bot) (tg.BotCommandScopeClass, error) {
	return &tg.BotCommandScopeChatAdmins{}, nil
}

func (s botCommandScopeChat) resolve(ctx context.Context, b *Bot) (tg.BotCommandScopeClass, error) {
	peer, err := b.resolveInputPeer(ctx, s.chat)
	if err != nil {
		return nil, err
	}
	return &tg.BotCommandScopePeer{Peer: peer}, nil
}

func (s botCommandScopeChatAdministrators) resolve(ctx context.Context, b *Bot) (tg.BotCommandScopeClass, error) {
	peer, err := b.resolveInputPeer(ctx, s.chat)
	if err != nil {
		return nil, err
	}
	return &tg.BotCommandScopePeerAdmins{Peer: peer}, nil
}

func (s botCommandScopeChatMember) resolve(ctx context.Context, b *Bot) (tg.BotCommandScopeClass, error) {
	peer, err := b.resolveInputPeer(ctx, s.chat)
	if err != nil {
		return nil, err
	}
	user, err := b.resolveInputUser(ctx, s.userID)
	if err != nil {
		return nil, err
	}
	return &tg.BotCommandScopePeerUser{Peer: peer, UserID: user}, nil
}

// CommandOption configures a commands call (scope and language code).
type CommandOption func(*commandConfig)

type commandConfig struct {
	scope    BotCommandScope
	langCode string
}

func (c commandConfig) resolveScope(ctx context.Context, b *Bot) (tg.BotCommandScopeClass, error) {
	if c.scope == nil {
		return &tg.BotCommandScopeDefault{}, nil
	}
	return c.scope.resolve(ctx, b)
}

// WithCommandScope restricts the commands call to the given scope. When unset
// the default scope is used.
func WithCommandScope(scope BotCommandScope) CommandOption {
	return func(c *commandConfig) { c.scope = scope }
}

// WithLanguageCode restricts the commands call to a two-letter language code.
func WithLanguageCode(code string) CommandOption {
	return func(c *commandConfig) { c.langCode = code }
}

// SetMyCommands sets the list of the bot's commands for the given scope and
// language.
func (b *Bot) SetMyCommands(ctx context.Context, commands []BotCommand, opts ...CommandOption) error {
	var cfg commandConfig
	for _, o := range opts {
		o(&cfg)
	}
	scope, err := cfg.resolveScope(ctx, b)
	if err != nil {
		return err
	}

	cmds := make([]tg.BotCommand, len(commands))
	for i, c := range commands {
		cmds[i] = tg.BotCommand{Command: c.Command, Description: c.Description}
	}

	if _, err := b.raw.BotsSetBotCommands(ctx, &tg.BotsSetBotCommandsRequest{
		Scope:    scope,
		LangCode: cfg.langCode,
		Commands: cmds,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// GetMyCommands returns the current list of the bot's commands for the given
// scope and language.
func (b *Bot) GetMyCommands(ctx context.Context, opts ...CommandOption) ([]BotCommand, error) {
	var cfg commandConfig
	for _, o := range opts {
		o(&cfg)
	}
	scope, err := cfg.resolveScope(ctx, b)
	if err != nil {
		return nil, err
	}

	cmds, err := b.raw.BotsGetBotCommands(ctx, &tg.BotsGetBotCommandsRequest{
		Scope:    scope,
		LangCode: cfg.langCode,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	out := make([]BotCommand, len(cmds))
	for i, c := range cmds {
		out[i] = BotCommand{Command: c.Command, Description: c.Description}
	}
	return out, nil
}

// DeleteMyCommands clears the list of the bot's commands for the given scope and
// language, restoring the default commands.
func (b *Bot) DeleteMyCommands(ctx context.Context, opts ...CommandOption) error {
	var cfg commandConfig
	for _, o := range opts {
		o(&cfg)
	}
	scope, err := cfg.resolveScope(ctx, b)
	if err != nil {
		return err
	}

	if _, err := b.raw.BotsResetBotCommands(ctx, &tg.BotsResetBotCommandsRequest{
		Scope:    scope,
		LangCode: cfg.langCode,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}
