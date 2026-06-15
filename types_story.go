package botapi

// Story represents a story.
type Story struct {
	// Chat is the chat that posted the story.
	Chat Chat `json:"chat"`
	// ID is the unique identifier of the story in the chat.
	ID int `json:"id"`
}

// InputStoryContent is a sealed union describing the content of a story to post.
//
// Concrete variants: InputStoryContentPhoto, InputStoryContentVideo.
type InputStoryContent interface {
	isInputStoryContent()
}

// InputStoryContentPhoto describes a photo to post as a story.
type InputStoryContentPhoto struct {
	// Photo is the photo to post: an uploaded file, a URL or an existing
	// file_id.
	Photo InputFile
}

// InputStoryContentVideo describes a video to post as a story.
type InputStoryContentVideo struct {
	// Video is the video to post: an uploaded file, a URL or an existing
	// file_id.
	Video InputFile
	// Duration is the precise duration of the video, in seconds.
	Duration float64
	// CoverFrameTimestamp is the timestamp, in seconds, of the frame used as the
	// story cover.
	CoverFrameTimestamp float64
	// IsAnimation reports whether the video has no sound and should be treated
	// as an animation.
	IsAnimation bool
}

func (InputStoryContentPhoto) isInputStoryContent() {}
func (InputStoryContentVideo) isInputStoryContent() {}
