/*
Package lawcase is about criminal law cases.
*/
package lawcase

import (
	"encoding/json"
	"fmt"
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

func (cs *Model) AddCase(c Case) (json.RawMessage, int) {
	// TODO: Maybe do some validation and the also return an error, then remove the panic line.
	newID := cs.maxCaseID() + 1
	d := decodedMsg{
		ID:     newID,
		Fields: c,
	}
	b, err := json.Marshal(d)
	if err != nil {
		panic(fmt.Sprintf("marshalling JSON event data: %v; this should never ever happen", err))
	}
	(*cs)[newID] = c
	return b, newID
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
