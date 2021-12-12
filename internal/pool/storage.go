package pool

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/go-faster/errors"

	"github.com/gotd/td/session"
)

type fileStorage struct {
	path string
	mux  sync.Mutex
}

type sessionFile struct {
	Data map[string][]byte `json:"data"`
}

func (s *fileStorage) Store(ctx context.Context, id string, data []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	var decoded sessionFile

	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) || len(b) == 0 {
		// Blank initial session.
	} else if err == nil {
		if err := json.Unmarshal(b, &decoded); err != nil {
			return errors.Wrap(err, "unmarshal session file")
		}
	}
	if decoded.Data == nil {
		decoded.Data = map[string][]byte{}
	}

	decoded.Data[id] = data

	if b, err = json.Marshal(&decoded); err != nil {
		return err
	}

	return os.WriteFile(s.path, b, 0o600)
}

func (s *fileStorage) Load(ctx context.Context, id string) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) || len(data) == 0 {
		return nil, session.ErrNotFound
	}

	var decoded sessionFile
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, err
	}

	if len(decoded.Data) == 0 {
		return nil, session.ErrNotFound
	}

	return decoded.Data[id], nil
}

type clientStorage struct {
	storage StateStorage
	id      string
}

func (c clientStorage) LoadSession(ctx context.Context) ([]byte, error) {
	data, err := c.storage.Load(ctx, c.id)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, session.ErrNotFound
	}
	return data, nil
}

func (c clientStorage) StoreSession(ctx context.Context, data []byte) error {
	return c.storage.Store(ctx, c.id, data)
}
