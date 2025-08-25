package response

import (
	"fmt"
	"io"

	"github.com/gabrielluizsf/tcp_to_http/pkg/headers"
)

func WriteHeaders(
	w io.Writer,
	h headers.Headers,
) (err error) {
	b := []byte{}
	for key, value := range h {
		b = fmt.Appendf(b, "%s: %s\r\n", key, value)
	}
	b = fmt.Append(b, "\r\n")
	_, err = w.Write(b)
	return
}
