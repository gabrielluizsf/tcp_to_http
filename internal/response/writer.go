package response

import (
	"io"

	"github.com/gabrielluizsf/tcp_to_http/internal/headers"
)

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	return WriteStatusLine(w.w, statusCode)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return WriteHeaders(w.w, headers)
}

func (w *Writer) WriteBody(data []byte) (int, error) {
	return w.w.Write(data)
}
