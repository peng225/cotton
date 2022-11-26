package web

import (
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/peng225/cotton/compress"
	cpath "github.com/peng225/cotton/path"
	"github.com/peng225/cotton/storage"
)

var (
	memStore       storage.MemoryStore
	dumpPostedData bool
)

func init() {
	memStore = *storage.NewMemoryStore()
}

func StartServer(port int, dump bool) {
	portStr := strconv.Itoa(port)
	dumpPostedData = dump

	http.HandleFunc("/", handler)
	http.HandleFunc("/ready", readyHandler)
	log.Printf("Start server. port = %s\n", portStr)
	log.Println(http.ListenAndServe(":"+portStr, nil))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !cpath.Valid(r.URL.Path) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		data, err = compress.Compress(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Encoding", "gzip")
	}
	w.Header().Add("Content-Length", strconv.Itoa(len(data)))
	writtenLength, err := w.Write(data)
	if err != nil || writtenLength != len(data) {
		log.Printf("Failed to write body. writtenLength = %d, err = %v", writtenLength, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
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
	if dumpPostedData {
		log.Printf("Posted data dump:\n%v\n", hex.Dump(body))
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
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		data, err = compress.Compress(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Encoding", "gzip")
	}
	w.Header().Add("Content-Length", strconv.Itoa(len(data)))
}
