package botapi

import (
	"strconv"
	"sync"

	"github.com/gotd/td/tg"
)

// businessDedupSize is the number of recent business messages remembered for
// deduplication. Telegram may redeliver updates (notably to bots after a
// reconnect), and the qts sequence does not always suppress this; the window
// only needs to cover the burst a redelivery replays.
const businessDedupSize = 4096

// businessDedup remembers recently delivered business messages so a redelivered
// update does not fire handlers — and reply — twice. It is a fixed-size set with
// FIFO eviction: safe for concurrent use, bounded in memory.
type businessDedup struct {
	mu   sync.Mutex
	seen map[string]struct{}
	ring []string
	pos  int
}

func newBusinessDedup(size int) *businessDedup {
	return &businessDedup{
		seen: make(map[string]struct{}, size),
		ring: make([]string, size),
	}
}

// fresh reports whether key has not been seen recently, recording it. A repeated
// key (a redelivered update) returns false.
func (d *businessDedup) fresh(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.seen[key]; ok {
		return false
	}

	if old := d.ring[d.pos]; old != "" {
		delete(d.seen, old)
	}

	d.ring[d.pos] = key
	d.pos = (d.pos + 1) % len(d.ring)
	d.seen[key] = struct{}{}

	return true
}

// businessMessageKey identifies a business message for deduplication. The edit
// date distinguishes a genuine edit (which should be handled) from a redelivery
// of the same message (which should not).
func businessMessageKey(connectionID string, m *tg.Message) string {
	edit, _ := m.GetEditDate()

	return connectionID + ":" + strconv.Itoa(m.ID) + ":" + strconv.Itoa(edit)
}
