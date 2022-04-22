/*
Package db provides the implementation of a JSON lines database (data file).
*/
package db

import (
	"encoding/json"
	"fmt"
)

type JSONLineDB struct {
	Data map[string]map[int]json.RawMessage
}

// New returns a new database instance.
func New() *JSONLineDB {
	data := map[string]map[int]json.RawMessage{}
	return &JSONLineDB{Data: data}
}

func (db *JSONLineDB) Insert(name string, data json.RawMessage) (int, error) {
	nextID := db.maxID(name) + 1
	objs := db.Data[name]
	if objs == nil {
		db.Data[name] = map[int]json.RawMessage{}
	}
	db.Data[name][nextID] = data
	return nextID, nil
}

func (db *JSONLineDB) maxID(name string) int {
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

func (db *JSONLineDB) Retrieve(name string, id int) (json.RawMessage, error) {
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

func (db *JSONLineDB) Update(name string, id int, data json.RawMessage) error {
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

func (db *JSONLineDB) RetrieveAll(name string) (map[int]json.RawMessage, error) {
	objs := db.Data[name]
	if objs == nil {
		return map[int]json.RawMessage{}, nil
	}
	return objs, nil
}
