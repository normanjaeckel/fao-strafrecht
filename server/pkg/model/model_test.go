package model_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/testutils"
)

func TestModel(t *testing.T) {
	es, _, cleanup := testutils.CreateEventstore(t)
	defer cleanup()

	randStr := "Eash3bae6b"
	id := 1
	msg := fmt.Sprintf(`{"Name": "Case", "ID": %d, "Fields": {"rubrum": "%s"}}`, id, randStr)
	firstCase := json.RawMessage([]byte(msg))
	es.Save(firstCase)

	t.Run("test insert and retrieve", func(t *testing.T) {

		m, err := model.New(es)
		if err != nil {
			t.Fatalf("creating model: %v", err)
		}
		got := m.RetrieveCase()[id].Rubrum
		if got != randStr {
			t.Fatalf("wrong model content for one case: expected rubrum %q, got %q", randStr, got)
		}
	})

	t.Run("test second insert via AddCase and retriev", func(t *testing.T) {
		m, err := model.New(es)
		if err != nil {
			t.Fatalf("creating model: %v", err)
		}

		secondCase := model.Case{
			Rubrum: randStr + randStr,
		}
		if err := m.AddCase(secondCase); err != nil {
			t.Fatalf("adding case: %v", err)
		}

		got := m.RetrieveCase()[2].Rubrum
		if got != randStr+randStr {
			t.Fatalf("wrong model content for one case: expected rubrum %q, got %q", randStr+randStr, got)
		}
	})
}
