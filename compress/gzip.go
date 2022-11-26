package compress

import (
	"bytes"
	"compress/gzip"
	"log"
)

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, err
	}

	// gzip writer must be closed buf.Bytes() is called.
	err = gzipWriter.Close()
	if err != nil {
		log.Fatalf("Failed to close gzip writer. err = %v", err)
	}
	return buf.Bytes(), nil
}
