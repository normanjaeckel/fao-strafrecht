package lawcase_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model/lawcase"
)

func TestLoad(t *testing.T) {
	t.Run("load one message", func(t *testing.T) {
		m := lawcase.Model{}
		expectedRubrum := "VisheRai6h 1"
		expectedID := 1
		msgJSON := fmt.Sprintf(`{"ID": %d, "Fields": {"Rubrum": "%s"}}`, expectedID, expectedRubrum)
		msg := json.RawMessage([]byte(msgJSON))

		if err := m.Load(msg); err != nil {
			t.Fatalf("loading message: %v", err)
		}

		got := m[1].Rubrum
		if expectedRubrum != got {
			t.Fatalf("wrong rubrum; expected %q, got %q", expectedRubrum, got)
		}
	})

	t.Run("load wrong message", func(t *testing.T) {
		m := lawcase.Model{}
		expectedRubrum := "VisheRai6h 2"
		msgJSON := fmt.Sprintf(`{"Fields": {"Rubrum": "%s"}}`, expectedRubrum)
		msg := json.RawMessage([]byte(msgJSON))

		err := m.Load(msg)

		expectedErrMsg := "message contains invalid id 0"
		if err == nil || err.Error() != expectedErrMsg {
			t.Fatalf("expected error %q, got %v", expectedErrMsg, err)
		}
		if len(m) != 0 {
			t.Fatalf("wrong lenght of model, expected 0, got %d", len(m))
		}
	})

	t.Run("load empty message", func(t *testing.T) {
		m := lawcase.Model{}

		err := m.Load(nil)

		expectedErrMsg := "message must not be nil"
		if err == nil || err.Error() != expectedErrMsg {
			t.Fatalf("expected error %q, got %v", expectedErrMsg, err)
		}
		if len(m) != 0 {
			t.Fatalf("wrong lenght of model, expected 0, got %d", len(m))
		}
	})
}

func TestAddCase(t *testing.T) {
	m := lawcase.Model{}

	t.Run("add one case", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		expectedRubrum := "rubrum yie5Athoh2"
		c := lawcase.Case{
			Rubrum: expectedRubrum,
		}

		id, err := m.AddCase(c, buf)

		if err != nil {
			t.Fatalf("adding case: %v", err)
		}
		if id != 1 {
			t.Fatalf("wrong id: expected 1, got %d", id)
		}
		content, err := io.ReadAll(buf)
		if err != nil {
			t.Fatalf("reading content from buffer given to AddCase: %v", err)
		}
		expectedMsg := []byte(fmt.Sprintf(
			`{"ID":1,"Fields":{"Rubrum":"%s","Az":"","Gericht":"","Beginn":"","Ende":"","Gegenstand":"","Art":"","Beschreibung":"","Stand":""}}`,
			expectedRubrum,
		))
		if !bytes.Equal(content, expectedMsg) {
			t.Fatalf("wrong message, expected %q, got %q", expectedMsg, content)
		}
		res, err := m.Retrieve(1)
		if err != nil {
			t.Fatalf("retrieving case: %v", err)
		}
		if res.Rubrum != expectedRubrum {
			t.Fatalf("wrong content of model")
		}
	})

	t.Run("add second case", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		id, err := m.AddCase(lawcase.Case{}, buf)

		if err != nil {
			t.Fatalf("adding case: %v", err)
		}
		if id != 2 {
			t.Fatalf("wrong id: expected 2, got %d", id)
		}
		content, err := io.ReadAll(buf)
		if err != nil {
			t.Fatalf("reading content from buffer given to AddCase: %v", err)
		}
		expectedMsg := []byte(`{"ID":2,"Fields":{"Rubrum":"","Az":"","Gericht":"","Beginn":"","Ende":"","Gegenstand":"","Art":"","Beschreibung":"","Stand":""}}`)
		if !bytes.Equal(content, expectedMsg) {
			t.Fatalf("wrong message, expected %q, got %q", expectedMsg, content)
		}
		if _, err := m.Retrieve(2); err != nil {
			t.Fatalf("retrieving case: %v", err)
		}
	})

	t.Run("retrieve not existing case", func(t *testing.T) {
		id := 42
		_, err := m.Retrieve(id)
		expectedErrMsg := fmt.Sprintf("case %d does not exist", id)
		if err == nil || err.Error() != expectedErrMsg {
			t.Fatalf("expected error %q, got %v", expectedErrMsg, err)
		}

	})
}
