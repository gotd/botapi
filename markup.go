package botapi

// ReplyMarkup is a sealed union of the markup objects a message can carry. The
// unexported marker method makes the set closed: only the types in this package
// can satisfy it, so a type switch over a ReplyMarkup is exhaustive.
//
// Concrete variants: *InlineKeyboardMarkup, *ReplyKeyboardMarkup,
// *ReplyKeyboardRemove, *ForceReply.
type ReplyMarkup interface {
	isReplyMarkup()
}

// WebAppInfo describes a Web App to be opened from a button.
type WebAppInfo struct {
	URL string `json:"url"`
}

// InlineKeyboardButton is one button of an inline keyboard. Exactly one of the
// optional action fields should be set.
type InlineKeyboardButton struct {
	Text                         string      `json:"text"`
	URL                          string      `json:"url,omitempty"`
	CallbackData                 string      `json:"callback_data,omitempty"`
	WebApp                       *WebAppInfo `json:"web_app,omitempty"`
	SwitchInlineQuery            *string     `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat *string     `json:"switch_inline_query_current_chat,omitempty"`
	Pay                          bool        `json:"pay,omitempty"`
}

// InlineKeyboardMarkup is an inline keyboard that appears beneath a message.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// KeyboardButtonPollType is the poll-type constraint of a poll-request button.
type KeyboardButtonPollType struct {
	Type PollType `json:"type,omitempty"`
}

// KeyboardButton is one button of a reply (custom) keyboard. Text is sent as a
// plain message when no request/web-app field is set.
type KeyboardButton struct {
	Text            string                  `json:"text"`
	RequestContact  bool                    `json:"request_contact,omitempty"`
	RequestLocation bool                    `json:"request_location,omitempty"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo             `json:"web_app,omitempty"`
}

// ReplyKeyboardMarkup is a custom keyboard with reply options.
type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	IsPersistent          bool               `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"`
	Selective             bool               `json:"selective,omitempty"`
}

// ReplyKeyboardRemove removes the current custom keyboard.
type ReplyKeyboardRemove struct {
	// RemoveKeyboard is always true; the type itself signals removal.
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

// ForceReply forces the user's client to display a reply interface.
type ForceReply struct {
	// ForceReply is always true; the type itself signals the intent.
	ForceReply            bool   `json:"force_reply"`
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	Selective             bool   `json:"selective,omitempty"`
}

func (*InlineKeyboardMarkup) isReplyMarkup() {}
func (*ReplyKeyboardMarkup) isReplyMarkup()  {}
func (*ReplyKeyboardRemove) isReplyMarkup()  {}
func (*ForceReply) isReplyMarkup()           {}

// --- Constructors / builders (the type-safe telegoutil equivalent) ---

// InlineKeyboard builds an inline keyboard from rows of buttons.
func InlineKeyboard(rows ...[]InlineKeyboardButton) *InlineKeyboardMarkup {
	return &InlineKeyboardMarkup{InlineKeyboard: rows}
}

// InlineRow groups inline buttons into a single keyboard row.
func InlineRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton { return buttons }

// InlineButtonURL builds an inline button that opens a URL.
func InlineButtonURL(text, url string) InlineKeyboardButton {
	return InlineKeyboardButton{Text: text, URL: url}
}

// InlineButtonData builds an inline button that sends callback data.
func InlineButtonData(text, data string) InlineKeyboardButton {
	return InlineKeyboardButton{Text: text, CallbackData: data}
}

// Keyboard builds a reply keyboard from rows of buttons.
func Keyboard(rows ...[]KeyboardButton) *ReplyKeyboardMarkup {
	return &ReplyKeyboardMarkup{Keyboard: rows}
}

// Row groups reply-keyboard buttons into a single row.
func Row(buttons ...KeyboardButton) []KeyboardButton { return buttons }

// Button builds a plain-text reply-keyboard button.
func Button(text string) KeyboardButton { return KeyboardButton{Text: text} }

// ButtonContact builds a reply-keyboard button that requests the user's contact.
func ButtonContact(text string) KeyboardButton {
	return KeyboardButton{Text: text, RequestContact: true}
}

// ButtonLocation builds a reply-keyboard button that requests the user's location.
func ButtonLocation(text string) KeyboardButton {
	return KeyboardButton{Text: text, RequestLocation: true}
}

// RemoveKeyboard builds a ReplyKeyboardRemove.
func RemoveKeyboard() *ReplyKeyboardRemove { return &ReplyKeyboardRemove{RemoveKeyboard: true} }
