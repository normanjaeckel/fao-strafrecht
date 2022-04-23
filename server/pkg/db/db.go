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

type entry struct {
	Name string
	ID   int
	Data json.RawMessage
}

type dbData map[string]map[int]json.RawMessage

// New returns a database instance. The given file is loaded.
func New(f *os.File) (*jsonLineDB, error) {
	data, err := load(f)
	if err != nil {
		return nil, fmt.Errorf("loading data from database file %q, %w", f.Name(), err)
	}
	return &jsonLineDB{
		Data: data,
		File: f,
	}, nil
}

func load(f *os.File) (dbData, error) {
	// See https://github.com/ostcar/timer/blob/main/model/model.go#L75

	dbData := dbData{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		e := entry{}
		json.Unmarshal(s.Bytes(), &e)

		objs := dbData[e.Name]
		if objs == nil {
			dbData[e.Name] = map[int]json.RawMessage{}
		}

		dbData[e.Name][e.ID] = e.Data
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("reading database file: %w", err)
	}

	return dbData, nil
}

func (db *jsonLineDB) Insert(name string, data json.RawMessage) (int, error) {
	nextID := db.maxID(name) + 1

	objs := db.Data[name]
	if objs == nil {
		db.Data[name] = map[int]json.RawMessage{}
	}

	db.Data[name][nextID] = data

	e := entry{
		Name: name,
		ID:   nextID,
		Data: data,
	}
	line, err := json.Marshal(e)
	if err != nil {
		return 0, fmt.Errorf("marshalling JSON line: %w", err)
	}

	line = append(line, '\n')
	if _, err := db.File.Write(line); err != nil {
		return 0, fmt.Errorf("writing to database file: %w", err)
	}

	return nextID, nil
}

func (db *jsonLineDB) maxID(name string) int {
	objs := db.Data[name]

	var result int
	for result = range objs {
		break
	}
	for n := range objs {
		if n > result {
			result = n
		}
	}
	return result
}

func (db *jsonLineDB) Retrieve(name string, id int) (json.RawMessage, error) {
	objs := db.Data[name]
	if objs == nil {
		return json.RawMessage{}, errorNotFound(name, id)
	}

	data, ok := objs[id]
	if !ok {
		return json.RawMessage{}, errorNotFound(name, id)
	}

	return data, nil
}

func (db *jsonLineDB) Update(name string, id int, data json.RawMessage) error {
	objs := db.Data[name]
	if objs == nil {
		return errorNotFound(name, id)
	}

	_, ok := objs[id]
	if !ok {
		return errorNotFound(name, id)
	}

	db.Data[name][id] = data
	return nil
}

func errorNotFound(name string, id int) error {
	return fmt.Errorf("object %s with id %d was not found", name, id)
}

func (db *jsonLineDB) RetrieveAll(name string) (map[int]json.RawMessage, error) {
	objs := db.Data[name]
	if objs == nil {
		return map[int]json.RawMessage{}, nil
	}

	return objs, nil
}
