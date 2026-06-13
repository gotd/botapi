package botapi

// InputMedia is a sealed union describing the content of a media message to be
// sent (e.g. in a media group or via an edit).
//
// Concrete variants: *InputMediaPhoto, *InputMediaVideo, *InputMediaAnimation,
// *InputMediaAudio, *InputMediaDocument.
//
// The Media (and Thumbnail) fields hold an InputFile; the send path (Phase 3)
// resolves each to a file_id, URL or upload.
type InputMedia interface {
	isInputMedia()
}

// InputMediaPhoto is a photo to be sent.
type InputMediaPhoto struct {
	Type            InputMediaType  `json:"type"`
	Media           InputFile       `json:"-"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	HasSpoiler      bool            `json:"has_spoiler,omitempty"`
}

// InputMediaVideo is a video to be sent.
type InputMediaVideo struct {
	Type              InputMediaType  `json:"type"`
	Media             InputFile       `json:"-"`
	Thumbnail         InputFile       `json:"-"`
	Caption           string          `json:"caption,omitempty"`
	ParseMode         ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities   []MessageEntity `json:"caption_entities,omitempty"`
	Width             int             `json:"width,omitempty"`
	Height            int             `json:"height,omitempty"`
	Duration          int             `json:"duration,omitempty"`
	SupportsStreaming bool            `json:"supports_streaming,omitempty"`
	HasSpoiler        bool            `json:"has_spoiler,omitempty"`
}

// InputMediaAnimation is an animation (GIF or silent H.264/MPEG-4 AVC) to be sent.
type InputMediaAnimation struct {
	Type            InputMediaType  `json:"type"`
	Media           InputFile       `json:"-"`
	Thumbnail       InputFile       `json:"-"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Width           int             `json:"width,omitempty"`
	Height          int             `json:"height,omitempty"`
	Duration        int             `json:"duration,omitempty"`
	HasSpoiler      bool            `json:"has_spoiler,omitempty"`
}

// InputMediaAudio is an audio file to be sent.
type InputMediaAudio struct {
	Type            InputMediaType  `json:"type"`
	Media           InputFile       `json:"-"`
	Thumbnail       InputFile       `json:"-"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Duration        int             `json:"duration,omitempty"`
	Performer       string          `json:"performer,omitempty"`
	Title           string          `json:"title,omitempty"`
}

// InputMediaDocument is a general file to be sent.
type InputMediaDocument struct {
	Type                        InputMediaType  `json:"type"`
	Media                       InputFile       `json:"-"`
	Thumbnail                   InputFile       `json:"-"`
	Caption                     string          `json:"caption,omitempty"`
	ParseMode                   ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities             []MessageEntity `json:"caption_entities,omitempty"`
	DisableContentTypeDetection bool            `json:"disable_content_type_detection,omitempty"`
}

func (*InputMediaPhoto) isInputMedia()     {}
func (*InputMediaVideo) isInputMedia()     {}
func (*InputMediaAnimation) isInputMedia() {}
func (*InputMediaAudio) isInputMedia()     {}
func (*InputMediaDocument) isInputMedia()  {}
