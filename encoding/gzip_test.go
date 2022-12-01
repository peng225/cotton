package encoding

import (
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzipCompress(t *testing.T) {
	cases := []struct {
		description string
		data        []byte
	}{
		{
			description: "success",
			data:        []byte{0, 0, 1, 1},
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			reader, errch := GzipCompress(tt.data)
			gzipReader, err := gzip.NewReader(reader)
			require.NoError(t, err)

			readData, err := io.ReadAll(gzipReader)
			require.NoError(t, err)
			require.NoError(t, <-errch)
			assert.Equal(t, tt.data, readData)
		})
	}
}
