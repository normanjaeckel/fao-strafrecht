package srv_test

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/srv"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/testutils"
)

func TestStart(t *testing.T) {
	logger := log.Default()

	es, _, cleanup := testutils.CreateEventstore(t, logger)
	defer cleanup()

	model, err := model.New(es)
	if err != nil {
		t.Fatalf("loading model: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan error, 1)

	go func() {
		ch <- srv.Start(ctx, logger, model, ":8080")
	}()
	cancel()

	srvErr := <-ch
	if srvErr != nil {
		t.Fatalf("got error from closed server: %v", srvErr)
	}

}

func TestClientHandler(t *testing.T) {
	logger := log.Default()

	es, _, cleanup := testutils.CreateEventstore(t, logger)
	defer cleanup()

	model, err := model.New(es)
	if err != nil {
		t.Fatalf("loading model: %v", err)
	}

	ts := httptest.NewServer(srv.Handler(model))
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

	t.Run("test retrieve cases", func(t *testing.T) {
		t.Fail()
	})

}
