package botapi

// CorrectMethod fixes legacy method name to actual.
func CorrectMethod(m string) string {
	switch m {
	case "getChatMembersCount":
		// See https://core.telegram.org/bots/api#june-25-2021.
		return "getChatMemberCount"
	case "kickChatMember":
		// See https://core.telegram.org/bots/api#june-25-2021.
		return "banChatMember"
	default:
		return m
	}
}
