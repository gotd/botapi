package botapi

import (
	"strconv"

	"github.com/gotd/td/tg"
)

// entitiesToTg converts Bot API message entities into MTProto entities.
//
// Offsets and lengths are passed through unchanged: both APIs count in UTF-16
// code units. The switch over MessageEntityType is exhaustive (exhaustive lint).
func entitiesToTg(entities []MessageEntity) []tg.MessageEntityClass {
	if len(entities) == 0 {
		return nil
	}

	out := make([]tg.MessageEntityClass, 0, len(entities))
	for _, e := range entities {
		off, length := e.Offset, e.Length
		switch e.Type {
		case EntityMention:
			out = append(out, &tg.MessageEntityMention{Offset: off, Length: length})
		case EntityHashtag:
			out = append(out, &tg.MessageEntityHashtag{Offset: off, Length: length})
		case EntityCashtag:
			out = append(out, &tg.MessageEntityCashtag{Offset: off, Length: length})
		case EntityBotCommand:
			out = append(out, &tg.MessageEntityBotCommand{Offset: off, Length: length})
		case EntityURL:
			out = append(out, &tg.MessageEntityURL{Offset: off, Length: length})
		case EntityEmail:
			out = append(out, &tg.MessageEntityEmail{Offset: off, Length: length})
		case EntityPhoneNumber:
			out = append(out, &tg.MessageEntityPhone{Offset: off, Length: length})
		case EntityBold:
			out = append(out, &tg.MessageEntityBold{Offset: off, Length: length})
		case EntityItalic:
			out = append(out, &tg.MessageEntityItalic{Offset: off, Length: length})
		case EntityUnderline:
			out = append(out, &tg.MessageEntityUnderline{Offset: off, Length: length})
		case EntityStrikethrough:
			out = append(out, &tg.MessageEntityStrike{Offset: off, Length: length})
		case EntitySpoiler:
			out = append(out, &tg.MessageEntitySpoiler{Offset: off, Length: length})
		case EntityBlockquote:
			out = append(out, &tg.MessageEntityBlockquote{Offset: off, Length: length})
		case EntityExpandableBlockquote:
			out = append(out, &tg.MessageEntityBlockquote{Offset: off, Length: length, Collapsed: true})
		case EntityCode:
			out = append(out, &tg.MessageEntityCode{Offset: off, Length: length})
		case EntityPre:
			out = append(out, &tg.MessageEntityPre{Offset: off, Length: length, Language: e.Language})
		case EntityTextLink:
			out = append(out, &tg.MessageEntityTextURL{Offset: off, Length: length, URL: e.URL})
		case EntityTextMention:
			var userID int64

			if e.User != nil {
				userID = e.User.ID
			}

			out = append(out, &tg.MessageEntityMentionName{Offset: off, Length: length, UserID: userID})
		case EntityCustomEmoji:
			id, _ := strconv.ParseInt(e.CustomEmojiID, 10, 64)

			out = append(out, &tg.MessageEntityCustomEmoji{Offset: off, Length: length, DocumentID: id})
		default:
			// Unknown entity type: skip rather than emit an invalid entity.
		}
	}

	return out
}

// entitiesFromTg converts MTProto entities into Bot API message entities.
//
// Text-mention entities carry only the user id here; full user resolution is
// performed by the message converter when peer data is available.
func entitiesFromTg(entities []tg.MessageEntityClass) []MessageEntity {
	if len(entities) == 0 {
		return nil
	}

	out := make([]MessageEntity, 0, len(entities))
	for _, e := range entities {
		me := MessageEntity{Offset: e.GetOffset(), Length: e.GetLength()}
		switch e := e.(type) {
		case *tg.MessageEntityMention:
			me.Type = EntityMention
		case *tg.MessageEntityHashtag:
			me.Type = EntityHashtag
		case *tg.MessageEntityCashtag:
			me.Type = EntityCashtag
		case *tg.MessageEntityBotCommand:
			me.Type = EntityBotCommand
		case *tg.MessageEntityURL:
			me.Type = EntityURL
		case *tg.MessageEntityEmail:
			me.Type = EntityEmail
		case *tg.MessageEntityPhone:
			me.Type = EntityPhoneNumber
		case *tg.MessageEntityBold:
			me.Type = EntityBold
		case *tg.MessageEntityItalic:
			me.Type = EntityItalic
		case *tg.MessageEntityUnderline:
			me.Type = EntityUnderline
		case *tg.MessageEntityStrike:
			me.Type = EntityStrikethrough
		case *tg.MessageEntitySpoiler:
			me.Type = EntitySpoiler
		case *tg.MessageEntityBlockquote:
			if e.Collapsed {
				me.Type = EntityExpandableBlockquote
			} else {
				me.Type = EntityBlockquote
			}
		case *tg.MessageEntityCode:
			me.Type = EntityCode
		case *tg.MessageEntityPre:
			me.Type = EntityPre
			me.Language = e.Language
		case *tg.MessageEntityTextURL:
			me.Type = EntityTextLink
			me.URL = e.URL
		case *tg.MessageEntityMentionName:
			me.Type = EntityTextMention
			me.User = &User{ID: e.UserID}
		case *tg.MessageEntityCustomEmoji:
			me.Type = EntityCustomEmoji
			me.CustomEmojiID = strconv.FormatInt(e.DocumentID, 10)
		default:
			// Unknown or bot-irrelevant entity (e.g. unknown future types): skip.
			continue
		}

		out = append(out, me)
	}

	return out
}
