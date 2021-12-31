package botapi

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func (b *BotAPI) convertToBotCommandScopeClass(
	ctx context.Context,
	input oas.OptBotCommandScope,
) (tg.BotCommandScopeClass, error) {
	scope, ok := input.Get()
	if !ok {
		return &tg.BotCommandScopeDefault{}, nil
	}
	switch scope.Type {
	case oas.BotCommandScopeDefaultBotCommandScope:
		return &tg.BotCommandScopeDefault{}, nil
	case oas.BotCommandScopeAllPrivateChatsBotCommandScope:
		return &tg.BotCommandScopeUsers{}, nil
	case oas.BotCommandScopeAllGroupChatsBotCommandScope:
		return &tg.BotCommandScopeChats{}, nil
	case oas.BotCommandScopeAllChatAdministratorsBotCommandScope:
		return &tg.BotCommandScopeChatAdmins{}, nil
	case oas.BotCommandScopeChatBotCommandScope:
		chatID := scope.BotCommandScopeChat.ChatID
		p, err := b.resolveID(ctx, chatID)
		if err != nil {
			return nil, errors.Wrap(err, "resolve chatID")
		}
		return &tg.BotCommandScopePeer{Peer: p.InputPeer()}, nil
	case oas.BotCommandScopeChatAdministratorsBotCommandScope:
		chatID := scope.BotCommandScopeChatAdministrators.ChatID
		p, err := b.resolveID(ctx, chatID)
		if err != nil {
			return nil, errors.Wrap(err, "resolve chatID")
		}
		return &tg.BotCommandScopePeerAdmins{Peer: p.InputPeer()}, nil
	case oas.BotCommandScopeChatMemberBotCommandScope:
		userID := scope.BotCommandScopeChatMember.UserID
		user, err := b.resolveUserID(ctx, userID)
		if err != nil {
			return nil, errors.Wrap(err, "resolve userID")
		}

		chatID := scope.BotCommandScopeChatMember.ChatID
		p, err := b.resolveID(ctx, chatID)
		if err != nil {
			return nil, errors.Wrap(err, "resolve chatID")
		}
		return &tg.BotCommandScopePeerUser{
			Peer:   p.InputPeer(),
			UserID: user.InputUser(),
		}, nil
	default:
		return nil, errors.Errorf("unknown peer type %q", scope.Type)
	}
}

// GetMyCommands implements oas.Handler.
func (b *BotAPI) GetMyCommands(ctx context.Context, req oas.OptGetMyCommands) (oas.ResultArrayOfBotCommand, error) {
	var (
		scope    tg.BotCommandScopeClass = &tg.BotCommandScopeDefault{}
		langCode string
	)

	if input, ok := req.Get(); ok {
		s, err := b.convertToBotCommandScopeClass(ctx, input.Scope)
		if err != nil {
			return oas.ResultArrayOfBotCommand{}, errors.Wrap(err, "convert scope")
		}
		scope = s
		langCode = input.LanguageCode.Value
	}

	cmds, err := b.raw.BotsGetBotCommands(ctx, &tg.BotsGetBotCommandsRequest{
		Scope:    scope,
		LangCode: langCode,
	})
	if err != nil {
		return oas.ResultArrayOfBotCommand{}, err
	}

	r := make([]oas.BotCommand, len(cmds))
	for i, cmd := range cmds {
		r[i] = oas.BotCommand{
			Command:     cmd.Command,
			Description: cmd.Description,
		}
	}

	return oas.ResultArrayOfBotCommand{
		Result: r,
		Ok:     true,
	}, nil
}

// SetMyCommands implements oas.Handler.
func (b *BotAPI) SetMyCommands(ctx context.Context, req oas.SetMyCommands) (oas.Result, error) {
	scope, err := b.convertToBotCommandScopeClass(ctx, req.Scope)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "convert scope")
	}

	commands := make([]tg.BotCommand, len(req.Commands))
	for i, cmd := range req.Commands {
		commands[i] = tg.BotCommand{
			Command:     cmd.Command,
			Description: cmd.Description,
		}
	}

	r, err := b.raw.BotsSetBotCommands(ctx, &tg.BotsSetBotCommandsRequest{
		Scope:    scope,
		LangCode: req.LanguageCode.Value,
		Commands: commands,
	})
	if err != nil {
		return oas.Result{}, err
	}

	return resultOK(r), nil
}

// DeleteMyCommands implements oas.Handler.
func (b *BotAPI) DeleteMyCommands(ctx context.Context, req oas.OptDeleteMyCommands) (oas.Result, error) {
	var (
		scope    tg.BotCommandScopeClass = &tg.BotCommandScopeDefault{}
		langCode string
	)

	if input, ok := req.Get(); ok {
		s, err := b.convertToBotCommandScopeClass(ctx, input.Scope)
		if err != nil {
			return oas.Result{}, errors.Wrap(err, "convert scope")
		}
		scope = s
		langCode = input.LanguageCode.Value
	}

	r, err := b.raw.BotsResetBotCommands(ctx, &tg.BotsResetBotCommandsRequest{
		Scope:    scope,
		LangCode: langCode,
	})
	if err != nil {
		return oas.Result{}, err
	}

	return resultOK(r), nil
}
