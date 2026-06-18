package botapi

import "github.com/gotd/td/tg"

// MessageEntity represents one special entity in a text message (e.g. a
// hashtag, link, or formatted run).
type MessageEntity struct {
	Type           MessageEntityType `json:"type"`
	Offset         int               `json:"offset"`
	Length         int               `json:"length"`
	URL            string            `json:"url,omitempty"`
	User           *User             `json:"user,omitempty"`
	Language       string            `json:"language,omitempty"`
	CustomEmojiID  string            `json:"custom_emoji_id,omitempty"`
	UnixTime       int               `json:"unix_time,omitempty"`
	DateTimeFormat string            `json:"date_time_format,omitempty"`
}

// Contact represents a phone contact.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
	VCard       string `json:"vcard,omitempty"`
}

// Dice represents an animated emoji with a random value.
type Dice struct {
	Emoji DiceEmoji `json:"emoji"`
	Value int       `json:"value"`
}

// Location represents a point on the map.
type Location struct {
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
}

// Venue represents a venue.
type Venue struct {
	Location        Location `json:"location"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareID    string   `json:"foursquare_id,omitempty"`
	FoursquareType  string   `json:"foursquare_type,omitempty"`
	GooglePlaceID   string   `json:"google_place_id,omitempty"`
	GooglePlaceType string   `json:"google_place_type,omitempty"`
}

// PollOption represents one answer option in a poll.
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

// Poll represents a poll.
type Poll struct {
	ID                    string          `json:"id"`
	Question              string          `json:"question"`
	Options               []PollOption    `json:"options"`
	TotalVoterCount       int             `json:"total_voter_count"`
	IsClosed              bool            `json:"is_closed,omitempty"`
	IsAnonymous           bool            `json:"is_anonymous,omitempty"`
	Type                  PollType        `json:"type"`
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers,omitempty"`
	CorrectOptionID       int             `json:"correct_option_id,omitempty"`
	Explanation           string          `json:"explanation,omitempty"`
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`
	OpenPeriod            int             `json:"open_period,omitempty"`
}

// PollAnswer represents a change of a user's answer in a non-anonymous poll.
type PollAnswer struct {
	PollID    string `json:"poll_id"`
	VoterChat *Chat  `json:"voter_chat,omitempty"`
	User      *User  `json:"user,omitempty"`
	OptionIDs []int  `json:"option_ids"`
}

// Message represents a message.
type Message struct {
	MessageID       int `json:"message_id"`
	MessageThreadID int `json:"message_thread_id,omitempty"`
	// BusinessConnectionID is the unique identifier of the business connection
	// the message was received from, for messages delivered as part of a
	// business connection.
	BusinessConnectionID string        `json:"business_connection_id,omitempty"`
	From                 *User         `json:"from,omitempty"`
	SenderChat           *Chat         `json:"sender_chat,omitempty"`
	Date                 int           `json:"date"`
	Chat                 Chat          `json:"chat"`
	ForwardOrigin        MessageOrigin `json:"forward_origin,omitempty"`
	ReplyToMessage       *Message      `json:"reply_to_message,omitempty"`
	ViaBot               *User         `json:"via_bot,omitempty"`
	EditDate             int           `json:"edit_date,omitempty"`
	HasProtectedContent  bool          `json:"has_protected_content,omitempty"`
	MediaGroupID         string        `json:"media_group_id,omitempty"`
	AuthorSignature      string        `json:"author_signature,omitempty"`

	Text     string          `json:"text,omitempty"`
	Entities []MessageEntity `json:"entities,omitempty"`

	Caption         string          `json:"caption,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	Animation *Animation  `json:"animation,omitempty"`
	Audio     *Audio      `json:"audio,omitempty"`
	Document  *Document   `json:"document,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
	Sticker   *Sticker    `json:"sticker,omitempty"`
	Video     *Video      `json:"video,omitempty"`
	VideoNote *VideoNote  `json:"video_note,omitempty"`
	Voice     *Voice      `json:"voice,omitempty"`

	Contact  *Contact  `json:"contact,omitempty"`
	Dice     *Dice     `json:"dice,omitempty"`
	Poll     *Poll     `json:"poll,omitempty"`
	Venue    *Venue    `json:"venue,omitempty"`
	Location *Location `json:"location,omitempty"`

	NewChatMembers []User   `json:"new_chat_members,omitempty"`
	LeftChatMember *User    `json:"left_chat_member,omitempty"`
	NewChatTitle   string   `json:"new_chat_title,omitempty"`
	PinnedMessage  *Message `json:"pinned_message,omitempty"`

	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// raw holds the original MTProto message this Message was converted from, or
	// nil for synthesized stubs (e.g. reply_to/pinned placeholders). It is not
	// serialized; reach it through Raw.
	raw *tg.Message

	// businessPeer is the input peer of a business message's chat, built from the
	// access hash delivered in the update's entities. A reply on behalf of a
	// business account must use this peer (the account's own access hash for the
	// chat) rather than the bot's stored one, which is invalid in the business
	// context (BUSINESS_PEER_INVALID). Nil for non-business messages.
	businessPeer tg.InputPeerClass
}

// Raw returns the original MTProto message this Message was converted from, for
// anything the typed Bot API surface does not expose. It is nil for messages
// that were not converted from a live update — synthesized reply_to and pinned
// placeholders, or values produced by JSON decoding.
func (m *Message) Raw() *tg.Message { return m.raw }
