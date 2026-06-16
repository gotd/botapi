package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func sampleChecklist() InputChecklist {
	return InputChecklist{
		Title:             "Groceries",
		Tasks:             []InputChecklistTask{{ID: 1, Text: "Milk"}, {ID: 2, Text: "Bread"}},
		OthersCanAddTasks: true,
	}
}

func TestSendChecklist(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)

	if _, err := b.SendChecklist(context.Background(), "bc1", userRef(10, 20), sampleChecklist()); err != nil {
		t.Fatalf("SendChecklist: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesSendMediaRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	sm, ok := wrapper.Query.(*tg.MessagesSendMediaRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	todo, ok := sm.Media.(*tg.InputMediaTodo)
	if !ok {
		t.Fatalf("media = %#v, want todo", sm.Media)
	}

	if todo.Todo.Title.Text != "Groceries" || len(todo.Todo.List) != 2 || !todo.Todo.OthersCanAppend {
		t.Fatalf("todo = %#v", todo.Todo)
	}

	if todo.Todo.List[0].ID != 1 || todo.Todo.List[0].Title.Text != "Milk" {
		t.Fatalf("task 0 = %#v", todo.Todo.List[0])
	}
}

func TestEditMessageChecklist(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)

	if _, err := b.EditMessageChecklist(context.Background(), "bc2", userRef(10, 20), 5, sampleChecklist()); err != nil {
		t.Fatalf("EditMessageChecklist: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesEditMessageRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc2" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	em, ok := wrapper.Query.(*tg.MessagesEditMessageRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if em.ID != 5 {
		t.Fatalf("message id = %d", em.ID)
	}

	media, ok := em.GetMedia()
	if !ok {
		t.Fatal("edit should set media")
	}

	if _, ok := media.(*tg.InputMediaTodo); !ok {
		t.Fatalf("media = %#v, want todo", media)
	}
}

func TestChecklistEmptyTasks(t *testing.T) {
	inv := newMockInvoker()

	if _, err := newMockBot(inv).SendChecklist(context.Background(), "bc1", userRef(10, 20),
		InputChecklist{Title: "x"}); err == nil {
		t.Fatal("expected error for checklist without tasks")
	}
}
