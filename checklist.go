package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

// InputChecklistTask describes a task to add to a checklist on creation or edit.
type InputChecklistTask struct {
	// ID is the unique identifier of the task, positive and unique among the
	// tasks of the checklist.
	ID int
	// Text is the task text, 1-100 characters after entities parsing.
	Text string
	// ParseMode is the formatting mode for Text.
	ParseMode ParseMode
	// TextEntities are explicit entities; they take precedence over ParseMode.
	TextEntities []MessageEntity
}

// InputChecklist describes a checklist to create or replace.
type InputChecklist struct {
	// Title is the checklist title, 1-255 characters after entities parsing.
	Title string
	// ParseMode is the formatting mode for Title.
	ParseMode ParseMode
	// TitleEntities are explicit title entities; they take precedence over
	// ParseMode.
	TitleEntities []MessageEntity
	// Tasks are the checklist tasks, 1-30 items.
	Tasks []InputChecklistTask
	// OthersCanAddTasks allows other users to add tasks to the checklist.
	OthersCanAddTasks bool
	// OthersCanMarkTasksAsDone allows other users to mark tasks as done or undone.
	OthersCanMarkTasksAsDone bool
}

// textWithEntities resolves text plus a parse mode or explicit entities into the
// MTProto TextWithEntities shape. Explicit entities take precedence.
func (b *Bot) textWithEntities(ctx context.Context, text string, mode ParseMode, explicit []MessageEntity) (tg.TextWithEntities, error) {
	if len(explicit) > 0 {
		return tg.TextWithEntities{Text: text, Entities: entitiesToTg(explicit)}, nil
	}

	msg, entities, err := b.styledMessage(ctx, text, mode)
	if err != nil {
		return tg.TextWithEntities{}, err
	}

	return tg.TextWithEntities{Text: msg, Entities: entities}, nil
}

// checklistMedia converts an InputChecklist into MTProto todo-list media.
func (b *Bot) checklistMedia(ctx context.Context, c InputChecklist) (*tg.InputMediaTodo, error) {
	if len(c.Tasks) == 0 {
		return nil, &Error{Code: 400, Description: "Bad Request: checklist must include at least one task"}
	}

	title, err := b.textWithEntities(ctx, c.Title, c.ParseMode, c.TitleEntities)
	if err != nil {
		return nil, err
	}

	list := make([]tg.TodoItem, len(c.Tasks))

	for i, task := range c.Tasks {
		taskTitle, err := b.textWithEntities(ctx, task.Text, task.ParseMode, task.TextEntities)
		if err != nil {
			return nil, err
		}

		list[i] = tg.TodoItem{ID: task.ID, Title: taskTitle}
	}

	return &tg.InputMediaTodo{Todo: tg.TodoList{
		OthersCanAppend:   c.OthersCanAddTasks,
		OthersCanComplete: c.OthersCanMarkTasksAsDone,
		Title:             title,
		List:              list,
	}}, nil
}

// SendChecklist sends a checklist on behalf of a connected business account.
func (b *Bot) SendChecklist(
	ctx context.Context, businessConnectionID string, chat ChatID, checklist InputChecklist, opts ...SendOption,
) (*Message, error) {
	opts = append([]SendOption{WithBusinessConnection(businessConnectionID)}, opts...)

	return b.sendResolvedMedia(ctx, chat, "", func(ctx context.Context, _ []styling.StyledTextOption) (message.MediaOption, error) {
		todo, err := b.checklistMedia(ctx, checklist)
		if err != nil {
			return nil, err
		}

		return message.Media(todo), nil
	}, opts...)
}

// EditMessageChecklist replaces the checklist in a message sent on behalf of a
// connected business account.
func (b *Bot) EditMessageChecklist(
	ctx context.Context, businessConnectionID string, chat ChatID, messageID int, checklist InputChecklist, opts ...SendOption,
) (*Message, error) {
	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	todo, err := b.checklistMedia(ctx, checklist)
	if err != nil {
		return nil, err
	}

	req := &tg.MessagesEditMessageRequest{Peer: peer, ID: messageID}
	req.SetMedia(todo)

	if cfg.markup != nil {
		mkp, err := replyMarkupToTg(cfg.markup)
		if err != nil {
			return nil, err
		}

		req.SetReplyMarkup(mkp)
	}

	resp, err := b.businessRaw(businessConnectionID).MessagesEditMessage(ctx, req)

	return b.sentMessage(ctx, peer, resp, err)
}
