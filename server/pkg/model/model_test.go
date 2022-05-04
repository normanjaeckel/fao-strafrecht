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
	es.Save(firstCase)

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

		if err := m.WriteEvent("Case", m.Case.AddCase(secondCase)); err != nil {
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
