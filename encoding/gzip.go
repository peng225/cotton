package encoding

import (
	"compress/gzip"
	"io"
)

var (
	gzipWriter *gzip.Writer
)

func init() {
	gzipWriter = gzip.NewWriter(nil)
}

func GzipCompress(data []byte) (io.ReadCloser, chan error) {
	reader, writer := io.Pipe()
	gzipWriter.Reset(writer)
	errch := make(chan error, 1)

	go func() {
		defer close(errch)
		sentErr := false
		sendErr := func(err error) {
			if !sentErr {
				errch <- err
				sentErr = true
			}
		}

		_, err := gzipWriter.Write(data)
		if err != nil {
			sendErr(err)
		}

		err = gzipWriter.Close()
		if err != nil {
			sendErr(err)
		}

		err = writer.Close()
		if err != nil && err != io.ErrClosedPipe {
			sendErr(err)
		}
	}()

	return reader, errch
}
