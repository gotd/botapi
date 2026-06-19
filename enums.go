package botapi

// This file defines the typed string enums of the Bot API surface. Each is a
// distinct named type over string so callers cannot pass an arbitrary string
// where an enum is expected, and the valid values are exported constants.

// ParseMode is the formatting mode for message text and captions.
//
// See https://core.telegram.org/bots/api#formatting-options.
type ParseMode string

const (
	// ParseModeNone sends text without any entity parsing.
	ParseModeNone ParseMode = ""
	// ParseModeHTML parses a subset of HTML tags.
	ParseModeHTML ParseMode = "HTML"
	// ParseModeMarkdownV2 parses MarkdownV2-style formatting.
	ParseModeMarkdownV2 ParseMode = "MarkdownV2"
	// ParseModeMarkdown is the legacy Markdown mode; prefer ParseModeMarkdownV2.
	ParseModeMarkdown ParseMode = "Markdown"
)

// ChatType is the kind of a chat.
type ChatType string

const (
	ChatTypePrivate    ChatType = "private"
	ChatTypeGroup      ChatType = "group"
	ChatTypeSupergroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
	// ChatTypeSender is used for inline queries sent from the inline mode of a
	// private chat with the bot.
	ChatTypeSender ChatType = "sender"
)

// ChatAction is a status reported to a chat via SendChatAction.
type ChatAction string

const (
	ChatActionTyping          ChatAction = "typing"
	ChatActionUploadPhoto     ChatAction = "upload_photo"
	ChatActionRecordVideo     ChatAction = "record_video"
	ChatActionUploadVideo     ChatAction = "upload_video"
	ChatActionRecordVoice     ChatAction = "record_voice"
	ChatActionUploadVoice     ChatAction = "upload_voice"
	ChatActionUploadDocument  ChatAction = "upload_document"
	ChatActionChooseSticker   ChatAction = "choose_sticker"
	ChatActionFindLocation    ChatAction = "find_location"
	ChatActionRecordVideoNote ChatAction = "record_video_note"
	ChatActionUploadVideoNote ChatAction = "upload_video_note"
)

// MessageEntityType is the kind of a MessageEntity.
type MessageEntityType string

const (
	EntityMention              MessageEntityType = "mention"
	EntityHashtag              MessageEntityType = "hashtag"
	EntityCashtag              MessageEntityType = "cashtag"
	EntityBotCommand           MessageEntityType = "bot_command"
	EntityURL                  MessageEntityType = "url"
	EntityEmail                MessageEntityType = "email"
	EntityPhoneNumber          MessageEntityType = "phone_number"
	EntityBold                 MessageEntityType = "bold"
	EntityItalic               MessageEntityType = "italic"
	EntityUnderline            MessageEntityType = "underline"
	EntityStrikethrough        MessageEntityType = "strikethrough"
	EntitySpoiler              MessageEntityType = "spoiler"
	EntityBlockquote           MessageEntityType = "blockquote"
	EntityExpandableBlockquote MessageEntityType = "expandable_blockquote"
	EntityCode                 MessageEntityType = "code"
	EntityPre                  MessageEntityType = "pre"
	EntityTextLink             MessageEntityType = "text_link"
	EntityTextMention          MessageEntityType = "text_mention"
	EntityCustomEmoji          MessageEntityType = "custom_emoji"
	EntityDateTime             MessageEntityType = "date_time"
)

// ChatMemberStatus is a member's status in a chat.
type ChatMemberStatus string

const (
	StatusCreator       ChatMemberStatus = "creator"
	StatusAdministrator ChatMemberStatus = "administrator"
	StatusMember        ChatMemberStatus = "member"
	StatusRestricted    ChatMemberStatus = "restricted"
	StatusLeft          ChatMemberStatus = "left"
	StatusBanned        ChatMemberStatus = "kicked"
)

// PollType is the kind of a poll.
type PollType string

const (
	PollRegular PollType = "regular"
	PollQuiz    PollType = "quiz"
)

// StickerType is the kind of a sticker (or sticker set).
type StickerType string

const (
	StickerRegular     StickerType = "regular"
	StickerMask        StickerType = "mask"
	StickerCustomEmoji StickerType = "custom_emoji"
)

// MessageOriginType discriminates a MessageOrigin variant.
type MessageOriginType string

const (
	OriginUser       MessageOriginType = "user"
	OriginHiddenUser MessageOriginType = "hidden_user"
	OriginChat       MessageOriginType = "chat"
	OriginChannel    MessageOriginType = "channel"
)

// ReactionTypeKind discriminates a ReactionType variant.
type ReactionTypeKind string

const (
	ReactionEmoji       ReactionTypeKind = "emoji"
	ReactionCustomEmoji ReactionTypeKind = "custom_emoji"
	ReactionPaid        ReactionTypeKind = "paid"
)

// MenuButtonType discriminates a MenuButton variant.
type MenuButtonType string

const (
	MenuButtonCommandsType MenuButtonType = "commands"
	MenuButtonWebAppType   MenuButtonType = "web_app"
	MenuButtonDefaultType  MenuButtonType = "default"
)

// InputMediaType discriminates an InputMedia variant.
type InputMediaType string

const (
	InputMediaPhotoType     InputMediaType = "photo"
	InputMediaVideoType     InputMediaType = "video"
	InputMediaAnimationType InputMediaType = "animation"
	InputMediaAudioType     InputMediaType = "audio"
	InputMediaDocumentType  InputMediaType = "document"
)

// DiceEmoji is the emoji on which a dice-style value is based.
type DiceEmoji string

const (
	DiceDie        DiceEmoji = "🎲"
	DiceDart       DiceEmoji = "🎯"
	DiceBasketball DiceEmoji = "🏀"
	DiceFootball   DiceEmoji = "⚽"
	DiceBowling    DiceEmoji = "🎳"
	DiceSlot       DiceEmoji = "🎰"
)
