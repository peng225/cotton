package web

import (
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/google/uuid"
	cpath "github.com/peng225/cotton/path"
	"github.com/peng225/cotton/storage"
)

var memStore storage.MemoryStore

func init() {
	memStore = *storage.NewMemoryStore()
}

func StartServer(port int) {
	portStr := strconv.Itoa(port)

	http.HandleFunc("/", handler)
	log.Printf("Start server. port = %s\n", portStr)
	log.Println(http.ListenAndServe(":"+portStr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if !cpath.Valid(r.URL.Path) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)
	case http.MethodDelete:
		deleteHandler(w, r)
	case http.MethodHead:
		headHandler(w, r)
	default:
		log.Printf("Invalid method: %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	data, err := memStore.Get(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
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
	key := path.Join(r.URL.Path, uuid.New().String())
	if !cpath.Valid(key) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	memStore.Add(key, body)
	w.Header().Add("Location", key)
	w.WriteHeader(http.StatusCreated)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	memStore.Delete(r.URL.Path)
}

func headHandler(w http.ResponseWriter, r *http.Request) {
	data, err := memStore.Get(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Length", strconv.Itoa(len(data)))
}
