// Package audit provides functionality to record and retrieve
// a history of operations performed on scopes within the store.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single auditable event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Scope     string    `json:"scope"`
	Action    string    `json:"action"`
	Detail    string    `json:"detail,omitempty"`
}

const logFileName = "audit.log"

// Log appends an audit entry to the store's audit log file.
func Log(storeDir, scope, action, detail string) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Scope:     scope,
		Action:    action,
		Detail:    detail,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	data = append(data, '\n')

	path := filepath.Join(storeDir, logFileName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}

// Read returns all audit entries from the store's audit log.
// Returns an empty slice if the log does not yet exist.
func Read(storeDir string) ([]Entry, error) {
	path := filepath.Join(storeDir, logFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read log: %w", err)
	}

	var entries []Entry
	dec := json.NewDecoder(bytesReader(data))
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// bytesReader wraps a byte slice for use with json.NewDecoder.
func bytesReader(b []byte) *bytesReaderImpl {
	return &bytesReaderImpl{data: b}
}

type bytesReaderImpl struct {
	data []byte
	pos  int
}

func (r *bytesReaderImpl) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("EOF")
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
