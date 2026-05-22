// Package history tracks per-scope change history by recording a lightweight
// journal of mutations (add, update, remove) with timestamps and optional
// actor labels. Entries are appended to a newline-delimited JSON file stored
// alongside the encrypted scope file.
package history

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ErrNoHistory is returned when no history file exists for the requested scope.
var ErrNoHistory = errors.New("history: no history found for scope")

// Op describes the kind of mutation that was recorded.
type Op string

const (
	OpAdd    Op = "add"
	OpUpdate Op = "update"
	OpRemove Op = "remove"
)

// Entry is a single history record for one key mutation.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Op        Op        `json:"op"`
	Key       string    `json:"key"`
	Actor     string    `json:"actor,omitempty"`
}

// historyPath returns the path to the history journal for the given scope.
func historyPath(storeDir, scope string) string {
	return filepath.Join(storeDir, scope+".history.jsonl")
}

// Record appends a new Entry to the history journal for the given scope.
// The journal is created if it does not already exist.
func Record(storeDir, scope string, op Op, key, actor string) error {
	if err := os.MkdirAll(storeDir, 0o700); err != nil {
		return err
	}

	f, err := os.OpenFile(historyPath(storeDir, scope), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Op:        op,
		Key:       key,
		Actor:     actor,
	}

	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	line = append(line, '\n')

	_, err = f.Write(line)
	return err
}

// Read returns all history entries for the given scope in chronological order.
// Returns ErrNoHistory if no journal file exists yet.
func Read(storeDir, scope string) ([]Entry, error) {
	f, err := os.Open(historyPath(storeDir, scope))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoHistory
		}
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Clear removes the history journal for the given scope.
// It is a no-op if no journal exists.
func Clear(storeDir, scope string) error {
	err := os.Remove(historyPath(storeDir, scope))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
