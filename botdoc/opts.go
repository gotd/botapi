package botdoc

import "strings"

var (
	wellKnownTypes = []string{
		"Update",
		"Message",
		"MessageId",
		"User",
		"Chat",
		"File",
		"Poll",
		"BotCommand",
		"GameHighScore",
		"WebhookInfo",
		"UserProfilePhotos",
		"ChatMember",
		"ChatInviteLink",
	}

	isIDLikeName = createMatcher([]string{
		"chat_id",
		"user_id",
	}, strings.Contains)
	isIDLikeDesc = createMatcher([]string{
		"Unique identifier for this user",
		"Unique identifier for this chat",
	}, strings.Contains)

	discriminatorFields = []string{
		"type",
		"chat_type",
		"source",
		"status",
	}
	isDiscriminatorField = createMatcher(discriminatorFields, strings.EqualFold)

	isExclusiveMinimum = createMatcher([]string{
		"width",
		"height",
		"duration",
	}, strings.Contains)
)

func createMatcher(s []string, fn func(a, b string) bool) func(string) bool {
	return func(input string) bool {
		for _, elem := range s {
			if fn(input, elem) {
				return true
			}
		}
		return false
	}
}
