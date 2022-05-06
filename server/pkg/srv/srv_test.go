package srv_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

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
	ts, _, cleanup := testutils.CreateServer(t, logger)
	defer cleanup()

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

func TestRetrieveCaseHandler(t *testing.T) {
	logger := log.Default()
	ts, _, cleanup := testutils.CreateServer(t, logger)
	defer cleanup()

	t.Run("test retrieve cases", func(t *testing.T) {
		path := "/api/case/retrieve"

		res, err := http.Get(ts.URL + path)
		if err != nil {
			t.Fatalf("issuing GET request to %q: %v", path, err)
		}
		respBody, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("reading response body: %v", err)
		}

		expected := "{}"
		if string(respBody) != expected {
			t.Fatalf("wrong response body: expected %q, got %q", expected, string(respBody))
		}

		expectedCTHeader := "application/json"
		gotCTHeader := res.Header.Get("Content-Type")
		if expectedCTHeader != gotCTHeader {
			t.Fatalf("wrong response Content-Type header: expected %q, got %q", expectedCTHeader, gotCTHeader)
		}
	})

	// TODO: Add test to retrieve one case
}

func TestNewCaseHandler(t *testing.T) {
	logger := log.Default()
	ts, filename, cleanup := testutils.CreateServer(t, logger)
	defer cleanup()

	path := "/api/case/new"

	t.Run("one POST request", func(t *testing.T) {
		reqBody := []byte(`{"Rubrum": "test_rubrum_beiTh9itha", "Az": "test_az_uwwe34sdf1"}`)

		res, err := http.Post(ts.URL+path, "application/json", bytes.NewReader(reqBody))
		if err != nil {
			t.Fatalf("issuing POST request to %q: %v", path, err)
		}

		if res.StatusCode != 200 {
			t.Fatalf("wrong status code: expected 200, got %d", res.StatusCode)
		}

		respBody, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("reading response body: %v", err)
		}

		expected := `{"id":1}`
		if string(respBody) != expected {
			t.Fatalf("wrong response body: expected %q, got %q", expected, string(respBody))
		}

		expectedCTHeader := "application/json"
		gotCTHeader := res.Header.Get("Content-Type")
		if expectedCTHeader != gotCTHeader {
			t.Fatalf("wrong response Content-Type header: expected %q, got %q", expectedCTHeader, gotCTHeader)
		}

		gotEventstore, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatalf("reading eventstore file: %v", err)
		}
		expectedEventstore := []byte(fmt.Sprintf(
			`{"Event":{"Name":"Case","Data":{"ID":1,"Fields":{"Rubrum":"test_rubrum_beiTh9itha","Az":"test_az_uwwe34sdf1"}}},"Timestamp":%d}`, time.Now().Unix(),
		))
		expectedEventstore = append(expectedEventstore, '\n')
		if !bytes.Equal(expectedEventstore, gotEventstore) {
			t.Fatalf("wrong content of eventstore: expected %q, got %q", expectedEventstore, gotEventstore)
		}
	})

	t.Run("invalid request, invalid JSON", func(t *testing.T) {
		reqBody := []byte(`invalid`)

		res, err := http.Post(ts.URL+path, "application/json", bytes.NewReader(reqBody))
		if err != nil {
			t.Fatalf("issuing POST request to %q: %v", path, err)
		}

		if res.StatusCode != 400 {
			t.Fatalf("wrong status code: expected 400, got %d", res.StatusCode)
		}

		respBody, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("reading response body: %v", err)
		}

		expected := "Error: decoding request: invalid character 'i' looking for beginning of value\n"
		if string(respBody) != expected {
			t.Fatalf("wrong response body: expected %q, got %q", expected, string(respBody))
		}

		expectedCTHeader := "text/plain; charset=utf-8"
		gotCTHeader := res.Header.Get("Content-Type")
		if expectedCTHeader != gotCTHeader {
			t.Fatalf("wrong response Content-Type header: expected %q, got %q", expectedCTHeader, gotCTHeader)
		}
	})

	t.Run("invalid request, bad values", func(t *testing.T) {
		reqBody := []byte(`{"unknown":"key"}`)

		res, err := http.Post(ts.URL+path, "application/json", bytes.NewReader(reqBody))
		if err != nil {
			t.Fatalf("issuing POST request to %q: %v", path, err)
		}

		if res.StatusCode != 400 {
			t.Fatalf("wrong status code: expected 400, got %d", res.StatusCode)
		}

		respBody, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("reading response body: %v", err)
		}

		expected := "Error: ...\n"
		if string(respBody) != expected {
			t.Fatalf("wrong response body: expected %q, got %q", expected, string(respBody))
		}

		expectedCTHeader := "text/plain; charset=utf-8"
		gotCTHeader := res.Header.Get("Content-Type")
		if expectedCTHeader != gotCTHeader {
			t.Fatalf("wrong response Content-Type header: expected %q, got %q", expectedCTHeader, gotCTHeader)
		}
	})
}
