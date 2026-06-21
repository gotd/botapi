package botapi

import (
	"html"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
)

var mdMap = map[MessageEntityType]string{
	EntityBold:   "*",
	EntityItalic: "_",
	EntityCode:   "`",
	EntityPre:    "```",
}

var mdV2Map = map[MessageEntityType]string{
	EntityBold:                 "*",
	EntityItalic:               "_",
	EntityCode:                 "`",
	EntityPre:                  "```",
	EntityUnderline:            "__",
	EntityStrikethrough:        "~",
	EntitySpoiler:              "||",
	EntityBlockquote:           ">",
	EntityExpandableBlockquote: "**>",
}

var htmlMap = map[MessageEntityType]string{
	EntityBold:                 "b",
	EntityItalic:               "i",
	EntityCode:                 "code",
	EntityPre:                  "pre",
	EntityUnderline:            "u",
	EntityStrikethrough:        "s",
	EntitySpoiler:              "span class=\"tg-spoiler\"",
	EntityBlockquote:           "blockquote",
	EntityExpandableBlockquote: "blockquote expandable",
}

// TextAndEntities gets message or caption text and entities
func (m *Message) TextAndEntities() (string, []MessageEntity) {
	if m.Text != "" {
		return m.Text, m.Entities
	}
	return m.Caption, m.CaptionEntities
}

// OriginalMD gets the original markdown formatting of a message text.
func (m *Message) OriginalMD() string {
	return getOrigMsgMD(utf16.Encode([]rune(m.Text)), m.Entities)
}

// OriginalMDV2 gets the original markdownV2 formatting of a message text.
func (m *Message) OriginalMDV2() string {
	return getOrigMsgMDV2(utf16.Encode([]rune(m.Text)), m.Entities)
}

// OriginalHTML gets the original HTML formatting of a message text.
func (m *Message) OriginalHTML() string {
	return getOrigMsgHTML(utf16.Encode([]rune(m.Text)), m.Entities)
}

// OriginalCaptionMD gets the original markdown formatting of a message caption.
func (m *Message) OriginalCaptionMD() string {
	return getOrigMsgMD(utf16.Encode([]rune(m.Caption)), m.CaptionEntities)
}

// OriginalCaptionMDV2 gets the original markdownV2 formatting of a message caption.
func (m *Message) OriginalCaptionMDV2() string {
	return getOrigMsgMDV2(utf16.Encode([]rune(m.Caption)), m.CaptionEntities)
}

// OriginalCaptionHTML gets the original HTML formatting of a message caption.
func (m *Message) OriginalCaptionHTML() string {
	return getOrigMsgHTML(utf16.Encode([]rune(m.Caption)), m.CaptionEntities)
}

// OriginalTextMD gets the original markdown formatting of a message text or caption.
func (m *Message) OriginalTextMD() string {
	text, ents := m.TextAndEntities()
	return getOrigMsgMD(utf16.Encode([]rune(text)), ents)
}

// OriginalTextMDV2 gets the original markdownV2 formatting of a message text or caption.
func (m *Message) OriginalTextMDV2() string {
	text, ents := m.TextAndEntities()
	return getOrigMsgMDV2(utf16.Encode([]rune(text)), ents)
}

// OriginalTextHTML gets the original HTML formatting of a message text caption.
func (m *Message) OriginalTextHTML() string {
	text, ents := m.TextAndEntities()
	return getOrigMsgHTML(utf16.Encode([]rune(text)), ents)
}

// Does not support nesting. only look at upper entities.
func getOrigMsgMD(utf16Data []uint16, ents []MessageEntity) string {
	out := strings.Builder{}
	prev := 0

	for _, ent := range getUpperEntities(ents) {
		newPrev := ent.Offset + ent.Length
		prevText := string(utf16.Decode(utf16Data[prev:ent.Offset]))

		text := utf16.Decode(utf16Data[ent.Offset:newPrev])
		pre, cleanCntnt, post := splitEdgeWhitespace(string(text), ent)
		cleanCntntRune := []rune(cleanCntnt)

		switch ent.Type {
		case EntityBold, EntityItalic, EntityCode:
			out.WriteString(prevText + pre + mdMap[ent.Type] + escapeContainedMDV1(cleanCntntRune, []rune(mdMap[ent.Type])) + mdMap[ent.Type] + post)
		case EntityPre:
			if ent.Language == "" {
				out.WriteString(prevText + pre + mdMap[ent.Type] +
					escapeContainedMDV1(cleanCntntRune, []rune(mdMap[ent.Type])) + mdMap[ent.Type] + post)
			} else {
				out.WriteString(prevText + pre + mdMap[ent.Type] +
					ent.Language + "\n" + escapeContainedMDV1(cleanCntntRune, []rune(mdMap[ent.Type])) + mdMap[ent.Type] + post)
			}
		case EntityTextMention:
			out.WriteString(prevText + pre + "[" + escapeContainedMDV1(cleanCntntRune, []rune("[]()")) + "](tg://user?id=" +
				strconv.FormatInt(ent.User.ID, 10) + ")" + post)
		case EntityTextLink:
			out.WriteString(prevText + pre + "[" + escapeContainedMDV1(cleanCntntRune, []rune("[]()")) + "](" + ent.URL + ")" + post)
		default:
			continue
		}

		prev = newPrev
	}

	out.WriteString(string(utf16.Decode(utf16Data[prev:])))

	return out.String()
}

func getOrigMsgHTML(utf16Data []uint16, ents []MessageEntity) string {
	if len(ents) == 0 {
		return html.EscapeString(string(utf16.Decode(utf16Data)))
	}

	bd := strings.Builder{}
	prev := 0

	for _, e := range getUpperEntities(ents) {
		data, end := fillNestedHTML(utf16Data, e, prev, getChildEntities(e, ents))
		bd.WriteString(data)

		prev = end
	}

	bd.WriteString(html.EscapeString(string(utf16.Decode(utf16Data[prev:]))))

	return bd.String()
}

func getOrigMsgMDV2(utf16Data []uint16, ents []MessageEntity) (origMsg string) {
	if len(ents) == 0 {
		return string(utf16.Decode(utf16Data))
	}

	bd := strings.Builder{}
	prev := 0

	for _, e := range getUpperEntities(ents) {
		data, end := fillNestedMarkdownV2(utf16Data, e, prev, getChildEntities(e, ents))
		bd.WriteString(data)

		prev = end
	}

	bd.WriteString(string(utf16.Decode(utf16Data[prev:])))

	return bd.String()
}

func fillNestedHTML(data []uint16, ent MessageEntity, start int, entities []MessageEntity) (finalHTML string, entEnd int) {
	entEnd = ent.Offset + ent.Length
	if len(entities) == 0 || entEnd < entities[0].Offset {
		// no nesting; just return straight away and move to next.
		return writeFinalHTML(data, ent, start, html.EscapeString(string(utf16.Decode(data[ent.Offset:entEnd])))), entEnd
	}

	subPrev := ent.Offset
	subEnd := ent.Offset
	bd := strings.Builder{}

	for _, e := range getUpperEntities(entities) {
		if e.Offset < subEnd || e == ent {
			continue
		}

		if e.Offset >= entEnd {
			break
		}

		out, end := fillNestedHTML(data, e, subPrev, getChildEntities(e, entities))
		bd.WriteString(out)

		subPrev = end
	}

	bd.WriteString(html.EscapeString(string(utf16.Decode(data[subPrev:entEnd]))))

	return writeFinalHTML(data, ent, start, bd.String()), entEnd
}

func fillNestedMarkdownV2(
	data []uint16,
	ent MessageEntity,
	start int,
	entities []MessageEntity,
) (finalMD string, entEnd int) {
	entEnd = ent.Offset + ent.Length
	if len(entities) == 0 || entEnd < entities[0].Offset {
		// no nesting; just return straight away and move to next.
		return writeFinalMarkdownV2(data, ent, start, string(utf16.Decode(data[ent.Offset:entEnd]))), entEnd
	}

	subPrev := ent.Offset
	subEnd := ent.Offset
	bd := strings.Builder{}

	for _, e := range getUpperEntities(entities) {
		if e.Offset < subEnd || e == ent {
			continue
		}

		if e.Offset >= entEnd {
			break
		}

		out, end := fillNestedMarkdownV2(data, e, subPrev, getChildEntities(e, entities))
		bd.WriteString(out)

		subPrev = end
	}

	bd.WriteString(string(utf16.Decode(data[subPrev:entEnd])))

	return writeFinalMarkdownV2(data, ent, start, bd.String()), entEnd
}

func writeFinalHTML(data []uint16, ent MessageEntity, start int, cntnt string) string {
	prevText := html.EscapeString(string(utf16.Decode(data[start:ent.Offset])))
	switch ent.Type {
	case EntityBold, EntityItalic, EntityCode, EntityUnderline, EntityStrikethrough, EntitySpoiler:
		return prevText + "<" + htmlMap[ent.Type] + ">" + cntnt + "</" + closeHTMLTag(htmlMap[ent.Type]) + ">"
	case EntityPre:
		if ent.Language == "" {
			return prevText + "<pre>" + cntnt + "</pre>"
		}

		return prevText + `<pre><code class="` + ent.Language + `">` + cntnt + "</code></pre>"
	case EntityCustomEmoji:
		return prevText + `<tg-emoji emoji-id="` + ent.CustomEmojiID + `">` + cntnt + "</tg-emoji>"
	case EntityDateTime:
		if ent.DateTimeFormat != "" {
			return prevText + `<tg-time unix="` + strconv.Itoa(ent.UnixTime) + `" format="` + ent.DateTimeFormat + `">` + cntnt + "</tg-time>"
		}

		return prevText + `<tg-time unix="` + strconv.Itoa(ent.UnixTime) + `">` + cntnt + "</tg-time>"
	case EntityTextMention:
		return prevText + `<a href="tg://user?id=` + strconv.FormatInt(ent.User.ID, 10) + `">` + cntnt + "</a>"
	case EntityTextLink:
		return prevText + `<a href="` + ent.URL + `">` + cntnt + "</a>"
	case EntityBlockquote:
		return prevText + `<blockquote>` + cntnt + "</blockquote>"
	case EntityExpandableBlockquote:
		return prevText + `<blockquote expandable>` + cntnt + "</blockquote>"
	default:
		return prevText + cntnt
	}
}

// closeHTMLTag makes sure to generate the correct HTML closing tag for a given opening tag.
func closeHTMLTag(s string) string {
	if !strings.HasPrefix(s, "span") {
		return s
	}

	return "span"
}

func writeFinalMarkdownV2(data []uint16, ent MessageEntity, start int, cntnt string) string {
	prevText := string(utf16.Decode(data[start:ent.Offset]))
	pre, cleanCntnt, post := splitEdgeWhitespace(cntnt, ent)

	switch ent.Type {
	case EntityBold, EntityItalic, EntityCode, EntityUnderline, EntityStrikethrough, EntitySpoiler:
		return prevText + pre + mdV2Map[ent.Type] + cleanCntnt + mdV2Map[ent.Type] + post
	case EntityPre:
		if ent.Language == "" {
			return prevText + pre + "```\n" + cleanCntnt + "```" + post
		}

		return prevText + pre + "```" + ent.Language + "\n" + cleanCntnt + "```" + post
	case EntityCustomEmoji:
		return prevText + pre + "![" + cleanCntnt + "](tg://emoji?id=" + ent.CustomEmojiID + ")" + post
	case EntityDateTime:
		if ent.DateTimeFormat != "" {
			return prevText + pre + "![" + cleanCntnt + "](tg://time?unix=" +
				strconv.Itoa(ent.UnixTime) + "&format=" + ent.DateTimeFormat + ")" + post
		}

		return prevText + pre + "![" + cleanCntnt + "](tg://time?unix=" + strconv.Itoa(ent.UnixTime) + ")" + post
	case EntityTextMention:
		return prevText + pre + "[" + cleanCntnt + "](tg://user?id=" + strconv.FormatInt(ent.User.ID, 10) + ")" + post
	case EntityTextLink:
		return prevText + pre + "[" + cleanCntnt + "](" + ent.URL + ")" + post
	case EntityBlockquote:
		return prevText + pre + ">" + strings.Join(strings.Split(cleanCntnt, "\n"), "\n>") + post
	case EntityExpandableBlockquote:
		return prevText + pre + "**>" + strings.Join(strings.Split(cleanCntnt, "\n"), "\n>") + "||" + post
	default:
		return prevText + cntnt
	}
}

func getUpperEntities(ents []MessageEntity) []MessageEntity {
	prev := 0
	uppers := make([]MessageEntity, 0, len(ents))

	for _, e := range ents {
		if e.Offset < prev {
			continue
		}

		uppers = append(uppers, e)
		prev = e.Offset + e.Length
	}

	return uppers
}

func getChildEntities(ent MessageEntity, ents []MessageEntity) []MessageEntity {
	end := ent.Offset + ent.Length
	children := make([]MessageEntity, 0, len(ents))

	for _, e := range ents {
		if e.Offset < ent.Offset || e == ent {
			continue
		}

		if e.Offset >= end {
			break
		}

		children = append(children, e)
	}

	return children
}

func splitEdgeWhitespace(text string, ent MessageEntity) (pre, cntnt, post string) {
	keepNewLines := ent.Type == EntityPre

	bd := strings.Builder{}
	rText := []rune(text)

	for i := 0; i < len(rText) && unicode.IsSpace(rText[i]) && (!keepNewLines || rText[i] != '\n'); i++ {
		bd.WriteRune(rText[i])
	}

	pre = bd.String()

	text = strings.TrimPrefix(text, pre)

	bd.Reset()

	for i := len(rText) - 1; i >= 0 && unicode.IsSpace(rText[i]); i-- {
		bd.WriteRune(rText[i])
	}

	post = bd.String()

	return pre, strings.TrimSuffix(text, post), post
}

func escapeContainedMDV1(data, mdType []rune) string {
	out := strings.Builder{}

	for _, x := range data {
		if slices.Contains(mdType, x) {
			out.WriteRune('\\')
		}

		out.WriteRune(x)
	}

	return out.String()
}
