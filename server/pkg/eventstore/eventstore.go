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
)

type jsonLineDS struct {
	Filename string
	File     *os.File
}

type line struct {
	Event json.RawMessage
	// TODO: Add number and timestamp
}

// New returns a eventstore instance.
func New(filename string) (*jsonLineDS, func() error, error) {
	f, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("opening datastore file: %w", err)
	}

	ds := jsonLineDS{
		Filename: filename,
		File:     f,
	}
	return &ds, f.Close, nil
}

func (ds *jsonLineDS) Save(event json.RawMessage) error {
	if !json.Valid(event) {
		return fmt.Errorf("invalid JSON encoding for event %q", string(event))
	}

	l := line{
		Event: event,
	}
	encodedLine, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("marshalling JSON line: %w", err)
	}
	encodedLine = append(encodedLine, '\n')

	if _, err := ds.File.Write(encodedLine); err != nil {
		return fmt.Errorf("writing to database file: %w", err)
	}

	return nil
}

func (ds *jsonLineDS) Retrieve(n int) ([]json.RawMessage, error) {
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

	return events, nil
}

// func (db *jsonLineDB) InsertCase(fields json.RawMessage) (int, error) {
// 	nextID := db.maxCaseID() + 1

// 	blob := caseBlob{
// 		ID:     nextID,
// 		Fields: fields,
// 	}
// 	encodedBlob, err := json.Marshal(blob)
// 	if err != nil {
// 		return 0, fmt.Errorf("marshalling case blob for JSON line: %w", err)
// 	}

// 	l := line{
// 		Name: "Case",
// 		Blob: encodedBlob,
// 	}
// 	encodedLine, err := json.Marshal(l)
// 	if err != nil {
// 		return 0, fmt.Errorf("marshalling JSON line: %w", err)
// 	}

// 	encodedLine = append(encodedLine, '\n')
// 	if _, err := db.File.Write(encodedLine); err != nil {
// 		return 0, fmt.Errorf("writing to database file: %w", err)
// 	}

// 	db.Data.Case[nextID] = fields

// 	return nextID, nil
// }

// func (db *jsonLineDB) maxCaseID() int {
// 	var result int
// 	for result = range db.Data.Case {
// 		break
// 	}
// 	for n := range db.Data.Case {
// 		if n > result {
// 			result = n
// 		}
// 	}
// 	return result
// }

// func (db *jsonLineDB) RetrieveCase(id int) (json.RawMessage, error) {
// 	fields, ok := db.Data.Case[id]
// 	if !ok {
// 		return json.RawMessage{}, errorNotFound("Case", id)
// 	}
// 	return fields, nil
// }

// func (db *jsonLineDB) UpdateCase(id int, fields json.RawMessage) error {
// 	if _, err := db.RetrieveCase(id); err != nil {
// 		return fmt.Errorf("preparing update: %w", err)
// 	}
// 	db.Data.Case[id] = fields
// 	return nil
// }

// func errorNotFound(name string, id int) error {
// 	return fmt.Errorf("object %s with id %d was not found", name, id)
// }

// func (db *jsonLineDB) RetrieveCaseAll() (map[int]json.RawMessage, error) {
// 	return db.Data.Case, nil
// }
