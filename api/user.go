package api

// User represents "User" api type.
// This object represents a Telegram user or bot.
//
// https://core.telegram.org/bots/api#user
type User struct {
	ID int `json:"id"`

	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsBot        bool   `json:"is_bot"`

	// Returns only in getMe.
	CanJoinGroups   bool `json:"can_join_groups"`
	CanReadMessages bool `json:"can_read_all_group_messages"`
	SupportsInline  bool `json:"supports_inline_queries"`
}
