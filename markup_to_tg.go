package botapi

import (
	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"
)

// replyMarkupToTg translates a Bot API ReplyMarkup into the MTProto
// tg.ReplyMarkupClass understood by the sender.
//
// The switch over the sealed ReplyMarkup union is exhaustive (enforced by the
// gochecksumtype linter): every variant is handled.
func replyMarkupToTg(m ReplyMarkup) (tg.ReplyMarkupClass, error) {
	switch m := m.(type) {
	case *InlineKeyboardMarkup:
		rows := make([]tg.KeyboardButtonRow, 0, len(m.InlineKeyboard))
		for _, row := range m.InlineKeyboard {
			buttons := make([]tg.KeyboardButtonClass, len(row))
			for i, btn := range row {
				b, err := inlineButtonToTg(btn)
				if err != nil {
					return nil, err
				}

				buttons[i] = b
			}

			rows = append(rows, tg.KeyboardButtonRow{Buttons: buttons})
		}

		return &tg.ReplyInlineMarkup{Rows: rows}, nil
	case *ReplyKeyboardMarkup:
		rows := make([]tg.KeyboardButtonRow, 0, len(m.Keyboard))
		for _, row := range m.Keyboard {
			buttons := make([]tg.KeyboardButtonClass, len(row))
			for i, btn := range row {
				buttons[i] = keyboardButtonToTg(btn)
			}

			rows = append(rows, tg.KeyboardButtonRow{Buttons: buttons})
		}

		res := &tg.ReplyKeyboardMarkup{
			Resize:     m.ResizeKeyboard,
			SingleUse:  m.OneTimeKeyboard,
			Selective:  m.Selective,
			Persistent: m.IsPersistent,
			Rows:       rows,
		}
		if m.InputFieldPlaceholder != "" {
			res.SetPlaceholder(m.InputFieldPlaceholder)
		}

		return res, nil
	case *ReplyKeyboardRemove:
		if m.Selective {
			return markup.SelectiveHide(), nil
		}

		return markup.Hide(), nil
	case *ForceReply:
		res := &tg.ReplyKeyboardForceReply{Selective: m.Selective}
		if m.InputFieldPlaceholder != "" {
			res.SetPlaceholder(m.InputFieldPlaceholder)
		}

		return res, nil
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: unsupported reply markup"}
	}
}

func inlineButtonToTg(btn InlineKeyboardButton) (tg.KeyboardButtonClass, error) {
	switch {
	case btn.URL != "":
		return markup.URL(btn.Text, btn.URL), nil
	case btn.CallbackData != "":
		return markup.Callback(btn.Text, []byte(btn.CallbackData)), nil
	case btn.WebApp != nil:
		return markup.WebView(btn.Text, btn.WebApp.URL), nil
	case btn.SwitchInlineQuery != nil:
		return markup.SwitchInline(btn.Text, *btn.SwitchInlineQuery, false), nil
	case btn.SwitchInlineQueryCurrentChat != nil:
		return markup.SwitchInline(btn.Text, *btn.SwitchInlineQueryCurrentChat, true), nil
	case btn.Pay:
		return markup.Buy(btn.Text), nil
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: text buttons are unallowed in the inline keyboard"}
	}
}

func keyboardButtonToTg(btn KeyboardButton) tg.KeyboardButtonClass {
	switch {
	case btn.RequestContact:
		return markup.RequestPhone(btn.Text)
	case btn.RequestLocation:
		return markup.RequestGeoLocation(btn.Text)
	case btn.RequestPoll != nil:
		return markup.RequestPoll(btn.Text, btn.RequestPoll.Type == PollQuiz)
	case btn.WebApp != nil:
		return markup.SimpleWebView(btn.Text, btn.WebApp.URL)
	default:
		return markup.Button(btn.Text)
	}
}
