package eventstore_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/testutils"
)

func TestEventstore(t *testing.T) {
	es, filename, cleanup := testutils.CreateEventstore(t)
	defer cleanup()

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

		f, err := os.Open(filename)
		if err != nil {
			t.Fatalf("opening datastore file %q: %v", filename, err)
		}
		b, err := io.ReadAll(f)
		if err != nil {
			t.Fatalf("reading datastore file %q: %v", filename, err)
		}
		var line struct {
			Timestamp int64
		}
		json.Unmarshal(bytes.Split(b, []byte("\n"))[1], &line)

		expectedTimestamp := time.Now().Unix()
		if line.Timestamp != expectedTimestamp {
			t.Fatalf("wrong timestamp: expected %d, got %d", expectedTimestamp, line.Timestamp)
		}
	})

	t.Run("test retrieve", func(t *testing.T) {
		data, err := es.Retrieve()
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
}
