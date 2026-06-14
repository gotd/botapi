package botapi

import "github.com/go-faster/jx"

// MessageOrigin is a sealed union describing the original sender of a forwarded
// message.
//
// Concrete variants: *MessageOriginUser, *MessageOriginHiddenUser,
// *MessageOriginChat, *MessageOriginChannel.
type MessageOrigin interface {
	isMessageOrigin()
	// Encode writes the origin as a JSON object, including its "type"
	// discriminator, through the jx encoder.
	Encode(e *jx.Encoder)
}

// MessageOriginUser is a message originally sent by a known user.
type MessageOriginUser struct {
	Type       MessageOriginType `json:"type"`
	Date       int               `json:"date"`
	SenderUser User              `json:"sender_user"`
}

// MessageOriginHiddenUser is a message originally sent by a user who hid their
// account.
type MessageOriginHiddenUser struct {
	Type           MessageOriginType `json:"type"`
	Date           int               `json:"date"`
	SenderUserName string            `json:"sender_user_name"`
}

// MessageOriginChat is a message originally sent on behalf of a chat.
type MessageOriginChat struct {
	Type            MessageOriginType `json:"type"`
	Date            int               `json:"date"`
	SenderChat      Chat              `json:"sender_chat"`
	AuthorSignature string            `json:"author_signature,omitempty"`
}

// MessageOriginChannel is a message originally sent to a channel.
type MessageOriginChannel struct {
	Type            MessageOriginType `json:"type"`
	Date            int               `json:"date"`
	Chat            Chat              `json:"chat"`
	MessageID       int               `json:"message_id"`
	AuthorSignature string            `json:"author_signature,omitempty"`
}

func (*MessageOriginUser) isMessageOrigin()       {}
func (*MessageOriginHiddenUser) isMessageOrigin() {}
func (*MessageOriginChat) isMessageOrigin()       {}
func (*MessageOriginChannel) isMessageOrigin()    {}
