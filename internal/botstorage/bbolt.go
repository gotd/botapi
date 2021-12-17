package botstorage

import (
	"context"
	"encoding/binary"

	"github.com/go-faster/errors"
	"github.com/gotd/td/bin"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"go.etcd.io/bbolt"
)

// BBoltStorage is bbolt-based storage.
type BBoltStorage struct {
	db *bbolt.DB
}

var _ interface {
	peers.Storage
	peers.Cache
	updates.StateStorage
} = (*BBoltStorage)(nil)

// NewBBoltStorage creates new BBoltStorage.
func NewBBoltStorage(db *bbolt.DB) *BBoltStorage {
	return &BBoltStorage{db: db}
}

var _ = map[string]struct{}{
	hashPrefix:        {},
	userPrefix:        {},
	userFullPrefix:    {},
	chatPrefix:        {},
	chatFullPrefix:    {},
	channelPrefix:     {},
	channelFullPrefix: {},
	statePrefix:       {},
	channelsPtsPrefix: {},
	sessionPrefix:     {},
}

const (
	hashPrefix        = "hashes_"
	userPrefix        = "users_"
	userFullPrefix    = "userFulls_"
	chatPrefix        = "chats_"
	chatFullPrefix    = "chatFulls_"
	channelPrefix     = "channels_"
	channelFullPrefix = "channelFulls_"
)

func formatInt(i int64) []byte {
	var keyBuf [8]byte
	binary.LittleEndian.PutUint64(keyBuf[:], uint64(i))
	return keyBuf[:]
}

func parseInt(v []byte) (int64, bool) {
	if len(v) < 8 {
		return 0, false
	}
	i := binary.LittleEndian.Uint64(v)
	return int64(i), true
}

func (b *BBoltStorage) viewBucket(prefix string, cb func(b *bbolt.Bucket, tx *bbolt.Tx) error) error {
	return b.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(prefix))
		if b == nil {
			return nil
		}
		return cb(b, tx)
	})
}

func (b *BBoltStorage) batchBucket(prefix string, cb func(b *bbolt.Bucket, tx *bbolt.Tx) error) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(prefix))
		if err != nil {
			return errors.Wrapf(err, "create %q bucket", prefix)
		}
		return cb(b, tx)
	})
}

// Save implements peers.Storage.
func (b *BBoltStorage) Save(_ context.Context, key peers.Key, value peers.Value) error {
	return b.batchBucket(hashPrefix+key.Prefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		return b.Put(formatInt(key.ID), formatInt(value.AccessHash))
	})
}

// Find implements peers.Storage.
func (b *BBoltStorage) Find(_ context.Context, key peers.Key) (value peers.Value, found bool, err error) {
	err = b.viewBucket(hashPrefix+key.Prefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		storageKey := formatInt(key.ID)
		val := b.Get(storageKey)
		// Value not found.
		if val == nil {
			return nil
		}
		id, ok := parseInt(val)
		if !ok {
			return errors.Errorf("got invalid value %+x", val)
		}
		value.AccessHash = id
		found = true
		return nil
	})

	return value, found, err
}

// SavePhone implements peers.Storage.
func (b *BBoltStorage) SavePhone(_ context.Context, phone string, key peers.Key) error {
	// FIXME(tdakkota): Implement it. We don't need it for bots
	//  However, we can use it as default cache.
	return nil
}

// FindPhone implements peers.Storage.
func (b *BBoltStorage) FindPhone(_ context.Context, phone string) (key peers.Key, value peers.Value, found bool, err error) {
	// FIXME(tdakkota): Implement it. We don't need it for bots
	//  However, we can use it as default cache.
	return key, value, found, err
}

// GetContactsHash implements peers.Storage.
func (b *BBoltStorage) GetContactsHash(ctx context.Context) (int64, error) {
	// FIXME(tdakkota): Implement it. We don't need it for bots
	//  However, we can use it as default cache.
	return 0, nil
}

// SaveContactsHash implements peers.Storage.
func (b *BBoltStorage) SaveContactsHash(_ context.Context, hash int64) error {
	// FIXME(tdakkota): Implement it. We don't need it for bots
	//  However, we can use it as default cache.
	return nil
}

func putMTProtoKey(b *bbolt.Bucket, e interface {
	GetID() int64
	bin.Encoder
}) error {
	var buf bin.Buffer
	if err := e.Encode(&buf); err != nil {
		return errors.Wrap(err, "encode")
	}
	id := e.GetID()
	if err := b.Put(formatInt(id), buf.Raw()); err != nil {
		return errors.Wrapf(err, "put %d", id)
	}
	return nil
}

func getMTProtoKey(b *bbolt.Bucket, id int64, d bin.Decoder) (bool, error) {
	key := formatInt(id)
	data := b.Get(key)
	if data == nil {
		return false, nil
	}
	buf := bin.Buffer{Buf: data}
	if err := d.Decode(&buf); err != nil {
		// Ignore decode errors.
		return false, nil
	}
	return true, nil
}

// SaveUsers implements BBoltStorage.
func (b *BBoltStorage) SaveUsers(_ context.Context, users ...*tg.User) error {
	return b.batchBucket(userPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		for _, user := range users {
			if err := putMTProtoKey(b, user); err != nil {
				return err
			}
		}
		return nil
	})
}

// SaveUserFulls implements BBoltStorage.
func (b *BBoltStorage) SaveUserFulls(_ context.Context, users ...*tg.UserFull) error {
	return b.batchBucket(userFullPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		for _, user := range users {
			if err := putMTProtoKey(b, user); err != nil {
				return err
			}
		}
		return nil
	})
}

// FindUser implements BBoltStorage.
func (b *BBoltStorage) FindUser(_ context.Context, id int64) (e *tg.User, found bool, err error) {
	// Use batch to delete invalid keys.
	err = b.viewBucket(userPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		e = new(tg.User)
		found, err = getMTProtoKey(b, id, e)
		return err
	})
	return e, found, err
}

// FindUserFull implements BBoltStorage.
func (b *BBoltStorage) FindUserFull(_ context.Context, id int64) (e *tg.UserFull, found bool, err error) {
	// Use batch to delete invalid keys.
	err = b.viewBucket(userFullPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		e = new(tg.UserFull)
		found, err = getMTProtoKey(b, id, e)
		return err
	})
	return e, found, err
}

// SaveChats implements BBoltStorage.
func (b *BBoltStorage) SaveChats(_ context.Context, chats ...*tg.Chat) error {
	return b.batchBucket(chatPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		for _, chat := range chats {
			if err := putMTProtoKey(b, chat); err != nil {
				return err
			}
		}
		return nil
	})
}

// SaveChatFulls implements BBoltStorage.
func (b *BBoltStorage) SaveChatFulls(_ context.Context, chats ...*tg.ChatFull) error {
	return b.batchBucket(chatFullPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		for _, chat := range chats {
			if err := putMTProtoKey(b, chat); err != nil {
				return err
			}
		}
		return nil
	})
}

// FindChat implements BBoltStorage.
func (b *BBoltStorage) FindChat(_ context.Context, id int64) (e *tg.Chat, found bool, err error) {
	err = b.viewBucket(chatPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		e = new(tg.Chat)
		found, err = getMTProtoKey(b, id, e)
		return err
	})
	return e, found, err
}

// FindChatFull implements BBoltStorage.
func (b *BBoltStorage) FindChatFull(_ context.Context, id int64) (e *tg.ChatFull, found bool, err error) {
	err = b.viewBucket(chatFullPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		e = new(tg.ChatFull)
		found, err = getMTProtoKey(b, id, e)
		return err
	})
	return e, found, err
}

// SaveChannels implements BBoltStorage.
func (b *BBoltStorage) SaveChannels(_ context.Context, channels ...*tg.Channel) error {
	return b.batchBucket(channelPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		for _, channel := range channels {
			if err := putMTProtoKey(b, channel); err != nil {
				return err
			}
		}
		return nil
	})
}

// SaveChannelFulls implements BBoltStorage.
func (b *BBoltStorage) SaveChannelFulls(_ context.Context, channels ...*tg.ChannelFull) error {
	return b.batchBucket(channelFullPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		for _, channel := range channels {
			if err := putMTProtoKey(b, channel); err != nil {
				return err
			}
		}
		return nil
	})
}

// FindChannel implements BBoltStorage.
func (b *BBoltStorage) FindChannel(_ context.Context, id int64) (e *tg.Channel, found bool, err error) {
	err = b.viewBucket(channelPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		e = new(tg.Channel)
		found, err = getMTProtoKey(b, id, e)
		return err
	})
	return e, found, err
}

// FindChannelFull implements BBoltStorage.
func (b *BBoltStorage) FindChannelFull(_ context.Context, id int64) (e *tg.ChannelFull, found bool, err error) {
	err = b.viewBucket(channelFullPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		e = new(tg.ChannelFull)
		found, err = getMTProtoKey(b, id, e)
		return err
	})
	return e, found, err
}

const (
	statePrefix = "state_"
	ptsKey      = "pts"
	qtsKey      = "qts"
	dateKey     = "date"
	seqKey      = "seq"

	channelsPtsPrefix = "channel_pts_"
)

func getStateField(b *bbolt.Bucket, key string, v *int) (bool, error) {
	bytesKey := []byte(key)
	data := b.Get(bytesKey)
	if data == nil {
		return false, nil
	}
	p, ok := parseInt(data)
	if !ok {
		return false, errors.Errorf("decode %q", key)
	}
	*v = int(p)
	return true, nil
}

// GetState implements updates.StateStorage.
func (b *BBoltStorage) GetState(_ int64) (state updates.State, found bool, err error) {
	err = b.viewBucket(statePrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		if ok, err := getStateField(b, ptsKey, &state.Pts); err != nil || !ok {
			return err
		}
		if ok, err := getStateField(b, qtsKey, &state.Qts); err != nil || !ok {
			return err
		}
		if ok, err := getStateField(b, dateKey, &state.Date); err != nil || !ok {
			return err
		}
		if ok, err := getStateField(b, seqKey, &state.Seq); err != nil || !ok {
			return err
		}
		found = true
		return nil
	})
	return state, found, err
}

func setStateField(b *bbolt.Bucket, key string, v int) error {
	return b.Put([]byte(key), formatInt(int64(v)))
}

// SetState implements updates.StateStorage.
func (b *BBoltStorage) SetState(_ int64, state updates.State) error {
	return b.batchBucket(statePrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		if err := setStateField(b, ptsKey, state.Pts); err != nil {
			return err
		}
		if err := setStateField(b, qtsKey, state.Qts); err != nil {
			return err
		}
		if err := setStateField(b, dateKey, state.Date); err != nil {
			return err
		}
		if err := setStateField(b, seqKey, state.Seq); err != nil {
			return err
		}
		return nil
	})
}

// SetPts implements updates.StateStorage.
func (b *BBoltStorage) SetPts(_ int64, pts int) error {
	return b.batchBucket(statePrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		if err := setStateField(b, ptsKey, pts); err != nil {
			return err
		}
		return nil
	})
}

// SetQts implements updates.StateStorage.
func (b *BBoltStorage) SetQts(_ int64, qts int) error {
	return b.batchBucket(statePrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		if err := setStateField(b, qtsKey, qts); err != nil {
			return err
		}
		return nil
	})
}

// SetDate implements updates.StateStorage.
func (b *BBoltStorage) SetDate(_ int64, date int) error {
	return b.batchBucket(statePrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		if err := setStateField(b, dateKey, date); err != nil {
			return err
		}
		return nil
	})
}

// SetSeq implements updates.StateStorage.
func (b *BBoltStorage) SetSeq(_ int64, seq int) error {
	return b.batchBucket(statePrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		if err := setStateField(b, seqKey, seq); err != nil {
			return err
		}
		return nil
	})
}

// SetDateSeq implements updates.StateStorage.
func (b *BBoltStorage) SetDateSeq(_ int64, date, seq int) error {
	return b.batchBucket(statePrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		if err := setStateField(b, dateKey, date); err != nil {
			return err
		}
		if err := setStateField(b, seqKey, seq); err != nil {
			return err
		}
		return nil
	})
}

// GetChannelPts implements updates.StateStorage.
func (b *BBoltStorage) GetChannelPts(_, channelID int64) (pts int, found bool, err error) {
	err = b.viewBucket(channelsPtsPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		v, ok := parseInt(b.Get(formatInt(channelID)))
		pts, found = int(v), ok
		return nil
	})
	return pts, found, err
}

// SetChannelPts implements updates.StateStorage.
func (b *BBoltStorage) SetChannelPts(_, channelID int64, pts int) error {
	return b.batchBucket(channelsPtsPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		return b.Put(formatInt(channelID), formatInt(int64(pts)))
	})
}

// ForEachChannels implements updates.StateStorage.
func (b *BBoltStorage) ForEachChannels(_ int64, f func(channelID int64, pts int) error) error {
	return b.viewBucket(channelsPtsPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		return b.ForEach(func(k, v []byte) error {
			channelID, ok := parseInt(k)
			if !ok {
				// Ignore invalid entries.
				return nil
			}
			pts, ok := parseInt(v)
			if !ok {
				// Ignore invalid entries.
				return nil
			}
			return f(channelID, int(pts))
		})
	})
}

const (
	sessionPrefix = "session_"
	sessionKey    = "session"
)

// LoadSession implements session.Storage.
func (b *BBoltStorage) LoadSession(_ context.Context) (data []byte, err error) {
	err = b.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(sessionPrefix))
		if b == nil {
			return session.ErrNotFound
		}

		data = b.Get([]byte(sessionKey))
		if data == nil {
			return session.ErrNotFound
		}
		return nil
	})
	return data, err
}

// StoreSession implements session.Storage.
func (b *BBoltStorage) StoreSession(_ context.Context, data []byte) error {
	return b.batchBucket(sessionPrefix, func(b *bbolt.Bucket, tx *bbolt.Tx) error {
		return b.Put([]byte(sessionKey), data)
	})
}
