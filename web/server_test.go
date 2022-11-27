package web

import (
	"bytes"
	"io"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestResponseWriter struct {
	statusCode int
	header     http.Header
}

func NewTestResponseWriter() *TestResponseWriter {
	return &TestResponseWriter{
		header: make(http.Header),
	}
}

func (trw *TestResponseWriter) Header() http.Header {
	return trw.header
}

func (trw *TestResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (trw *TestResponseWriter) WriteHeader(statusCode int) {
	trw.statusCode = statusCode
}

func TestPostHandler(t *testing.T) {
	largeData := make([]byte, maxBlobSize+1)
	cases := []struct {
		description        string
		body               io.Reader
		expectedStatusCode int
	}{
		{
			description:        "nil body",
			body:               nil,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description:        "empty body",
			body:               bytes.NewBuffer([]byte{}),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:        "too large data",
			body:               bytes.NewBuffer(largeData),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			w := NewTestResponseWriter()
			req, err := http.NewRequest(http.MethodPost, "http://"+path.Join("localhost:8080", "/test/path"), tt.body)
			require.NoError(t, err)
			postHandler(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.statusCode)
		})
	}
}
