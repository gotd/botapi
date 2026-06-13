package botapi

// Group is a set of handlers that share predicates and middleware. Every handler
// registered through the group inherits the group's predicates (as additional
// guards) and its middleware (applied inside the global middleware).
//
// Create one with Bot.Group:
//
//	admin := bot.Group(ChatTypeIs(ChatTypeSupergroup))
//	admin.Use(Recover())
//	admin.OnCommand("ban", banHandler)
type Group struct {
	bot   *Bot
	preds []Predicate
	mws   []Middleware
}

// Group returns a new handler group guarded by the given predicates.
func (b *Bot) Group(predicates ...Predicate) *Group {
	return &Group{bot: b, preds: predicates}
}

// Use adds middleware applied to every handler registered through this group.
// Returns the group for chaining.
func (g *Group) Use(mws ...Middleware) *Group {
	g.mws = append(g.mws, mws...)
	return g
}

// Handle registers a handler in the group with optional extra predicates.
func (g *Group) Handle(h Handler, predicates ...Predicate) {
	preds := make([]Predicate, 0, len(g.preds)+len(predicates))
	preds = append(preds, g.preds...)
	preds = append(preds, predicates...)
	g.bot.onWith(h, g.mws, preds)
}

// OnMessage registers a message handler in the group.
func (g *Group) OnMessage(h Handler, predicates ...Predicate) {
	g.Handle(h, prepend(hasMessage, predicates)...)
}

// OnCallbackQuery registers a callback-query handler in the group.
func (g *Group) OnCallbackQuery(h Handler, predicates ...Predicate) {
	g.Handle(h, prepend(hasCallbackQuery, predicates)...)
}

// OnCommand registers a command handler in the group.
func (g *Group) OnCommand(name string, h Handler, predicates ...Predicate) {
	g.OnMessage(h, prepend(Command(name), predicates)...)
}
