package web

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/peng225/cotton/encoding"
	cpath "github.com/peng225/cotton/path"
	"github.com/peng225/cotton/storage"
)

const (
	maxBlobSize = 10 * 1024 * 1024
	chunkSize   = 100 * 1024
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

func readChunk(readCloser io.ReadCloser, buf []byte) (int, error) {
	readSize := 0
	for readSize < chunkSize {
		n, err := readCloser.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Read failed. err = %v", err)
			return 0, err
		}
		readSize += n
		if err == io.EOF {
			return readSize, io.EOF
		}
		buf = buf[n:]
		time.Sleep(100 * time.Microsecond)
	}
	return readSize, nil
}

func chunkedTransfer(w http.ResponseWriter, buf []byte) error {
	n, err := w.Write(buf)
	if err != nil {
		return err
	} else if n != len(buf) {
		return fmt.Errorf("Write length is too short. n = %d, len(buf) = %d", n, len(buf))
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("w should implement http.Flushter interface.")
	}
	flusher.Flush()

	return nil
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	data, err := memStore.Get(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Add("Content-Encoding", "gzip")
		// Because "Transfer-Encoding: chunked" is automatically added by calling Flusher.Flush(),
		// it is not needed to add it explicitly.
		readCloser, errch := encoding.GzipCompress(data)
		defer readCloser.Close()
		buf := make([]byte, chunkSize)
		for {
			n, readErr := readChunk(readCloser, buf)
			if readErr != nil && readErr != io.EOF {
				log.Printf("readChunk failed. err = %v", readErr)
				break
			}
			if n != 0 {
				ctErr := chunkedTransfer(w, buf[:n])
				if ctErr != nil {
					log.Printf("chunkedTransfer failed. err = %v", ctErr)
					break
				}
			}

			if readErr == io.EOF {
				break
			}
		}
		if err := <-errch; err != nil {
			log.Printf("GzipCompress failed. err = %v", err)
		}
	} else {
		w.Header().Add("Content-Length", strconv.Itoa(len(data)))
		writtenLength, err := w.Write(data)
		if err != nil || writtenLength != len(data) {
			log.Printf("Failed to write body. writtenLength = %d, err = %v", writtenLength, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(body) == 0 || len(body) > maxBlobSize {
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

	// Accept-Encoding header is ignored in case of HEAD request.
	w.Header().Add("Content-Length", strconv.Itoa(len(data)))
}
