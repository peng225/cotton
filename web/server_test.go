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

func TestPutHandler(t *testing.T) {
	exampleKey := "/test/data/3905f7d8-852f-4df3-bd8c-2fbe8e54c01a"
	cases := []struct {
		description        string
		key                string
		expectedStatusCode int
	}{
		{
			description:        "Good UUID",
			key:                exampleKey,
			expectedStatusCode: http.StatusCreated,
		},
		{
			description:        "UUID is too short",
			key:                exampleKey[:len(exampleKey)-2],
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:        "UUID contains no hyphens",
			key:                "/test/data/3905f7d8852f4df3bd8c2fbe8e54c01a",
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			w := NewTestResponseWriter()
			req, err := http.NewRequest(http.MethodPut, "http://"+path.Join("localhost:8080", tt.key), bytes.NewBuffer([]byte("test data")))
			require.NoError(t, err)
			putHandler(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.statusCode)
			if tt.expectedStatusCode == http.StatusCreated {
				assert.Equal(t, exampleKey, w.Header().Get("Location"))
			}
		})
	}
}
