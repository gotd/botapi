package peers

import (
	"context"
	"sync"

	"github.com/gotd/td/tg"
)

// InmemoryStorage stores users and chats info in memory.
type InmemoryStorage struct {
	chats    map[int64]tg.FullChat
	chatsMux sync.RWMutex

	usersMux sync.RWMutex
	users    map[int64]*tg.User
}

// NewInmemoryStorage creates new InmemoryStorage.
func NewInmemoryStorage() *InmemoryStorage {
	return &InmemoryStorage{
		chats: map[int64]tg.FullChat{},
		users: map[int64]*tg.User{},
	}
}

// SaveUsers implements FileStorage.
func (f *InmemoryStorage) SaveUsers(ctx context.Context, users ...*tg.User) error {
	f.usersMux.Lock()
	defer f.usersMux.Unlock()

	for _, u := range users {
		f.users[u.GetID()] = u
	}

	return nil
}

// SaveChats implements InmemoryStorage.
func (f *InmemoryStorage) SaveChats(ctx context.Context, chats ...tg.FullChat) error {
	f.chatsMux.Lock()
	defer f.chatsMux.Unlock()

	for _, u := range chats {
		f.chats[u.GetID()] = u
	}

	return nil
}

// FindUser implements InmemoryStorage.
func (f *InmemoryStorage) FindUser(ctx context.Context, id int64) (*tg.User, bool, error) {
	f.usersMux.RLock()
	defer f.usersMux.RUnlock()

	v, ok := f.users[id]
	return v, ok, nil
}

// FindChat implements InmemoryStorage.
func (f *InmemoryStorage) FindChat(ctx context.Context, id int64) (tg.FullChat, bool, error) {
	f.chatsMux.RLock()
	defer f.chatsMux.RUnlock()

	v, ok := f.chats[id]
	return v, ok, nil
}
