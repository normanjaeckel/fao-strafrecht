package eventstore_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/eventstore"
)

func TestDB(t *testing.T) {
	dir, err := os.MkdirTemp("", "fao-strafrecht-")
	if err != nil {
		t.Fatalf("creating tmp directors: %v", err)
	}

	es, close, err := eventstore.New(path.Join(dir, "ds.jsonl"))
	if err != nil {
		t.Fatalf("loading eventstore: %v", err)
	}
	defer close()

	testData0 := json.RawMessage(`{"foo":"bar 0"}`)
	testData1 := json.RawMessage(`{"foo":"bar 1"}`)

	t.Run("test invalid event", func(t *testing.T) {
		invalidEvent := json.RawMessage("invalid")
		err := es.Save(invalidEvent)
		if err == nil {
			t.Fatalf("expected error but got nil")
		}
		errMsg := `invalid JSON encoding for event "invalid"`
		if err.Error() != errMsg {
			t.Fatalf("wrong error message: expected %q, got %q", errMsg, err.Error())
		}
	})

	t.Run("test first save", func(t *testing.T) {
		if err := es.Save(testData0); err != nil {
			t.Fatalf("saving test data %v", err)
		}
	})

	t.Run("test second save", func(t *testing.T) {
		if err := es.Save(testData1); err != nil {
			t.Fatalf("saving test data %v", err)
		}
	})

	t.Run("test retrieve", func(t *testing.T) {
		data, err := es.Retrieve(0)
		if err != nil {
			t.Fatalf("retrieving data: %v", err)
		}
		if len(data) != 2 {
			t.Fatalf("length of retrieved data must be 2 but is %d", len(data))
		}

		if !bytes.Equal(data[0], testData0) {
			t.Fatalf("wrong content: expected %q, got %q", testData0, data[0])
		}
		if !bytes.Equal(data[1], testData1) {
			t.Fatalf("wrong content: expected %q, got %q", testData1, data[1])
		}
	})

	// TODO: Retrieve from position.

}

// 	firstData := SomeCaseFields{
// 		Foo: "foo",
// 		Bar: "bar",
// 	}

// 	t.Run("simple insert and retrieve", func(t *testing.T) {
// 		id := insertCase(t, db, firstData)

// 		expected := 1
// 		if id != expected {
// 			t.Fatalf("wrong id, expected %d, got %d", expected, id)
// 		}

// 		result := SomeCaseFields{}
// 		retrieveCase(t, db, id, &result)

// 		if firstData != result {
// 			t.Fatalf("wrong retrieved data, expected %v, got %v", firstData, result)
// 		}
// 	})

// 	t.Run("second insert", func(t *testing.T) {
// 		data := SomeCaseFields{
// 			Foo: "foo2",
// 		}

// 		id := insertCase(t, db, data)

// 		expected := 2
// 		if id != expected {
// 			t.Fatalf("wrong id, expected %d, got %d", expected, id)
// 		}
// 	})

// 	t.Run("update", func(t *testing.T) {
// 		newData := SomeCaseFields{
// 			Foo: "foo3",
// 		}
// 		updateCase(t, db, 1, newData)

// 		result := SomeCaseFields{}
// 		retrieveCase(t, db, 1, &result)

// 		expected := SomeCaseFields{
// 			Foo: newData.Foo,
// 			Bar: "",
// 		}

// 		if expected != result {
// 			t.Fatalf("wrong retrieved data, expected %v, got %v", expected, result)
// 		}
// 	})

// 	t.Run("retrieve all", func(t *testing.T) {
// 		result, err := db.RetrieveCaseAll()
// 		if err != nil {
// 			t.Fatalf("retrieving all data: %v", err)
// 		}
// 		expectedLen := 2
// 		got := len(result)
// 		if expectedLen != got {
// 			t.Fatalf("wrong length of result, expected %d, got %d", expectedLen, got)
// 		}
// 	})
// }

// func TestDBPersistence(t *testing.T) {
// 	dir, err := os.MkdirTemp("", "fao-strafrecht-*")
// 	if err != nil {
// 		t.Fatalf("creating tmp directors: %v", err)
// 	}
// 	defer os.RemoveAll(dir)

// 	data := SomeCaseFields{
// 		Foo: "foo1\nfoo2",
// 		Bar: "bar",
// 	}
// 	p := path.Join(dir, "db.jsonl")

// 	f1, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)
// 	if err != nil {
// 		t.Fatalf("opening database file: %v", err)
// 	}
// 	defer f1.Close()
// 	db1, err := ds.New(f1)
// 	if err != nil {
// 		t.Fatalf("loading db: %v", err)
// 	}
// 	insertCase(t, db1, data)
// 	f1.Close()

// 	f2, err := os.OpenFile(p, os.O_APPEND|os.O_RDWR, 0644)
// 	if err != nil {
// 		t.Fatalf("opening database file: %v", err)
// 	}
// 	defer f2.Close()
// 	db2, err := db.New(f2)
// 	if err != nil {
// 		t.Fatalf("loading db: %v", err)
// 	}
// 	result := SomeCaseFields{}
// 	retrieveCase(t, db2, 1, &result)

// 	if data != result {
// 		t.Fatalf("wrong retrieved data, expected %v, got %v", data, result)
// 	}
// }

// type DatabaseOld interface {
// 	InsertCase(fields json.RawMessage) (int, error)
// 	UpdateCase(id int, fields json.RawMessage) error
// 	RetrieveCase(id int) (json.RawMessage, error)
// 	RetrieveCaseAll() (map[int]json.RawMessage, error)
// }

// func insertCase(t testing.TB, db DatabaseOld, data interface{}) int {
// 	t.Helper()
// 	encodedData, err := json.Marshal(data)
// 	if err != nil {
// 		t.Fatalf("marshalling json: %v", err)
// 	}
// 	id, err := db.InsertCase(encodedData)
// 	if err != nil {
// 		t.Fatalf("inserting data: %v", err)
// 	}
// 	return id
// }

// func retrieveCase(t testing.TB, db DatabaseOld, id int, v interface{}) {
// 	t.Helper()
// 	result, err := db.RetrieveCase(id)
// 	if err != nil {
// 		t.Fatalf("retrieving data: %v", err)
// 	}

// 	if err := json.Unmarshal(result, v); err != nil {
// 		t.Fatalf("unmarhalling json: %v", err)
// 	}
// }

// func updateCase(t testing.TB, db DatabaseOld, id int, data interface{}) {
// 	t.Helper()
// 	encodedData, err := json.Marshal(data)
// 	if err != nil {
// 		t.Fatalf("marshalling json: %v", err)
// 	}
// 	if err := db.UpdateCase(id, encodedData); err != nil {
// 		t.Fatalf("updating data: %v", err)
// 	}
// }
