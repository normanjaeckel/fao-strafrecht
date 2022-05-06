/*
Package lawcase is about criminal law cases.
*/
package lawcase

import (
	"encoding/json"
	"fmt"
	"io"
)

type Model map[int]Case

type Case struct {
	Rubrum string `json:"Rubrum"`
	Az     string `json:"Az"`
}

type decodedMsg struct {
	ID     int  `json:"ID"`
	Fields Case `json:"Fields"`
}

func (cs *Model) Load(msg json.RawMessage) error {
	if msg == nil {
		return fmt.Errorf("message must not be nil")
	}
	var d decodedMsg
	if err := json.Unmarshal(msg, &d); err != nil {
		return fmt.Errorf("unmarshalling JSON: %v", err)
	}
	if d.ID < 1 {
		return fmt.Errorf("message contains invalid id %d", d.ID)
	}
	(*cs)[d.ID] = d.Fields
	return nil
}

func (cs *Model) AddCase(c Case, w io.Writer) (int, error) {
	// TODO: Validate case: https://pkg.go.dev/github.com/go-playground/validator

	newID := cs.maxCaseID() + 1
	d := decodedMsg{
		ID:     newID,
		Fields: c,
	}
	b, err := json.Marshal(d)
	if err != nil {
		return 0, fmt.Errorf("marshalling JSON event data: %w", err)
	}
	if _, err := w.Write(b); err != nil {
		return 0, fmt.Errorf("writing event data: %w", err)
	}
	(*cs)[newID] = c
	return newID, nil
}

func (cs Model) maxCaseID() int {
	var result int
	for result = range cs {
		break
	}
	for n := range cs {
		if n > result {
			result = n
		}
	}
	return result
}

func (cs Model) Retrieve(id int) (Case, error) {
	c, ok := cs[id]
	if !ok {
		return Case{}, fmt.Errorf("case %d does not exist", id)
	}
	return c, nil
}
