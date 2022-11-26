package web

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup() {
	go func() {
		StartServer(8080, false)
	}()
}

func waitServerReady(t *testing.T) {
	t.Helper()
	require.Eventually(t, func() bool {
		resp, err := http.Get("http://localhost:8080/ready")
		if err != nil {
			return false
		}
		return resp.StatusCode == http.StatusOK
	}, time.Second*10, time.Millisecond*200)
}

func TestAllMethods(t *testing.T) {
	waitServerReady(t)
	// POST
	resp, err := http.Post("http://localhost:8080/test/data", "text/plain;charset=UTF-8", bytes.NewReader([]byte("test data")))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	location, err := resp.Location()
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(location.Path, "/test/data"), "Location: ", location.Path)

	// GET
	url := "http://" + path.Join("localhost:8080", location.Path)
	resp, err = http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, []byte("test data"), body)
	assert.True(t, resp.Uncompressed)
	resp.Body.Close()

	// HEAD
	resp, err = http.Head(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(len([]byte("test data"))), resp.ContentLength)

	// DELETE
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// GET again
	resp, err = http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}
