package request

import (
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/stringx"
)

func TestRequestLineParse(t *testing.T) {
	reader := &chunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept:*/*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := NewFromReader(reader)
	assert.NoError(t, err)
	assert.Equal(t,r.Line.Method, "GET")
	assert.Equal(t, r.Line.Target, "/")
	assert.Equal(t, r.Line.Version, "1.1")
	
	reader = &chunkReader{
		data: "GET /coffe HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept:*/*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err = NewFromReader(reader)
	assert.NoError(t, err)
	assert.Equal(t,r.Line.Method, "GET")
	assert.Equal(t, r.Line.Target, "/coffe")
	assert.Equal(t, r.Line.Version, "1.1")

	_, err = NewFromReader(stringx.NewReader("/coffe HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept:*/*\r\n\r\n"))
	assert.Error(t, err)
}
