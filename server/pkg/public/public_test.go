package public_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/public"
)

func TestFiles(t *testing.T) {
	ts := httptest.NewServer(public.Files())
	defer ts.Close()

	t.Run("test index.html", func(t *testing.T) {
		body := httpGET(t, ts, "/")
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

	t.Run("test elm.js", func(t *testing.T) {
		body := httpGET(t, ts, "/assets/elm.js")

		expected, err := os.ReadFile(path.Join("files", "assets", "elm.js"))
		if err != nil {
			t.Fatalf("reading elm.js: %v", err)
		}

		if !bytes.Equal(body, expected) {
			t.Fatalf("wrong content of elm.js")
		}
	})

	t.Run("test 404", func(t *testing.T) {
		body := httpGET(t, ts, "/unknown/path")
		expected := []byte("404 page not found\n")
		if !bytes.Equal(body, expected) {
			t.Fatalf("wrong content of body in 404 case: expected %q, got %q", expected, body)
		}
	})
}

func httpGET(t testing.TB, ts *httptest.Server, urlPath string) []byte {
	t.Helper()

	res, err := http.Get(ts.URL + urlPath)
	if err != nil {
		t.Fatalf("issuing GET request to %q: %v", urlPath, err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}

	return body
}
