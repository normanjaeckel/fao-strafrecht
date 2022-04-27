/*
Package model provides handlers to deal with the model.
*/
package model

import (
	"encoding/json"
	"net/http"
)

type Eventstore interface {
	//Save(json.RawMessage) error
	Retrieve() ([]json.RawMessage, error)
}

type CaseHandler struct {
	Eventstore Eventstore
}

func NewCaseHandler(es Eventstore) *CaseHandler {
	return &CaseHandler{
		Eventstore: es,
	}
}

func (h CaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/retrieve", RetrieveCases(h.Eventstore))
	mux.ServeHTTP(w, r)
}

func RetrieveCases(es Eventstore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Go ahead here.
	}
}
