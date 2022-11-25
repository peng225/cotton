package web

import (
	"io"
	"log"
	"net/http"

	"github.com/peng225/cotton/storage"
)

var memStore storage.MemoryStore

func init() {
	memStore = *storage.NewMemoryStore()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)
	default:
		log.Printf("Invalid method: %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	data, err := memStore.Get(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Write(data)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	memStore.Add(r.URL.Path, body)
}
