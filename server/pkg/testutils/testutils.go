/*
Package testutils contains some helpers for tests.
*/
package testutils

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/eventstore"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/srv"
)

type Eventstore interface {
	io.Writer
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
		defer close()
		defer os.RemoveAll(dir)
	}

	return es, filename, cleanupFn
}

func CreateServer(t testing.TB, logger eventstore.Logger) (*httptest.Server, string, func()) {

	es, filename, esCleanup := CreateEventstore(t, logger)

	model, err := model.New(es)
	if err != nil {
		t.Fatalf("loading model: %v", err)
	}

	ts := httptest.NewServer(srv.Handler(logger, model))

	cleanupFn := func() {
		defer ts.Close()
		defer esCleanup()
	}

	return ts, filename, cleanupFn
}
