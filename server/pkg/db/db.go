/*
Package db provides the implementation of a JSON lines database (data file).
*/
package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type jsonLineDB struct {
	Data dbData
	File *os.File
}

type dbData struct {
	Case map[int]json.RawMessage
}

type line struct {
	Name string
	Blob json.RawMessage
}

type caseBlob struct {
	ID     int
	Fields json.RawMessage
}

// New returns a database instance. The given file is loaded.
func New(f *os.File) (*jsonLineDB, error) {
	data := dbData{
		Case: map[int]json.RawMessage{},
	}
	db := jsonLineDB{
		Data: data,
		File: f,
	}
	if err := db.load(); err != nil {
		return nil, fmt.Errorf("loading data from database file %q, %w", f.Name(), err)
	}
	return &db, nil
}

func (db *jsonLineDB) load() error {
	// See https://github.com/ostcar/timer/blob/main/model/model.go#L75

	s := bufio.NewScanner(db.File)
	for s.Scan() {
		l := line{}
		json.Unmarshal(s.Bytes(), &l)

		switch l.Name {

		case "Case":
			blob := caseBlob{}
			if err := json.Unmarshal(l.Blob, &blob); err != nil {
				return fmt.Errorf("unmarshalling JSON case blob: %w", err)
			}
			db.Data.Case[blob.ID] = blob.Fields

		default:
			return fmt.Errorf("invalid line in database: %q", s.Text())
		}
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("reading database file: %w", err)
	}

	return nil
}

func (db *jsonLineDB) InsertCase(fields json.RawMessage) (int, error) {
	nextID := db.maxCaseID() + 1

	blob := caseBlob{
		ID:     nextID,
		Fields: fields,
	}
	encodedBlob, err := json.Marshal(blob)
	if err != nil {
		return 0, fmt.Errorf("marshalling case blob for JSON line: %w", err)
	}

	l := line{
		Name: "Case",
		Blob: encodedBlob,
	}
	encodedLine, err := json.Marshal(l)
	if err != nil {
		return 0, fmt.Errorf("marshalling JSON line: %w", err)
	}

	encodedLine = append(encodedLine, '\n')
	if _, err := db.File.Write(encodedLine); err != nil {
		return 0, fmt.Errorf("writing to database file: %w", err)
	}

	db.Data.Case[nextID] = fields

	return nextID, nil
}

func (db *jsonLineDB) maxCaseID() int {
	var result int
	for result = range db.Data.Case {
		break
	}
	for n := range db.Data.Case {
		if n > result {
			result = n
		}
	}
	return result
}

func (db *jsonLineDB) RetrieveCase(id int) (json.RawMessage, error) {
	fields, ok := db.Data.Case[id]
	if !ok {
		return json.RawMessage{}, errorNotFound("Case", id)
	}
	return fields, nil
}

func (db *jsonLineDB) UpdateCase(id int, fields json.RawMessage) error {
	if _, err := db.RetrieveCase(id); err != nil {
		return fmt.Errorf("preparing update: %w", err)
	}
	db.Data.Case[id] = fields
	return nil
}

func errorNotFound(name string, id int) error {
	return fmt.Errorf("object %s with id %d was not found", name, id)
}

func (db *jsonLineDB) RetrieveCaseAll() (map[int]json.RawMessage, error) {
	return db.Data.Case, nil
}
