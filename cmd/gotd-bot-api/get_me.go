package main

import (
	"context"
	"encoding/json"

	"github.com/gotd/botapi/api"
	"github.com/gotd/td/tg"
)

func convertUser(res *tg.User) api.User {
	// TDLib uses special flag USER_FLAG_IS_INLINE_BOT, which is not defined in schema.
	//
	// See links for reference.
	//
	// User object JSON encoding
	// https://github.com/tdlib/telegram-bot-api/blob/81f298361cf80d1d6c70a074ff88534bd3d450b3/telegram-bot-api/Client.cpp#L335
	//
	// TDLib API (td_api.tl) user constructor <-> BotAPI UserInfo structure conversion.
	// https://github.com/tdlib/telegram-bot-api/blob/81f298361cf80d1d6c70a074ff88534bd3d450b3/telegram-bot-api/Client.cpp#L8211
	//
	// TDLib User type <-> TDLib API (td_api.tl) user constructor conversion.
	// https://github.com/tdlib/td/blob/c45535d607463adb0cd20fcadf43e8f793b1fb24/td/telegram/ContactsManager.cpp#L15782-L15783
	//
	// Telegram API user constructor <-> TDLib User type conversion.
	// https://github.com/tdlib/td/blob/c45535d607463adb0cd20fcadf43e8f793b1fb24/td/telegram/ContactsManager.cpp#L8156
	//
	// USER_FLAG_IS_INLINE_BOT definition.
	// https://github.com/tdlib/td/blob/c45535d607463adb0cd20fcadf43e8f793b1fb24/td/telegram/ContactsManager.h#L990
	isInlineBot := res.Flags.Has(19)

	return api.User{
		ID:              res.ID,
		FirstName:       res.FirstName,
		LastName:        res.LastName,
		Username:        res.Username,
		LanguageCode:    res.LangCode,
		IsBot:           res.Bot,
		CanJoinGroups:   !res.BotNochats,
		CanReadMessages: res.BotChatHistory,
		SupportsInline:  isInlineBot,
	}
}

func getMe(ctx context.Context, h handleContext) error {
	res, err := h.Client.Self(ctx)
	if err != nil {
		return err
	}

	return json.NewEncoder(h.Writer).Encode(api.Response{
		Result: convertUser(res),
	})
}
