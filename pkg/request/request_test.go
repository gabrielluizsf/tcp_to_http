package request

import (
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/stringx"
)

func TestRequestLineParse(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept:*/*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := NewFromReader(reader)
	assert.NoError(t, err)
	assert.Equal(t, r.Line.Method, "GET")
	assert.Equal(t, r.Line.Target, "/")
	assert.Equal(t, r.Line.Version, "1.1")

	reader = &chunkReader{
		data:            "GET /coffe HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept:*/*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err = NewFromReader(reader)
	assert.NoError(t, err)
	assert.Equal(t, r.Line.Method, "GET")
	assert.Equal(t, r.Line.Target, "/coffe")
	assert.Equal(t, r.Line.Version, "1.1")

	_, err = NewFromReader(stringx.NewReader("/coffe HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept:*/*\r\n\r\n"))
	assert.Error(t, err)
}

func TestRequestHeadersParse(t *testing.T) {
	t.Run("Standard Headers", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			numBytesPerRead: 3,
		}
		r, err := NewFromReader(reader)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, "localhost:42069", r.Headers["host"])
		assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
		assert.Equal(t, "*/*", r.Headers["accept"])
	})

	t.Run("Malformed Header", func(t *testing.T) {
		reader := &chunkReader{
			data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
			numBytesPerRead: 3,
		}
		_, err := NewFromReader(reader)
		assert.Error(t, err)
	})
}

func TestRequestBodyParse(t *testing.T) {
	t.Run("Standard Body", func(t *testing.T) {
		reader := &chunkReader{
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"Content-Length: 13\r\n" +
				"\r\n" +
				"hello world!\n",
			numBytesPerRead: 3,
		}
		r, err := NewFromReader(reader)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, "hello world!\n", string(r.Body))
	})

	t.Run("Body shorter than reported content length", func(t *testing.T) {
		reader := &chunkReader{
			data: "POST /submit HTTP/1.1\r\n" +
				"Host: localhost:42069\r\n" +
				"Content-Length: 20\r\n" +
				"\r\n" +
				"partial content",
			numBytesPerRead: 3,
		}
		_, err := NewFromReader(reader)
		assert.Error(t, err)
	})
}
