package model_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model/lawcase"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/testutils"
)

func TestModel(t *testing.T) {
	logger := log.Default()
	es, _, cleanup := testutils.CreateEventstore(t, logger)
	defer cleanup()

	randStr := "Eash3bae6b"
	id := 1
	msg := fmt.Sprintf(`{"Name": "Case", "Data": {"ID": %d, "Fields": {"rubrum": "%s"}}}`, id, randStr)
	firstCase := json.RawMessage([]byte(msg))
	if _, err := es.Write(firstCase); err != nil {
		t.Fatalf("writing first case: %v", err)
	}

	t.Run("test insert and retrieve", func(t *testing.T) {

		m, err := model.New(es)
		if err != nil {
			t.Fatalf("creating model: %v", err)
		}
		c, err := m.Case.Retrieve(id)
		if err != nil {
			t.Fatalf("retrieving case: %v", err)
		}
		got := c.Rubrum
		if got != randStr {
			t.Fatalf("wrong model content for one case: expected rubrum %q, got %q", randStr, got)
		}
	})

	t.Run("test second insert via AddCase and retrieve", func(t *testing.T) {
		m, err := model.New(es)
		if err != nil {
			t.Fatalf("creating model: %v", err)
		}

		secondCase := lawcase.Case{
			Rubrum: randStr + randStr,
		}

		if _, err := m.Case.AddCase(secondCase, m.WriteEvent("Case")); err != nil {
			t.Fatalf("adding case: %v", err)
		}

		c, err := m.Case.Retrieve(2)
		if err != nil {
			t.Fatalf("retrieving case: %v", err)
		}
		got := c.Rubrum
		if got != randStr+randStr {
			t.Fatalf("wrong model content for one case: expected rubrum %q, got %q", randStr+randStr, got)
		}
	})
}
