package server

import (
	"io"

	"github.com/gabrielluizsf/tcp_to_http/internal/request"
	"github.com/gabrielluizsf/tcp_to_http/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he *HandlerError) Error() string {
	return he.Message
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
