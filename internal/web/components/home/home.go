package home

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

func AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /words", spawn_words)
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

type Words struct {
	Words string `json:"words"`
}

func spawn_words(w http.ResponseWriter, r *http.Request) {
	slog.Info("spawn_words")
	storage := &Words{}
	if err := datastar.ReadSignals(r, storage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	words := strings.Split(storage.Words, " ")
	sse := datastar.NewSSE(w, r)

	err := sse.PatchSignals([]byte(`{words: ''}`))
	if err != nil {
		slog.Error("Could not patch signals", "error", err.Error())
	}
	// Execute JavaScript in the browser
	err = sse.ExecuteScript(`console.log('Starting to spawn words shortly')`)
	if err != nil {
		slog.Error("Error patching", "error", err.Error())
	}
	var components = make([]templ.Component, len(words))
	const font_size = 20
	characters := 0

	for i, word := range words {
		id := fmt.Sprintf("component_%v", i)
		components[i] = CharacterBlock(
			id,
			word,
			font_size,
			characters*font_size/2,
			-characters*3,
			0,
			6,
		)
		characters += len(word) + 1
	}
	slog.Info("Patch in the element for", "words", words, "components", components)

	err = sse.PatchElementTempl(WordRain(components))
	if err != nil {
		slog.Error("Patching issue", "error", err.Error())
	}

	time.Sleep(5 * time.Second)

	_ = sse.PatchElementTempl(StartWordRain())
	// appendoption, err := datastar.ElementPatchModeFromString("mode append")
	// if err != nil {
	// 	slog.Error("error patching", "error", err.Error())
	// 	return
	// }
	// i := 0
	// for {
	// 	index := i % len(words)
	//
	// 	err = sse.PatchElementTempl(
	// 		CharacterBlock(strconv.Itoa(i), words[index], 20, 100, 0, 3, 10),
	//		appendoption
	// 	)
	//
	// 	if err != nil {
	// 		slog.Error("Error patching", "error", err.Error())
	// 	}
	// 	time.Sleep(5 * time.Second)
	// }
}
