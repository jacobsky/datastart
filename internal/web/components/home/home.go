package home

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

var words = []string{
	"hello", "world",
}

func AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /words", spawn_words)
	mux.HandleFunc("DELETE /word/{id}", delete_word)
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

func spawn_words(w http.ResponseWriter, r *http.Request) {
	slog.Info("spawn_words")
	see := datastar.NewSSE(w, r)
	i := 0
	for {
		index := i % len(words)

		err := see.PatchElementTempl(CharacterBlock(strconv.Itoa(i), words[index], 20, 100, 0, 3, 10))
		if err != nil {
			slog.Error("Error patching", "error", err.Error())
		}
	}
}

func delete_word(w http.ResponseWriter, r *http.Request) {

	slog.Info("spawn_words")
}
