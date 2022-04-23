package db_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/db"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/deps"
)

type SomeCaseFields struct {
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
	firstData := SomeCaseFields{
		Foo: "foo",
		Bar: "bar",
	}

	t.Run("simple insert and retrieve", func(t *testing.T) {
		id := insertCase(t, db, firstData)

		expected := 1
		if id != expected {
			t.Fatalf("wrong id, expected %d, got %d", expected, id)
		}

		result := SomeCaseFields{}
		retrieveCase(t, db, id, &result)

		if firstData != result {
			t.Fatalf("wrong retrieved data, expected %v, got %v", firstData, result)
		}
	})

	t.Run("second insert", func(t *testing.T) {
		data := SomeCaseFields{
			Foo: "foo2",
		}

		id := insertCase(t, db, data)

		expected := 2
		if id != expected {
			t.Fatalf("wrong id, expected %d, got %d", expected, id)
		}
	})

	t.Run("update", func(t *testing.T) {
		newData := SomeCaseFields{
			Foo: "foo3",
		}
		updateCase(t, db, 1, newData)

		result := SomeCaseFields{}
		retrieveCase(t, db, 1, &result)

		expected := SomeCaseFields{
			Foo: newData.Foo,
			Bar: "",
		}

		if expected != result {
			t.Fatalf("wrong retrieved data, expected %v, got %v", expected, result)
		}
	})

	t.Run("retrieve all", func(t *testing.T) {
		result, err := db.RetrieveCaseAll()
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

	data := SomeCaseFields{
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
	insertCase(t, db1, data)
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
	result := SomeCaseFields{}
	retrieveCase(t, db2, 1, &result)

	if data != result {
		t.Fatalf("wrong retrieved data, expected %v, got %v", data, result)
	}
}

func insertCase(t testing.TB, db deps.Database, data interface{}) int {
	t.Helper()
	encodedData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshalling json: %v", err)
	}
	id, err := db.InsertCase(encodedData)
	if err != nil {
		t.Fatalf("inserting data: %v", err)
	}
	return id
}

func retrieveCase(t testing.TB, db deps.Database, id int, v interface{}) {
	t.Helper()
	result, err := db.RetrieveCase(id)
	if err != nil {
		t.Fatalf("retrieving data: %v", err)
	}

	if err := json.Unmarshal(result, v); err != nil {
		t.Fatalf("unmarhalling json: %v", err)
	}
}

func updateCase(t testing.TB, db deps.Database, id int, data interface{}) {
	t.Helper()
	encodedData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshalling json: %v", err)
	}
	if err := db.UpdateCase(id, encodedData); err != nil {
		t.Fatalf("updating data: %v", err)
	}
}
