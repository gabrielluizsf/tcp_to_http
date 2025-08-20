package headers

import (
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/stringx"
)

func TestHeadersParse(t *testing.T) {
	t.Run("Valid single header", func(t *testing.T) {
		headers := New()
		data := []byte("Host: localhost:42069\r\nFooFoo:            barbar\r\n\r\n")
		n, done, err := headers.Parse(data)
		assert.NoError(t, err)
		assert.NotNil(t, headers)
		assert.Equal(t, headers["Host"], "localhost:42069")
		assert.Equal(t, headers["FooFoo"], "barbar")
		assert.Equal(t, headers["MissingKey"], stringx.Empty)
		assert.Equal(t, n, 52)
		assert.True(t, done)

	})
	t.Run("Invalid spacing header", func(t *testing.T) {
		headers := New()
		data := []byte("       Host : localhost:42069       \r\n\r\n")
		n, done, err := headers.Parse(data)
		assert.Error(t, err)
		assert.Equal(t, n, 0)
		assert.False(t, done)
	})
}
