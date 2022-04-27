/*
Package eventstore provides the implementation of a JSON lines datastore (data
file) for all events.
*/
package eventstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger interface {
	Printf(format string, v ...any)
}

type jsonLineDS struct {
	Logger   Logger
	Filename string
	File     *os.File
}

type line struct {
	Event     json.RawMessage
	Timestamp int64
}

// New returns a eventstore instance.
func New(logger Logger, filename string) (*jsonLineDS, func() error, error) {
	f, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("opening datastore file: %w", err)
	}

	logger.Printf("Opened datastore file %s", filename)

	ds := jsonLineDS{
		Logger:   logger,
		Filename: filename,
		File:     f,
	}
	return &ds, f.Close, nil
}

// Save writes one event into the eventstore
func (ds *jsonLineDS) Save(event json.RawMessage) error {
	if !json.Valid(event) {
		return fmt.Errorf("invalid JSON encoding for event %q", string(event))
	}

	l := line{
		Event:     event,
		Timestamp: time.Now().Unix(),
	}
	encodedLine, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("marshalling JSON line: %w", err)
	}
	encodedLine = append(encodedLine, '\n')

	if _, err := ds.File.Write(encodedLine); err != nil {
		return fmt.Errorf("writing to database file: %w", err)
	}

	ds.Logger.Printf("Wrote event to datastore file: %s", string(event))

	return nil
}

// Retrieve retrieves all events from the eventstore.
func (ds *jsonLineDS) Retrieve() ([]json.RawMessage, error) {
	f, err := os.Open(ds.Filename)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", ds.Filename, err)
	}
	defer f.Close()

	var events []json.RawMessage

	s := bufio.NewScanner(f)
	for s.Scan() {
		l := line{}
		json.Unmarshal(s.Bytes(), &l)
		events = append(events, l.Event)
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("reading database file: %w", err)
	}

	ds.Logger.Printf("Retrieved all events from datastore file")

	return events, nil
}
