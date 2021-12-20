package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func TestBotAPI_convertToTelegramInlineButton(t *testing.T) {
	tests := []struct {
		name    string
		input   oas.InlineKeyboardButton
		want    tg.KeyboardButtonClass
		wantErr bool
	}{
		{
			"URL",
			oas.InlineKeyboardButton{
				Text: "aboba",
				URL:  oas.NewOptString("https://ya.ru"),
			},
			markup.URL("aboba", "https://ya.ru"),
			false,
		},
		{
			"Callback",
			oas.InlineKeyboardButton{
				Text:         "aboba",
				CallbackData: oas.NewOptString("data"),
			},
			markup.Callback("aboba", []byte("data")),
			false,
		},
		{
			"Game",
			oas.InlineKeyboardButton{
				Text:         "aboba",
				CallbackGame: &oas.CallbackGame{},
			},
			markup.Game("aboba"),
			false,
		},
		{
			"Pay",
			oas.InlineKeyboardButton{
				Text: "aboba",
				Pay:  oas.NewOptBool(true),
			},
			markup.Buy("aboba"),
			false,
		},
		{
			"SwitchInlineQuery",
			oas.InlineKeyboardButton{
				Text:              "aboba",
				SwitchInlineQuery: oas.NewOptString("query"),
			},
			markup.SwitchInline("aboba", "query", false),
			false,
		},
		{
			"SwitchInlineQueryCurrentChat",
			oas.InlineKeyboardButton{
				Text:                         "aboba",
				SwitchInlineQueryCurrentChat: oas.NewOptString("query"),
			},
			markup.SwitchInline("aboba", "query", true),
			false,
		},
		{
			"LoginURL",
			oas.InlineKeyboardButton{
				Text: "aboba",
				LoginURL: oas.NewOptLoginUrl(oas.LoginUrl{
					URL:                "https://ya.ru",
					ForwardText:        oas.NewOptString("forward text"),
					BotUsername:        oas.OptString{},
					RequestWriteAccess: oas.NewOptBool(true),
				}),
			},
			&tg.InputKeyboardButtonURLAuth{
				RequestWriteAccess: true,
				Text:               "aboba",
				FwdText:            "forward text",
				URL:                "https://ya.ru",
				Bot:                &tg.InputUserSelf{},
			},
			false,
		},
		{
			"Text",
			oas.InlineKeyboardButton{
				Text: "aboba",
			},
			nil,
			true,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
				r, err := api.convertToTelegramInlineButton(ctx, tt.input)
				if tt.wantErr {
					a.Error(err)
					return
				}
				a.NoError(err)
				a.Equal(tt.want, r)
			})
		})
	}
}

func Test_convertToBotAPIInlineButton(t *testing.T) {
	tests := []struct {
		name  string
		input tg.KeyboardButtonClass
		want  oas.InlineKeyboardButton
	}{
		{
			"URL",
			&tg.KeyboardButtonURL{
				Text: "aboba",
				URL:  "https://ya.ru",
			},
			oas.InlineKeyboardButton{
				Text: "aboba",
				URL:  oas.NewOptString("https://ya.ru"),
			},
		},
		{
			"Callback",
			&tg.KeyboardButtonCallback{
				Text: "aboba",
				Data: []byte("data"),
			},
			oas.InlineKeyboardButton{
				Text:         "aboba",
				CallbackData: oas.NewOptString("data"),
			},
		},
		{
			"Game",
			&tg.KeyboardButtonGame{
				Text: "aboba",
			},
			oas.InlineKeyboardButton{
				Text:         "aboba",
				CallbackGame: &oas.CallbackGame{},
			},
		},
		{
			"Game",
			&tg.KeyboardButtonBuy{
				Text: "aboba",
			},
			oas.InlineKeyboardButton{
				Text: "aboba",
				Pay:  oas.NewOptBool(true),
			},
		},
		{
			"Pay",
			&tg.KeyboardButtonBuy{
				Text: "aboba",
			},
			oas.InlineKeyboardButton{
				Text: "aboba",
				Pay:  oas.NewOptBool(true),
			},
		},
		{
			"SwitchInlineQuery",
			markup.SwitchInline("aboba", "query", false),
			oas.InlineKeyboardButton{
				Text:              "aboba",
				SwitchInlineQuery: oas.NewOptString("query"),
			},
		},
		{
			"SwitchInlineQueryCurrentChat",
			markup.SwitchInline("aboba", "query", true),
			oas.InlineKeyboardButton{
				Text:                         "aboba",
				SwitchInlineQueryCurrentChat: oas.NewOptString("query"),
			},
		},
		{
			"LoginURL",
			&tg.KeyboardButtonURLAuth{
				Text:    "aboba",
				FwdText: "forward text",
				URL:     "https://ya.ru",
			},
			oas.InlineKeyboardButton{
				Text: "aboba",
				URL:  oas.NewOptString("https://ya.ru"),
			},
		},
	}

	var (
		inputs  []tg.KeyboardButtonClass
		results []oas.InlineKeyboardButton
	)
	for _, tt := range tests {
		inputs = append(inputs, tt.input)
		results = append(results, tt.want)
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, convertToBotAPIInlineButton(tt.input))
		})
	}

	m := convertToBotAPIInlineReplyMarkup(&tg.ReplyInlineMarkup{
		Rows: []tg.KeyboardButtonRow{
			{Buttons: inputs},
			{Buttons: inputs[:1]},
		},
	})
	a := require.New(t)
	a.Len(m.InlineKeyboard, 2)
	a.Equal(results, m.InlineKeyboard[0])
	a.Equal(results[:1], m.InlineKeyboard[1])
}

func Test_convertToTelegramButton(t *testing.T) {
	obj := oas.NewKeyboardButtonObjectKeyboardButton
	tests := []struct {
		name  string
		input oas.KeyboardButton
		want  tg.KeyboardButtonClass
	}{
		{
			"StringText",
			oas.NewStringKeyboardButton("aboba"),
			&tg.KeyboardButton{Text: "aboba"},
		},
		{
			"Text",
			obj(oas.KeyboardButtonObject{
				Text: "aboba",
			}),
			&tg.KeyboardButton{Text: "aboba"},
		},
		{
			"RequestLocation",
			obj(oas.KeyboardButtonObject{
				Text:            "aboba",
				RequestLocation: oas.NewOptBool(true),
			}),
			&tg.KeyboardButtonRequestGeoLocation{Text: "aboba"},
		},
		{
			"RequestContact",
			obj(oas.KeyboardButtonObject{
				Text:           "aboba",
				RequestContact: oas.NewOptBool(true),
			}),
			&tg.KeyboardButtonRequestPhone{Text: "aboba"},
		},
		{
			"RequestPoll",
			obj(oas.KeyboardButtonObject{
				Text: "aboba",
				RequestPoll: oas.NewOptKeyboardButtonPollType(oas.KeyboardButtonPollType{
					Type: oas.NewOptString("quiz"),
				}),
			}),
			&tg.KeyboardButtonRequestPoll{Quiz: true, Text: "aboba"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, convertToTelegramButton(tt.input))
		})
	}
}

func TestBotAPI_convertToTelegramReplyMarkup(t *testing.T) {
	tests := []struct {
		name    string
		input   *oas.SendMessageReplyMarkup
		want    tg.ReplyMarkupClass
		wantErr bool
	}{
		{
			"Inline",
			&oas.SendMessageReplyMarkup{
				Type: oas.InlineKeyboardMarkupSendMessageReplyMarkup,
				InlineKeyboardMarkup: oas.InlineKeyboardMarkup{
					InlineKeyboard: [][]oas.InlineKeyboardButton{
						{
							{
								Text:         "aboba",
								CallbackData: oas.NewOptString("data"),
							},
						},
					},
				},
			},
			&tg.ReplyInlineMarkup{
				Rows: []tg.KeyboardButtonRow{
					{
						Buttons: []tg.KeyboardButtonClass{
							markup.Callback("aboba", []byte("data")),
						},
					},
				},
			},
			false,
		},
		{
			"Reply",
			&oas.SendMessageReplyMarkup{
				Type: oas.ReplyKeyboardMarkupSendMessageReplyMarkup,
				ReplyKeyboardMarkup: oas.ReplyKeyboardMarkup{
					Keyboard: [][]oas.KeyboardButton{
						{
							oas.NewStringKeyboardButton("aboba"),
						},
					},
					ResizeKeyboard:        oas.NewOptBool(true),
					OneTimeKeyboard:       oas.NewOptBool(true),
					InputFieldPlaceholder: oas.NewOptString("placeholder"),
					Selective:             oas.NewOptBool(true),
				},
			},
			&tg.ReplyKeyboardMarkup{
				Resize:    true,
				SingleUse: true,
				Selective: true,
				Rows: []tg.KeyboardButtonRow{
					{
						Buttons: []tg.KeyboardButtonClass{
							markup.Button("aboba"),
						},
					},
				},
				Placeholder: "placeholder",
			},
			false,
		},
		{
			"Hide",
			&oas.SendMessageReplyMarkup{
				Type: oas.ReplyKeyboardRemoveSendMessageReplyMarkup,
				ReplyKeyboardRemove: oas.ReplyKeyboardRemove{
					RemoveKeyboard: true,
				},
			},
			&tg.ReplyKeyboardHide{
				Selective: false,
			},
			false,
		},
		{
			"SelectiveHide",
			&oas.SendMessageReplyMarkup{
				Type: oas.ReplyKeyboardRemoveSendMessageReplyMarkup,
				ReplyKeyboardRemove: oas.ReplyKeyboardRemove{
					RemoveKeyboard: true,
					Selective:      oas.NewOptBool(true),
				},
			},
			&tg.ReplyKeyboardHide{
				Selective: true,
			},
			false,
		},
		{
			"ForceReply",
			&oas.SendMessageReplyMarkup{
				Type: oas.ForceReplySendMessageReplyMarkup,
				ForceReply: oas.ForceReply{
					ForceReply:            true,
					InputFieldPlaceholder: oas.NewOptString("placeholder"),
					Selective:             oas.NewOptBool(true),
				},
			},
			&tg.ReplyKeyboardForceReply{
				Selective:   true,
				Placeholder: "placeholder",
			},
			false,
		},
		{
			"UnknownType",
			&oas.SendMessageReplyMarkup{
				Type: "aboba",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
				got, err := api.convertToTelegramReplyMarkup(ctx, tt.input)
				if tt.wantErr {
					a.Error(err)
					return
				}
				a.NoError(err)
				setFlags(got)
				setFlags(tt.want)
				a.Equal(tt.want, got)
			})
		})
	}
}
