package srv

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model/lawcase"
)

type CaseHandler struct {
	Model *model.Model
}

func NewCaseHandler(m *model.Model) *CaseHandler {
	return &CaseHandler{Model: m}
}

func (h CaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/retrieve", RetrieveCases(h.Model))
	mux.HandleFunc("/new", NewCase(h.Model))
	mux.ServeHTTP(w, r)
}

func RetrieveCases(m *model.Model) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, err := json.Marshal(m.Case)
		if err != nil {
			w.WriteHeader(500) // TODO: Log error
		}
		if _, err := w.Write(b); err != nil {
			w.WriteHeader(500) // TODO: Log error
		}
	}
}

func NewCase(m *model.Model) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d := json.NewDecoder(r.Body)
		c := lawcase.Case{}
		if err := d.Decode(&c); err != nil {
			w.WriteHeader(500) // TODO: Log error
		}
		msg, id := m.Case.AddCase(c)
		if err := m.WriteEvent("Case", msg); err != nil {
			w.WriteHeader(500) // TODO: Log error
		}

		w.Header().Set("Content-Type", "application/json")
		respBody := []byte(fmt.Sprintf(`{"id":%d}`, id))
		if _, err := w.Write(respBody); err != nil {
			w.WriteHeader(500) // TODO: Log error
		}
	}
}

// https://pkg.go.dev/github.com/go-playground/validator
