/*
Package testutils contains some helpers for tests.
*/
package testutils

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/eventstore"
)

type Eventstore interface {
	Save(json.RawMessage) error
	Retrieve() ([]json.RawMessage, error)
}

func CreateEventstore(t testing.TB, logger eventstore.Logger) (Eventstore, string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "fao-strafrecht-")
	if err != nil {
		t.Fatalf("creating tmp directory: %v", err)
	}

	filename := path.Join(dir, "ds.jsonl")

	es, close, err := eventstore.New(logger, filename)
	if err != nil {
		t.Fatalf("loading eventstore at %q: %v", filename, err)
	}

	cleanupFn := func() {
		close()
		os.RemoveAll(dir)
	}

	return es, filename, cleanupFn
}
