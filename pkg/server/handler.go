package server

import (
	"github.com/gabrielluizsf/tcp_to_http/pkg/request"
	"github.com/gabrielluizsf/tcp_to_http/pkg/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he *HandlerError) Error() string {
	return he.Message
}

type Handler func(w *response.Writer, req *request.Request)
