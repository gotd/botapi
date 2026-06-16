package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// PreparedKeyboardButton is the result of SavePreparedKeyboardButton: a button a
// user of a Mini App can later add to a reply keyboard.
type PreparedKeyboardButton struct {
	// ID is the unique identifier of the prepared button.
	ID string `json:"id"`
}

// SavePreparedKeyboardButton stores a keyboard button that a user of a Mini App
// can later add to a reply keyboard. The button must request users or a chat
// (i.e. have RequestUsers or RequestChat set).
func (b *Bot) SavePreparedKeyboardButton(ctx context.Context, userID int64, button KeyboardButton) (*PreparedKeyboardButton, error) {
	if button.RequestUsers == nil && button.RequestChat == nil {
		return nil, &Error{Code: 400, Description: "Bad Request: button must request users or a chat"}
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.BotsRequestWebViewButton(ctx, &tg.BotsRequestWebViewButtonRequest{
		UserID: user,
		Button: keyboardButtonToTg(button),
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	return &PreparedKeyboardButton{ID: res.WebappReqID}, nil
}
