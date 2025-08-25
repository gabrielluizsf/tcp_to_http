package server

import (
	"github.com/gabrielluizsf/tcp_to_http/internal/request"
	"github.com/gabrielluizsf/tcp_to_http/internal/response"
	"github.com/i9si-sistemas/stringx"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
	DELETE Method = "DELETE"
)

func MethodMiddleware(method Method, handler Handler) Handler {
	return func(w *response.Writer, req *request.Request) {
		if !stringx.New(req.Line.Method).Equal(string(method)) {
			headers := response.GetDefaultHeaders(0)
			w.WriteStatusLine(response.StatusMethodNotAllowed)
			w.WriteHeaders(headers)
			return
		}
		handler(w, req)
	}
}

func (s *Server) Get(endpoint string, handler Handler) {
	s.handlers[endpoint] = MethodMiddleware(GET, handler)
}

func (s *Server) Post(endpoint string, handler Handler) {
	s.handlers[endpoint] = MethodMiddleware(POST, handler)
}

func (s *Server) Put(endpoint string, handler Handler) {
	s.handlers[endpoint] = MethodMiddleware(PUT, handler)
}

func (s *Server) Patch(endpoint string, handler Handler) {
	s.handlers[endpoint] = MethodMiddleware(PATCH, handler)
}

func (s *Server) Delete(endpoint string, handler Handler) {
	s.handlers[endpoint] = MethodMiddleware(DELETE, handler)
}
