package peers

import (
	"context"

	"go.uber.org/multierr"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

// Storage represents peer storage.
type Storage interface {
	SaveUsers(ctx context.Context, users ...*tg.User) error
	SaveChats(ctx context.Context, chats ...tg.FullChat) error

	FindUser(ctx context.Context, id int64) (*tg.User, bool, error)
	FindChat(ctx context.Context, id int64) (tg.FullChat, bool, error)
}

func save(ctx context.Context, s Storage, from interface {
	MapChats() tg.ChatClassArray
	MapUsers() tg.UserClassArray
}) error {
	return multierr.Append(
		s.SaveChats(ctx, from.MapChats().AppendOnlyFull(nil)...),
		s.SaveUsers(ctx, from.MapUsers().AppendOnlyNotEmpty(nil)...),
	)
}

// UpdateHook is update hook for Storage.
func UpdateHook(s Storage, next telegram.UpdateHandler) telegram.UpdateHandler {
	return telegram.UpdateHandlerFunc(func(ctx context.Context, u tg.UpdatesClass) error {
		var err error
		switch v := u.(type) {
		case *tg.UpdatesCombined:
			err = save(ctx, s, v)
		case *tg.Updates:
			err = save(ctx, s, v)
		}
		return multierr.Append(err, next.Handle(ctx, u))
	})
}

// AccessHasher is implementation of updates.ChannelAccessHasher based on Storage.
type AccessHasher struct {
	Storage Storage
}

// SetChannelAccessHash implements updates.ChannelAccessHasher.
func (a AccessHasher) SetChannelAccessHash(userID, channelID, accessHash int64) error {
	// TODO: update access hash?
	return nil
}

// GetChannelAccessHash implements updates.ChannelAccessHasher.
func (a AccessHasher) GetChannelAccessHash(userID, channelID int64) (accessHash int64, found bool, err error) {
	v, ok, err := a.Storage.FindChat(context.TODO(), channelID)
	if err != nil {
		return 0, false, err
	}
	if !ok {
		return 0, false, nil
	}
	nonForbidden, ok := v.(interface {
		GetAccessHash() int64
	})
	if !ok {
		return 0, false, nil
	}
	return nonForbidden.GetAccessHash(), true, nil
}
