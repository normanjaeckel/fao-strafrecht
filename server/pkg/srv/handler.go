package srv

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model"
	"github.com/normanjaeckel/fao-strafrecht/server/pkg/model/lawcase"
)

type CaseHandler struct {
	Logger Logger
	Model  *model.Model
}

func NewCaseHandler(logger Logger, m *model.Model) *CaseHandler {
	return &CaseHandler{
		Logger: logger,
		Model:  m,
	}
}

func (h CaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/retrieve", h.RetrieveCases())
	mux.HandleFunc("/new", h.NewCase())
	mux.ServeHTTP(w, r)
}

func (h CaseHandler) RetrieveCases() func(http.ResponseWriter, *http.Request) {
	return methodAllowed(
		http.MethodGet,
		func(w http.ResponseWriter, r *http.Request) {
			b, err := json.Marshal(h.Model.Case)
			if err != nil {
				msg := fmt.Sprintf("Error: marshalling JSON: %v", err)
				h.Logger.Printf(msg)
				http.Error(w, msg, 500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(b); err != nil {
				msg := fmt.Sprintf("Error: writing response body: %v", err)
				h.Logger.Printf(msg)
				http.Error(w, msg, 500)
				return
			}
		},
	)
}

func (h CaseHandler) NewCase() func(http.ResponseWriter, *http.Request) {
	return methodAllowed(
		http.MethodPost,
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "Error: Content-Type must be application/json", http.StatusBadRequest)
				return
			}

			d := json.NewDecoder(r.Body)
			c := lawcase.Case{}
			if err := d.Decode(&c); err != nil {
				http.Error(w, fmt.Sprintf("Error: decoding request: %v", err), http.StatusBadRequest)
				return
			}

			v := validator.New()
			if err := v.Struct(c); err != nil {
				http.Error(w, fmt.Sprintf("Error: invalid request:\n%v", err), http.StatusBadRequest)
				return
			}

			id, err := h.Model.Case.AddCase(c, h.Model.WriteEvent("Case"))
			if err != nil {
				msg := fmt.Sprintf("Error: adding case: %v", err)
				h.Logger.Printf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			respBody := []byte(fmt.Sprintf(`{"id":%d}`, id))
			if _, err := w.Write(respBody); err != nil {
				msg := fmt.Sprintf("Error: writing response body: %v", err)
				h.Logger.Printf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		},
	)
}

func methodAllowed(method string, fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fn(w, r)
	}
	return wrapper
}
