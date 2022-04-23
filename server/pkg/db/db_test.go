package db_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/db"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/deps"
)

type SomeData struct {
	Foo string
	Bar string
}

func TestDB(t *testing.T) {
	f, err := os.CreateTemp("", "db-*.jsonl")
	if err != nil {
		t.Fatalf("creating tmp file: %v", err)
	}
	defer f.Close()

	db, err := db.New(f)
	if err != nil {
		t.Fatalf("loading db: %v", err)
	}
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

func TestDBPersistence(t *testing.T) {
	dir, err := os.MkdirTemp("", "fao-strafrecht-*")
	if err != nil {
		t.Fatalf("creating tmp directors: %v", err)
	}
	defer os.RemoveAll(dir)

	name := "my collection"
	data := SomeData{
		Foo: "foo1\nfoo2",
		Bar: "bar",
	}
	p := path.Join(dir, "db.jsonl")

	f1, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("opening database file: %v", err)
	}
	defer f1.Close()
	db1, err := db.New(f1)
	if err != nil {
		t.Fatalf("loading db: %v", err)
	}
	insert(t, db1, name, data)
	f1.Close()

	f2, err := os.OpenFile(p, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("opening database file: %v", err)
	}
	defer f2.Close()
	db2, err := db.New(f2)
	if err != nil {
		t.Fatalf("loading db: %v", err)
	}
	result := SomeData{}
	retrieve(t, db2, name, 1, &result)

	if data != result {
		t.Fatalf("wrong retrieved data, expected %v, got %v", data, result)
	}
}

func insert(t testing.TB, db deps.Database, name string, data interface{}) int {
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

func retrieve(t testing.TB, db deps.Database, name string, id int, v interface{}) {
	t.Helper()
	result, err := db.Retrieve(name, id)
	if err != nil {
		t.Fatalf("retrieving data: %v", err)
	}

	if err := json.Unmarshal(result, v); err != nil {
		t.Fatalf("unmarhalling json: %v", err)
	}
}

func update(t testing.TB, db deps.Database, name string, id int, data interface{}) {
	t.Helper()
	encodedData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshalling json: %v", err)
	}
	if err := db.Update(name, id, encodedData); err != nil {
		t.Fatalf("updating data: %v", err)
	}
}
