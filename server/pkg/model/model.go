/*
Package model gives access to all model objects.
*/
package model

import (
	"encoding/json"
	"fmt"
)

type Eventstore interface {
	Save(json.RawMessage) error
	Retrieve() ([]json.RawMessage, error)
}

type Model struct {
	Eventstore Eventstore
	Data       ModelData
}

type ModelData struct {
	Case  map[int]Case
	Theme string
}

func initData() ModelData {
	return ModelData{
		Case: map[int]Case{},
	}
}

func New(es Eventstore) (*Model, error) {
	m := Model{
		Eventstore: es,
		Data:       initData(),
	}

	events, err := es.Retrieve()
	if err != nil {
		return nil, fmt.Errorf("retrieving events from eventstore: %w", err)
	}

	for i, e := range events {
		var decodedEvent struct {
			Name string `json:"Name"`
		}
		if err := json.Unmarshal(e, &decodedEvent); err != nil {
			return nil, fmt.Errorf("unmarshalling JSON (line %d): %w", i+1, err)
		}

		switch decodedEvent.Name {
		case "Case":
			id, c := jsonToCase(e)
			m.Data.Case[id] = c
		case "Theme":
			return nil, fmt.Errorf("not implemented")
		default:
			return nil, fmt.Errorf("invalid event %q", string(e))
		}
	}

	return &m, nil
}

// Case

type Case struct {
	Rubrum string `json:"Rubrum"`
	Az     string `json:"Az"`
}

func jsonToCase(msg json.RawMessage) (int, Case) {
	var decodedMsg struct {
		ID     int  `json:"ID"`
		Fields Case `json:"Fields"`
	}
	if err := json.Unmarshal(msg, &decodedMsg); err != nil {
		panic(fmt.Sprintf("unmarshalling JSON: %v; this should never ever happen", err))
	}
	return decodedMsg.ID, decodedMsg.Fields
}

func (m Model) RetrieveCase() map[int]Case {
	return m.Data.Case
}

func (m *Model) AddCase(c Case) error {
	newID := m.maxCaseID() + 1

	event := struct {
		Name   string
		ID     int
		Fields Case
	}{
		Name:   "Case",
		ID:     newID,
		Fields: c,
	}
	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshalling JSON event: %w", err)
	}
	if err := m.Eventstore.Save(b); err != nil {
		return fmt.Errorf("saving new case to eventstore: %w", err)
	}

	m.Data.Case[newID] = c
	return nil
}

func (m Model) maxCaseID() int {
	var result int
	for result = range m.Data.Case {
		break
	}
	for n := range m.Data.Case {
		if n > result {
			result = n
		}
	}
	return result
}

// Theme

func (m Model) RetrieveTheme() string {
	return m.Data.Theme
}

func (m *Model) SetTheme(t string) error {
	return fmt.Errorf("not implemented")
}

// type CaseHandler struct {
// }

// func NewCaseHandler() *CaseHandler {
// 	return &CaseHandler{}
// }

// func (h CaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/retrieve", RetrieveCases())
// 	mux.ServeHTTP(w, r)
// }

// func RetrieveCases() func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// TODO: Go ahead here.
// 	}
// }
