package botapi

// Update represents one incoming update. At most one of the optional fields is
// present in any given update.
type Update struct {
	UpdateID           int                 `json:"update_id"`
	Message            *Message            `json:"message,omitempty"`
	EditedMessage      *Message            `json:"edited_message,omitempty"`
	ChannelPost        *Message            `json:"channel_post,omitempty"`
	EditedChannelPost  *Message            `json:"edited_channel_post,omitempty"`
	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`
	MyChatMember       *ChatMemberUpdated  `json:"my_chat_member,omitempty"`
	ChatMember         *ChatMemberUpdated  `json:"chat_member,omitempty"`

	// BusinessConnection is set when the bot is connected to or disconnected from
	// a business account, or a connection setting was changed.
	BusinessConnection *BusinessConnection `json:"business_connection,omitempty"`
	// BusinessMessage is a new message from a connected business account.
	BusinessMessage *Message `json:"business_message,omitempty"`
	// EditedBusinessMessage is an edited message from a connected business
	// account.
	EditedBusinessMessage *Message `json:"edited_business_message,omitempty"`
	// DeletedBusinessMessages reports messages deleted from a connected business
	// account.
	DeletedBusinessMessages *BusinessMessagesDeleted `json:"deleted_business_messages,omitempty"`

	// botUsername is this bot's @username (without the @), set by the router so
	// command predicates can tell a command targeted at this bot ("/cmd@me")
	// from one targeted at another ("/cmd@other"). Not part of the Bot API
	// payload, so it is unexported and not serialized.
	botUsername string
}
