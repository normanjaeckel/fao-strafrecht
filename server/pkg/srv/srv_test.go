package srv_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/srv"
)

type FakeDB struct{}

func (db *FakeDB) Insert(name string, data json.RawMessage) (int, error) {
	return 0, fmt.Errorf("Not implemented")
}

func (db *FakeDB) Retrieve(name string, id int) (json.RawMessage, error) {
	return json.RawMessage{}, fmt.Errorf("Not implemented")
}

func (db *FakeDB) Update(name string, id int, data json.RawMessage) error {
	return fmt.Errorf("Not implemented")
}

func (db *FakeDB) RetrieveAll(name string) (map[int]json.RawMessage, error) {
	return nil, fmt.Errorf("Not implemented")
}

func TestStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan error, 1)

	logger := log.Default()
	db := FakeDB{}

	go func() {
		ch <- srv.Start(ctx, logger, &db, ":8080")
	}()
	cancel()

	srvErr := <-ch
	if srvErr != nil {
		t.Fatalf("got error from closed server: %v", srvErr)
	}

}

func TestClientHandler(t *testing.T) {
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	t.Run("test root path", func(t *testing.T) {
		path := "/"

		res, err := http.Get(ts.URL + path)
		if err != nil {
			t.Fatalf("issuing GET request to %q: %v", path, err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("reading response body: %v", err)
		}

		expected := `<!doctype html>

<html lang="de">

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="description" content="">
  <title>Fachanwalt für Strafrecht | Fallliste</title>
`
		if !strings.HasPrefix(string(body), expected) {
			t.Fatalf("wrong beginning of response body: expected %q, got (full data) %q", expected, string(body))
		}
	})

}
