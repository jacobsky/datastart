package home

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

//	type GameBoard struct {
//		Width
//		[Position]string
//	}
type CharacterPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var characterPosition = CharacterPosition{0, 0}

func AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("PATCH /move", move)
	mux.HandleFunc("GET /position", position)
	mux.Handle("/", NewHandler())
}

type Handler struct{}

func NewHandler() http.Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templ.Handler(Home()).ServeHTTP(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func move(w http.ResponseWriter, r *http.Request) {
	log.Printf("Info {}")
}

func position(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	positionjson, err := json.Marshal(characterPosition)
	if err != nil {
		log.Printf("json error: %v", err)
	}
	err = sse.PatchSignals([]byte(positionjson))
	if err != nil {
		log.Printf("json error: %v", err)
	}
}
