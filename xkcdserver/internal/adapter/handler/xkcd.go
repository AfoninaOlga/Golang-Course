package handler

import (
	"encoding/json"
	"github.com/AfoninaOlga/xkcd/xkcdserver/internal/core/port"
	"log"
	"net/http"
	"sync"
)

type XkcdHandler struct {
	svc port.ComicService
	mtx *sync.Mutex
}

func NewXkcdHandler(svc port.ComicService) *XkcdHandler {
	return &XkcdHandler{svc: svc, mtx: &sync.Mutex{}}
}

func (xh *XkcdHandler) Search(w http.ResponseWriter, req *http.Request) {
	text := req.URL.Query().Get("search")

	if text == "" {
		log.Println("Got empty query")
		http.Error(w, "Search query should not be empty", http.StatusBadRequest)
		return
	}

	comics := xh.svc.Search(req.Context(), text)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comics); err != nil {
		log.Panic("Error encoding found URLs:", err)
	}
}

func (xh *XkcdHandler) Update(w http.ResponseWriter, req *http.Request) {
	if xh.mtx.TryLock() {
		defer xh.mtx.Unlock()
		added := xh.svc.LoadComics(req.Context())
		resp := struct {
			Added int `json:"added"`
		}{Added: added}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Panic("Error encoding response:", err)
			return
		}
	} else {
		http.Error(w, "Update is already in progress", http.StatusServiceUnavailable)
	}
}
