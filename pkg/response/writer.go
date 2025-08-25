package response

import (
	"encoding/json"
	"io"

	"github.com/gabrielluizsf/tcp_to_http/pkg/headers"
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

func (w *Writer) Send(data []byte) (int, error) {
	return w.WriteBody(data)
}

func (w *Writer) SendString(data string) (int, error) {
	return w.WriteBody([]byte(data))
}

func (w *Writer) JSON(data any) (int, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return w.WriteBody(b)
}
