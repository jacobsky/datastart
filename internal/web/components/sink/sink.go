package sink

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

func AddRoutes(mux *http.ServeMux) {

	sink := NewHandler()
	mux.Handle("/sink", sink)
	mux.HandleFunc("/sink/testsse", testSee)
}

type Handler struct{}

func NewHandler() http.Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templ.Handler(Sink()).ServeHTTP(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func testSee(w http.ResponseWriter, r *http.Request) {
	see := datastar.NewSSE(w, r)
	event_format := `{'data': 'this is event #%v'}`
	for i := range 10 {

		event := fmt.Sprintf(event_format, i+1)
		err := see.PatchSignals([]byte(event))

		if err != nil {
			slog.Error("Error has occurred", "error", err)
		}
		time.Sleep(1 * time.Second)
	}

	err := see.PatchSignals([]byte(`{'data': 'no more events resetting'}`))

	if err != nil {
		slog.Error("Error has occurred", "error", err)
	}
	time.Sleep(1 * time.Second)
	err = see.PatchSignals([]byte(`{'data': '...'}`))

	if err != nil {
		slog.Error("Error has occurred", "error", err)
	}
}
