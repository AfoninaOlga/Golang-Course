package handler

import (
	"encoding/json"
	"github.com/AfoninaOlga/xkcd/internal/core/port"
	"log"
	"net/http"
)

type XkcdHandler struct {
	svc port.ComicService
}

func NewXkcdHandler(svc port.ComicService) *XkcdHandler {
	return &XkcdHandler{svc: svc}
}

func (xh *XkcdHandler) Search(w http.ResponseWriter, req *http.Request) {
	text := req.URL.Query().Get("search")

	if text == "" {
		log.Println("Got empty query")
		http.Error(w, "Search query should not be empty", http.StatusBadRequest)
		return
	}

	comics := xh.svc.Search(text)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comics); err != nil {
		log.Panic("Error encoding found URLs:", err)
	}
}

func (xh *XkcdHandler) Update(w http.ResponseWriter, req *http.Request) {
	added := xh.svc.LoadComics()
	resp := struct {
		Added int `json:"added"`
	}{Added: added}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Panic("Error encoding response:", err)
		return
	}
}
