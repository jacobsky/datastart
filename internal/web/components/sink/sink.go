package sink

import (
	"encoding/json"
	"fmt"
	"log"
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
	mux.HandleFunc("/sink/complexsee", complexSee)

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

type complexSignal struct {
	Input    string `json:"input"`
	Repeated string `json:"repeated"`
	Data     string `json:"data"`
	Complex  struct {
		IsPressed struct {
			Client bool `json:"client"`
			Server bool `json:"server"`
		} `json:"isPressed"`
	} `json:"complex"`
}

func complexSee(w http.ResponseWriter, r *http.Request) {
	see := datastar.NewSSE(w, r)
	signals := &complexSignal{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		log.Printf("%v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	signals.Complex.IsPressed.Server = !signals.Complex.IsPressed.Server
	jsondata, _ := json.Marshal(signals)
	err := see.PatchSignals([]byte(jsondata))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func testSee(w http.ResponseWriter, r *http.Request) {
	see := datastar.NewSSE(w, r)
	event_format := `{'data': 'this is event #%v'}`
	time.Sleep(5 * time.Second)
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
