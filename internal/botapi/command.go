package botapi

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func (b *BotAPI) convertFromBotCommandScopeClass(
	ctx context.Context,
	scope *oas.BotCommandScope,
) (tg.BotCommandScopeClass, error) {
	if scope == nil {
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
			return nil, errors.Errorf("resolve chatID: %w", err)
		}
		return &tg.BotCommandScopePeer{Peer: p}, nil
	case oas.BotCommandScopeChatAdministratorsBotCommandScope:
		chatID := scope.BotCommandScopeChatAdministrators.ChatID
		p, err := b.resolveID(ctx, chatID)
		if err != nil {
			return nil, errors.Errorf("resolve chatID: %w", err)
		}
		return &tg.BotCommandScopePeerAdmins{Peer: p}, nil
	case oas.BotCommandScopeChatMemberBotCommandScope:
		userID := scope.BotCommandScopeChatMember.UserID
		user, err := b.resolveUserID(ctx, userID)
		if err != nil {
			return nil, errors.Errorf("resolve userID: %w", err)
		}

		chatID := scope.BotCommandScopeChatMember.ChatID
		p, err := b.resolveID(ctx, chatID)
		if err != nil {
			return nil, errors.Errorf("resolve chatID: %w", err)
		}
		return &tg.BotCommandScopePeerUser{
			Peer:   p,
			UserID: user,
		}, nil
	default:
		return nil, errors.Errorf("unknown peer type %q", scope.Type)
	}
}

// GetMyCommands implements oas.Handler.
func (b *BotAPI) GetMyCommands(ctx context.Context, req oas.GetMyCommands) (oas.ResultArrayOfBotCommand, error) {
	scope, err := b.convertFromBotCommandScopeClass(ctx, req.Scope)
	if err != nil {
		return oas.ResultArrayOfBotCommand{}, errors.Errorf("convert scope: %w", err)
	}

	cmds, err := b.client.API().BotsGetBotCommands(ctx, &tg.BotsGetBotCommandsRequest{
		Scope:    scope,
		LangCode: req.LanguageCode.Value,
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
	scope, err := b.convertFromBotCommandScopeClass(ctx, req.Scope)
	if err != nil {
		return oas.Result{}, errors.Errorf("convert scope: %w", err)
	}

	commands := make([]tg.BotCommand, len(req.Commands))
	for i, cmd := range req.Commands {
		commands[i] = tg.BotCommand{
			Command:     cmd.Command,
			Description: cmd.Description,
		}
	}

	r, err := b.client.API().BotsSetBotCommands(ctx, &tg.BotsSetBotCommandsRequest{
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
func (b *BotAPI) DeleteMyCommands(ctx context.Context, req oas.DeleteMyCommands) (oas.Result, error) {
	scope, err := b.convertFromBotCommandScopeClass(ctx, req.Scope)
	if err != nil {
		return oas.Result{}, errors.Errorf("convert scope: %w", err)
	}

	r, err := b.client.API().BotsResetBotCommands(ctx, &tg.BotsResetBotCommandsRequest{
		Scope:    scope,
		LangCode: req.LanguageCode.Value,
	})
	if err != nil {
		return oas.Result{}, err
	}

	return resultOK(r), nil
}
