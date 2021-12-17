package botapi

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram/peers"

	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func (b *BotAPI) convertToTelegramInlineButton(
	ctx context.Context,
	button oas.InlineKeyboardButton,
) (tg.KeyboardButtonClass, error) {
	switch {
	case button.URL.Set:
		return markup.URL(button.Text, button.URL.Value), nil
	case button.CallbackData.Set:
		return markup.Callback(button.Text, []byte(button.CallbackData.Value)), nil
	case button.CallbackGame != nil:
		return markup.Game(button.Text), nil
	case button.Pay.Value:
		return markup.Buy(button.Text), nil
	case button.SwitchInlineQuery.Set:
		return markup.SwitchInline(button.Text, button.SwitchInlineQuery.Value, false), nil
	case button.SwitchInlineQueryCurrentChat.Set:
		return markup.SwitchInline(button.Text, button.SwitchInlineQuery.Value, true), nil
	case button.LoginURL.Set:
		loginURL := button.LoginURL.Value

		var user tg.InputUserClass = &tg.InputUserSelf{}
		if v, ok := loginURL.BotUsername.Get(); ok && v != "" {
			p, err := b.peers.ResolveDomain(ctx, loginURL.BotUsername.Value)
			if err != nil {
				return nil, errors.Wrap(err, "resolve bot")
			}

			u, ok := p.(peers.User)
			if !ok {
				return nil, &BadRequestError{Message: "given username is not bot"}
			}
			user = u.InputUser()
		}

		return &tg.InputKeyboardButtonURLAuth{
			RequestWriteAccess: loginURL.RequestWriteAccess.Value,
			Text:               button.Text,
			FwdText:            loginURL.ForwardText.Value,
			URL:                loginURL.URL,
			Bot:                user,
		}, nil
	default:
		return nil, &BadRequestError{Message: "text buttons are unallowed in the inline keyboard"}
	}
}

func (b *BotAPI) convertToTelegramButton(kb oas.KeyboardButton) tg.KeyboardButtonClass {
	if text, ok := kb.GetString(); ok {
		return markup.Button(text)
	}

	button := kb.KeyboardButtonObject
	if button.RequestLocation.Value || button.RequestContact.Value {
		return markup.RequestPhone(button.Text)
	}

	if poll, ok := button.RequestPoll.Get(); ok {
		return markup.RequestPoll(button.Text, poll.Type.Value == "quiz")
	}

	return markup.Button(button.Text)
}

func (b *BotAPI) convertToTelegramReplyMarkup(
	ctx context.Context,
	m *oas.SendMessageReplyMarkup,
) (tg.ReplyMarkupClass, error) {
	switch m.Type {
	case oas.InlineKeyboardMarkupSendMessageReplyMarkup:
		rows := m.InlineKeyboardMarkup.InlineKeyboard
		result := &tg.ReplyInlineMarkup{Rows: make([]tg.KeyboardButtonRow, 0, len(rows))}
		for _, row := range rows {
			resultRow := make([]tg.KeyboardButtonClass, len(row))
			for i, button := range row {
				resultButton, err := b.convertToTelegramInlineButton(ctx, button)
				if err != nil {
					return nil, errors.Wrapf(err, "convert button %d", i)
				}
				resultRow[i] = resultButton
			}
			result.Rows = append(result.Rows, tg.KeyboardButtonRow{Buttons: resultRow})
		}
		return result, nil
	case oas.ReplyKeyboardMarkupSendMessageReplyMarkup:
		mark := m.ReplyKeyboardMarkup
		rows := mark.Keyboard

		result := &tg.ReplyKeyboardMarkup{
			Resize:    mark.ResizeKeyboard.Value,
			SingleUse: mark.OneTimeKeyboard.Value,
			Selective: mark.Selective.Value,
			Rows:      make([]tg.KeyboardButtonRow, 0, len(rows)),
		}
		if v, ok := mark.InputFieldPlaceholder.Get(); ok {
			result.SetPlaceholder(v)
		}
		for _, row := range rows {
			resultRow := make([]tg.KeyboardButtonClass, len(row))
			for _, button := range row {
				resultRow = append(resultRow, b.convertToTelegramButton(button))
			}
			result.Rows = append(result.Rows, tg.KeyboardButtonRow{Buttons: resultRow})
		}
		return result, nil
	case oas.ReplyKeyboardRemoveSendMessageReplyMarkup:
		if v, ok := m.ReplyKeyboardRemove.Selective.Get(); ok && v {
			return markup.SelectiveHide(), nil
		}
		return markup.Hide(), nil
	case oas.ForceReplySendMessageReplyMarkup:
		mark := m.ForceReply
		result := &tg.ReplyKeyboardForceReply{
			Selective:   mark.Selective.Value,
			Placeholder: mark.InputFieldPlaceholder.Value,
		}
		return result, nil
	default:
		return nil, errors.Errorf("unknown type %q", m.Type)
	}
}

func convertToBotAPIInlineReplyMarkup(mkp *tg.ReplyInlineMarkup) oas.InlineKeyboardMarkup {
	resultRows := make([][]oas.InlineKeyboardButton, len(mkp.Rows))
	for i, row := range mkp.Rows {
		resultRow := make([]oas.InlineKeyboardButton, len(row.Buttons))
		for i, b := range row.Buttons {
			button := oas.InlineKeyboardButton{Text: b.GetText()}
			switch b := b.(type) {
			case *tg.KeyboardButtonURL:
				button.URL.SetTo(b.URL)
			case *tg.KeyboardButtonCallback:
				button.CallbackData.SetTo(string(b.Data))
			case *tg.KeyboardButtonSwitchInline:
				if b.SamePeer {
					button.SwitchInlineQueryCurrentChat.SetTo(b.Query)
				} else {
					button.SwitchInlineQuery.SetTo(b.Query)
				}
			case *tg.KeyboardButtonGame:
				button.CallbackGame = new(oas.CallbackGame)
			case *tg.KeyboardButtonBuy:
				button.Pay.SetTo(true)
			case *tg.KeyboardButtonURLAuth:
				// Quote: login_url buttons are represented as ordinary url buttons.
				//
				// See Message definition
				// See https://github.com/tdlib/telegram-bot-api/blob/90f52477814a2d8a08c9ffb1d780fd179815d715/telegram-bot-api/Client.cpp#L1526
				button.URL.SetTo(b.URL)
			}
			resultRow[i] = button
		}
		resultRows[i] = resultRow
	}

	return oas.InlineKeyboardMarkup{
		InlineKeyboard: resultRows,
	}
}
