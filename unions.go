package botapi

import (
	"encoding/json"
	"strconv"
)

// ChatID identifies a target chat. The Bot API accepts either a numeric chat id
// or an @username; this sealed union represents exactly those two cases, so an
// illegal "both/neither" state is unrepresentable.
//
// Construct with ID or Username.
type ChatID interface {
	isChatID()
	// json.Marshaler so a ChatID serializes to the bare int or string the wire
	// format expects.
	json.Marshaler
}

// ChatIDInt is a numeric chat identifier.
type ChatIDInt int64

// ChatIDUsername is an @username target (with or without the leading @).
type ChatIDUsername string

func (ChatIDInt) isChatID()      {}
func (ChatIDUsername) isChatID() {}

// MarshalJSON encodes the id as a JSON number.
func (c ChatIDInt) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(c), 10)), nil
}

// MarshalJSON encodes the username as a JSON string.
func (c ChatIDUsername) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

// ID targets a chat by numeric identifier.
func ID(id int64) ChatID { return ChatIDInt(id) }

// Username targets a chat by @username.
func Username(username string) ChatID { return ChatIDUsername(username) }

// InputFile is a sealed union describing a file to send: an existing Telegram
// file_id, an HTTP URL Telegram fetches, or a local upload.
//
// Construct with FileID, FileURL, FileFromPath, FileFromBytes or FileFromReader.
type InputFile interface {
	isInputFile()
}

// InputFileID references a file already on Telegram's servers by file_id.
type InputFileID string

// InputFileURL references a file by HTTP URL for Telegram to fetch.
type InputFileURL string

// InputFileUpload is a local file to be uploaded. Exactly one source is set;
// the send path (Phase 3) chooses the uploader accordingly.
type InputFileUpload struct {
	// Name is the filename reported to Telegram.
	Name string
	// Path, when non-empty, is read from disk.
	Path string
	// Bytes, when non-nil, is the in-memory content.
	Bytes []byte
	// Reader, when non-nil, streams the content.
	Reader interface{ Read([]byte) (int, error) }
}

func (InputFileID) isInputFile()      {}
func (InputFileURL) isInputFile()     {}
func (*InputFileUpload) isInputFile() {}

// FileID references an existing Telegram file by its file_id.
func FileID(id string) InputFile { return InputFileID(id) }

// FileURL references a remote file by URL for Telegram to fetch.
func FileURL(url string) InputFile { return InputFileURL(url) }

// FileFromPath uploads a local file from disk.
func FileFromPath(path string) InputFile { return &InputFileUpload{Path: path} }

// FileFromBytes uploads in-memory content under the given filename.
func FileFromBytes(name string, data []byte) InputFile {
	return &InputFileUpload{Name: name, Bytes: data}
}

// FileFromReader uploads streamed content under the given filename.
func FileFromReader(name string, r interface{ Read([]byte) (int, error) }) InputFile {
	return &InputFileUpload{Name: name, Reader: r}
}

// ReactionType is a sealed union of reaction kinds.
//
// Concrete variants: ReactionTypeEmoji, ReactionTypeCustomEmoji, ReactionTypePaid.
type ReactionType interface {
	isReactionType()
}

// ReactionTypeEmoji is a reaction with a standard emoji.
type ReactionTypeEmoji struct {
	Type  ReactionTypeKind `json:"type"`
	Emoji string           `json:"emoji"`
}

// ReactionTypeCustomEmoji is a reaction with a custom emoji.
type ReactionTypeCustomEmoji struct {
	Type          ReactionTypeKind `json:"type"`
	CustomEmojiID string           `json:"custom_emoji_id"`
}

// ReactionTypePaid is a paid (star) reaction.
type ReactionTypePaid struct {
	Type ReactionTypeKind `json:"type"`
}

func (ReactionTypeEmoji) isReactionType()       {}
func (ReactionTypeCustomEmoji) isReactionType() {}
func (ReactionTypePaid) isReactionType()        {}

// Emoji builds a standard-emoji reaction.
func Emoji(emoji string) ReactionType {
	return ReactionTypeEmoji{Type: ReactionEmoji, Emoji: emoji}
}

// CustomEmoji builds a custom-emoji reaction.
func CustomEmoji(id string) ReactionType {
	return ReactionTypeCustomEmoji{Type: ReactionCustomEmoji, CustomEmojiID: id}
}

// MenuButton is a sealed union describing the bot's menu button.
//
// Concrete variants: MenuButtonCommands, MenuButtonWebApp, MenuButtonDefault.
type MenuButton interface {
	isMenuButton()
}

// MenuButtonCommands opens the bot's command list.
type MenuButtonCommands struct {
	Type MenuButtonType `json:"type"`
}

// MenuButtonWebApp opens a Web App.
type MenuButtonWebApp struct {
	Type   MenuButtonType `json:"type"`
	Text   string         `json:"text"`
	WebApp WebAppInfo     `json:"web_app"`
}

// MenuButtonDefault is the default menu button.
type MenuButtonDefault struct {
	Type MenuButtonType `json:"type"`
}

func (MenuButtonCommands) isMenuButton() {}
func (MenuButtonWebApp) isMenuButton()   {}
func (MenuButtonDefault) isMenuButton()  {}
