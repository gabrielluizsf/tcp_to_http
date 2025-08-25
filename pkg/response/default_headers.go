package response

import (
	"fmt"

	"github.com/gabrielluizsf/tcp_to_http/pkg/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.New()
	h.Set("Content-Length", fmt.Sprint(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	return h
}
