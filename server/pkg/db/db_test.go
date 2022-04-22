package db_test

import (
	"encoding/json"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/db"
)

type SomeData struct {
	Foo string
	Bar string
}

func TestDB(t *testing.T) {
	db := db.New()
	name := "my collection"
	firstData := SomeData{
		Foo: "foo",
		Bar: "bar",
	}

	t.Run("simple insert and retrieve", func(t *testing.T) {
		id := insert(t, db, name, firstData)

		expected := 1
		if id != expected {
			t.Fatalf("wrong id, expected %d, got %d", expected, id)
		}

		result := SomeData{}
		retrieve(t, db, name, id, &result)

		if firstData != result {
			t.Fatalf("wrong retrieved data, expected %v, got %v", firstData, result)
		}
	})

	t.Run("second insert", func(t *testing.T) {
		data := SomeData{
			Foo: "foo2",
		}

		id := insert(t, db, name, data)

		expected := 2
		if id != expected {
			t.Fatalf("wrong id, expected %d, got %d", expected, id)
		}
	})

	t.Run("update", func(t *testing.T) {
		newData := SomeData{
			Foo: "foo3",
		}
		update(t, db, name, 1, newData)

		result := SomeData{}
		retrieve(t, db, name, 1, &result)

		expected := SomeData{
			Foo: newData.Foo,
			Bar: "",
		}

		if expected != result {
			t.Fatalf("wrong retrieved data, expected %v, got %v", expected, result)
		}
	})

	t.Run("retriev all", func(t *testing.T) {
		result, err := db.RetrieveAll(name)
		if err != nil {
			t.Fatalf("retrieving all data: %v", err)
		}
		expectedLen := 2
		got := len(result)
		if expectedLen != got {
			t.Fatalf("wrong length of result, expected %d, got %d", expectedLen, got)
		}
	})
}

func insert(t testing.TB, db *db.JSONLineDB, name string, data interface{}) int {
	t.Helper()
	encodedData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshalling json: %v", err)
	}
	id, err := db.Insert(name, encodedData)
	if err != nil {
		t.Fatalf("inserting data: %v", err)
	}
	return id
}

func retrieve(t testing.TB, db *db.JSONLineDB, name string, id int, v interface{}) {
	t.Helper()
	result, err := db.Retrieve(name, id)
	if err != nil {
		t.Fatalf("retrieving data: %v", err)
	}

	if err := json.Unmarshal(result, v); err != nil {
		t.Fatalf("unmarhalling json: %v", err)
	}
}

func update(t testing.TB, db *db.JSONLineDB, name string, id int, data interface{}) {
	t.Helper()
	encodedData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshalling json: %v", err)
	}
	if err := db.Update(name, id, encodedData); err != nil {
		t.Fatalf("updating data: %v", err)
	}
}
