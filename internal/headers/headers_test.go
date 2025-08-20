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
		assert.Equal(t, n, 52)
		assert.True(t, done)

		value, ok := headers.Get("Host")
		assert.True(t, ok)
		assert.Equal(t, value, "localhost:42069")

		value, ok = headers.Get("Foofoo")
		assert.True(t, ok)
		assert.Equal(t, value, "barbar")

		value, ok = headers.Get("Missingkey")
		assert.False(t, ok)
		assert.Equal(t, value, stringx.Empty)
	})
	t.Run("Invalid spacing header", func(t *testing.T) {
		headers := New()
		data := []byte("       Host : localhost:42069       \r\n\r\n")
		n, done, err := headers.Parse(data)
		assert.Error(t, err)
		assert.Equal(t, n, 0)
		assert.False(t, done)
	})
	t.Run("Invalid Caracter header key", func(t *testing.T) {
		headers := New()
		data := []byte("H©st: localhost:42069\r\n\r\n")
		n, done, err := headers.Parse(data)
		assert.Error(t, err)
		assert.Equal(t, n, 0)
		assert.False(t, done)
	})
}
