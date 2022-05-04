/*
Package model gives access to all model objects.
*/
package model

import (
	"encoding/json"
	"fmt"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model/lawcase"
)

type Eventstore interface {
	Save(json.RawMessage) error
	Retrieve() ([]json.RawMessage, error)
}

type Model struct {
	eventstore Eventstore
	Case       lawcase.Model
}

type decodedEvent struct {
	Name string          `json:"Name"`
	Data json.RawMessage `json:"Data"`
}

func New(es Eventstore) (*Model, error) {
	m := Model{
		eventstore: es,
		Case:       lawcase.Model{},
	}

	events, err := es.Retrieve()
	if err != nil {
		return nil, fmt.Errorf("retrieving events from eventstore: %w", err)
	}

	for i, e := range events {
		var d decodedEvent
		if err := json.Unmarshal(e, &d); err != nil {
			return nil, fmt.Errorf("unmarshalling JSON (line %d): %w", i+1, err)
		}

		switch d.Name {
		case "Case":
			if err := m.Case.Load(d.Data); err != nil {
				return nil, fmt.Errorf("loading case: %w", err)
			}
		case "Theme":
			return nil, fmt.Errorf("not implemented")
		default:
			return nil, fmt.Errorf("invalid event %q", string(e))
		}
	}

	return &m, nil
}

func (m *Model) WriteEvent(name string, data json.RawMessage) error {
	d := decodedEvent{
		Name: name,
		Data: data,
	}
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("marshalling JSON event: %w", err)
	}
	if err := m.eventstore.Save(b); err != nil {
		return fmt.Errorf("saving new event to eventstore: %w", err)
	}
	return nil
}

// Theme

// func (m Model) RetrieveTheme() string {
// 	return m.Theme
// }

// func (m *Model) SetTheme(t string) error {
// 	return fmt.Errorf("not implemented")
// }
