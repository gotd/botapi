package botapi

import "context"

// BusinessContext is a handle scoped to a single business connection. It lets the
// bot act on behalf of a connected business account without threading the
// connection id through every call.
//
// Obtain one with Bot.Business, passing the business_connection_id received in
// business updates (or returned by GetBusinessConnection). Its methods are thin
// wrappers over the SetBusinessAccount* / GetBusinessAccount* / DeleteBusinessMessages
// methods on *Bot, and require the same business bot rights.
type BusinessContext struct {
	bot          *Bot
	connectionID string
}

// Business returns a handle that scopes business-account operations to the given
// connection id.
func (b *Bot) Business(connectionID string) *BusinessContext {
	return &BusinessContext{bot: b, connectionID: connectionID}
}

// ConnectionID returns the business connection id this context is bound to.
func (c *BusinessContext) ConnectionID() string { return c.connectionID }

// Connection returns information about this business connection.
func (c *BusinessContext) Connection(ctx context.Context) (*BusinessConnection, error) {
	return c.bot.GetBusinessConnection(ctx, c.connectionID)
}

// SetName changes the first and last name of the connected business account. The
// bot must have the can_edit_name business bot right.
func (c *BusinessContext) SetName(ctx context.Context, firstName, lastName string) error {
	return c.bot.SetBusinessAccountName(ctx, c.connectionID, firstName, lastName)
}

// SetBio changes the bio of the connected business account. The bot must have the
// can_edit_bio business bot right.
func (c *BusinessContext) SetBio(ctx context.Context, bio string) error {
	return c.bot.SetBusinessAccountBio(ctx, c.connectionID, bio)
}

// SetUsername changes the username of the connected business account. The bot
// must have the can_edit_username business bot right.
func (c *BusinessContext) SetUsername(ctx context.Context, username string) error {
	return c.bot.SetBusinessAccountUsername(ctx, c.connectionID, username)
}

// SetProfilePhoto changes the profile photo of the connected business account.
// When isPublic is true the photo becomes the public (fallback) photo. The bot
// must have the can_edit_profile_photo business bot right.
func (c *BusinessContext) SetProfilePhoto(ctx context.Context, photo InputProfilePhoto, isPublic bool) error {
	return c.bot.SetBusinessAccountProfilePhoto(ctx, c.connectionID, photo, isPublic)
}

// RemoveProfilePhoto removes the profile photo of the connected business account.
// When isPublic is true the public (fallback) photo is removed. The bot must have
// the can_edit_profile_photo business bot right.
func (c *BusinessContext) RemoveProfilePhoto(ctx context.Context, isPublic bool) error {
	return c.bot.RemoveBusinessAccountProfilePhoto(ctx, c.connectionID, isPublic)
}

// StarBalance returns the Telegram Stars balance of the connected business
// account. The bot must have the can_view_gifts_and_stars business bot right.
func (c *BusinessContext) StarBalance(ctx context.Context) (StarAmount, error) {
	return c.bot.GetBusinessAccountStarBalance(ctx, c.connectionID)
}

// DeleteMessages deletes messages on behalf of the connected business account.
// The messages must all belong to the same chat.
func (c *BusinessContext) DeleteMessages(ctx context.Context, messageIDs []int) error {
	return c.bot.DeleteBusinessMessages(ctx, c.connectionID, messageIDs)
}

// SendMessage sends a text message to a chat on behalf of the connected business
// account and returns the sent message.
func (c *BusinessContext) SendMessage(ctx context.Context, chat ChatID, text string, opts ...SendOption) (*Message, error) {
	return c.bot.SendMessage(ctx, chat, text, append(opts, WithBusinessConnection(c.connectionID))...)
}

// SendPhoto sends a photo to a chat on behalf of the connected business account.
func (c *BusinessContext) SendPhoto(
	ctx context.Context, chat ChatID, photo InputFile, caption string, opts ...SendOption,
) (*Message, error) {
	return c.bot.SendPhoto(ctx, chat, photo, caption, append(opts, WithBusinessConnection(c.connectionID))...)
}

// SendDocument sends a general file to a chat on behalf of the connected business
// account.
func (c *BusinessContext) SendDocument(
	ctx context.Context, chat ChatID, document InputFile, caption string, opts ...SendOption,
) (*Message, error) {
	return c.bot.SendDocument(ctx, chat, document, caption, append(opts, WithBusinessConnection(c.connectionID))...)
}
