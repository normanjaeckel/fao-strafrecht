package srv

import (
	"encoding/json"
	"net/http"

	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model"
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
	mux.ServeHTTP(w, r)
}

func RetrieveCases(m *model.Model) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, err := json.Marshal(m.Case)
		if err != nil {
			w.WriteHeader(500)
		}
		if _, err := w.Write(b); err != nil {
			w.WriteHeader(500)
		}
	}
}
