package main

import "github.com/gotd/botapi"

// displayName renders a friendly name for a (possibly nil) user.
func displayName(u *botapi.User) string {
	switch {
	case u == nil:
		return "stranger"
	case u.Username != "":
		return "@" + u.Username
	default:
		return u.FirstName
	}
}

// reverse returns s with its runes in reverse order.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return string(r)
}

// --- incoming-media predicates, shared by media.go ---

func hasPhoto(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && len(m.Photo) > 0
}

func hasDocument(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && m.Document != nil
}

func hasSticker(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && m.Sticker != nil
}

func hasLocation(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && m.Location != nil
}

func hasContact(c *botapi.Context) bool {
	m := c.Message()
	return m != nil && m.Contact != nil
}

const helpText = `<b>gotd/botapi demo bot</b>

<b>Formatting</b>
/html — HTML formatting
/md — MarkdownV2 formatting
/rich — structured rich message

<b>Content</b>
/poll — send a poll
/dice — roll a dice
/location, /venue, /contact — places &amp; people
/silent — message with no notification
/protect — protected-content message
/remind — background send after a delay

<b>Media</b>
/photo — photo by URL
/document — generated in-memory file
/album — media-group album
/typing — chat action

<b>Streaming (live drafts)</b>
/stream — stream a rich message, then persist it
/streamtext — stream a plain-text message, then persist it

<b>Keyboards</b>
/keyboard — inline keyboard + callbacks
/removekbd — remove the reply keyboard

<b>Editing</b>
/edit — send then edit a message
/forward — forward your message back
/selfdestruct — send then delete a message
/ref — this chat's serializable PeerRef

<b>In groups</b>
/info — chat info &amp; your status
/pin — react to and pin your command
/raw — the raw MTProto escape hatch

<b>Also try</b>
• send me a photo, document, sticker, location or contact
• edit one of your messages
• type <code>@thisbot hello</code> in any chat (inline mode)
• say "hi"`
