package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	Sequence     int       `json:"sequence"`
	Timestamp    time.Time `json:"timestamp"`
	Actor        string    `json:"actor"`
	Action       string    `json:"action"`
	Payload      string    `json:"payload"`
	PreviousHash string    `json:"previous_hash"`
	Hash         string    `json:"hash"`
}
type Trail struct {
	mu      sync.RWMutex
	entries []Entry
}

func (t *Trail) Append(actor, action string, payload any) (Entry, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return Entry{}, fmt.Errorf("marshal audit payload: %w", err)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	prev := ""
	if len(t.entries) > 0 {
		prev = t.entries[len(t.entries)-1].Hash
	}
	e := Entry{Sequence: len(t.entries) + 1, Timestamp: time.Now().UTC(), Actor: actor, Action: action, Payload: string(b), PreviousHash: prev}
	sum := sha256.Sum256([]byte(fmt.Sprintf("%d|%s|%s|%s|%s|%s", e.Sequence, e.Timestamp, e.Actor, e.Action, e.Payload, prev)))
	e.Hash = hex.EncodeToString(sum[:])
	t.entries = append(t.entries, e)
	return e, nil
}
func (t *Trail) Verify() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	prev := ""
	for i, e := range t.entries {
		if e.Sequence != i+1 || e.PreviousHash != prev {
			return false
		}
		sum := sha256.Sum256([]byte(fmt.Sprintf("%d|%s|%s|%s|%s|%s", e.Sequence, e.Timestamp, e.Actor, e.Action, e.Payload, e.PreviousHash)))
		if e.Hash != hex.EncodeToString(sum[:]) {
			return false
		}
		prev = e.Hash
	}
	return true
}
func (t *Trail) Entries() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return append([]Entry(nil), t.entries...)
}
